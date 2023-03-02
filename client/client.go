// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	glmClient "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	"github.com/hpe-hcss/lh-cdc-singularity/restclient"
	log "github.com/hpe-storage/common-host-libs/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

const AUTH_SVC_INFO string = "/info/authsvcinfo"

type ClientInterface interface {
	CreateVolume(volume *model.Volume) (*model.Volume, error)
	DeleteVolume(volumeId string) (*model.Volume, error)
	GetVolume(volumeId string) (*model.Volume, error)
	ListVolumes() (*[]model.Volume, error)
	CreateVolumeAttachment(attachment *model.VolumeAttachment) (*model.VolumeAttachment, error)
	DeleteVolumeAttachment(attachmentId string) (*model.VolumeAttachment, error)
	GetVolumeAttachment(attachmentId string) (*model.VolumeAttachment, error)
	ListVolumeAttachments() (*[]model.VolumeAttachment, error)
	ListVolumeFlavors() (*[]model.VolumeFlavor, error)
	Login() error
	Logout()
	GetGlmCredentials() (map[string]string, error)
}

var UndefinedResponseError = errors.New("undefined response type")
var TokenExpiredError = errors.New("Token is expired")
var UndefinedResponseMsg = "undefined response type"

type errorMsg struct {
	Message string
}

type Client struct {
	Url          string
	UserName     string
	Password     string
	SessionToken string
	MembershipID string
}

func NewClient(URL, username, password, membershipID string) ClientInterface {
	return &Client{
		URL,
		username,
		password,
		"",
		membershipID,
	}
}

func (cli *Client) Login() error {
	log.Infof("Login function")
	client := restclient.NewRestClient()
	userInfo := restclient.UserInfo{
		UserName: cli.UserName,
		UserPwd:  cli.Password,
	}
	header := map[string]string{}

	// Get authorization service information
	statusCode, responseBody, err := client.ExecuteRestRequest("GET", AUTH_SVC_INFO, "",
		"", userInfo, cli.Url, "", 0, nil, header, "")
	if err != nil {
		msg := fmt.Sprintf("ExecuteRestRequest error: %+v", err)
		log.Errorf(msg)
		return errors.New(msg)
	}
	log.Infof("responseBody %+v", responseBody)
	if statusCode != restclient.StatusCodeOk {
		return status.Errorf(codes.Internal, fmt.Sprintf("failed to get auth service "+
			"info, error: %v", statusCode))
	}

	// parse the authsvcinfo
	var info model.AuthSvcInfo
	if err := json.Unmarshal(responseBody, &info); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("failed to parse auth service info, error: %v", statusCode))
	}

	log.Infof("AuthSvcInfo is %+v", info)
	authUrl := info.AuthURL
	authClientId := info.AuthClientID
	authAudience := info.AuthAudience

	var authReq model.AuthRequest

	authReq.ClientID = authClientId
	authReq.Username = cli.UserName
	authReq.Password = cli.Password
	authReq.Realm = constants.PASSWORD_REALM
	authReq.GrantType = constants.GRANT_PASSWORD_REALM
	authReq.Scope = constants.OPEN_ID
	authReq.Audience = authAudience

	body, _ := json.Marshal(&authReq)

	if statusCode, responseBody, err = client.ExecuteRestRequest("POST", "/oauth/token", "",
		"", userInfo, authUrl, "", 0, body, header, ""); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("failed to get OAUTH token, error: %v", err))
	}
	if statusCode != restclient.StatusCodeOk {
		return status.Errorf(codes.Internal, fmt.Sprintf("failed to get OAUTH token, error: %v", statusCode))
	}

	// parse JWT
	var resp model.AuthResponse
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("failed to parse OAUTH token, error: %v", err))
	}
	log.Infof("resp is %+v", resp)
	cli.SessionToken = resp.IDToken
	// save session token in the conf file.
	glm := &model.GLMCredDetails{}
	err = glm.UpdateGLMCredentials(constants.SESSION_TOKEN, resp.IDToken)
	if err != nil {
		msg := fmt.Sprintf("error: %+v", err)
		log.Errorf(msg)
		return errors.New(msg)
	}
	return nil
}

func (cli *Client) Logout() {
	//Need to write the definition
	log.Infof("Logout")
}

func (cli *Client) CreateVolume(vol *model.Volume) (*model.Volume, error) {
	log.Infof("Create volume req %v\n", vol)
	var errMsg errorMsg
	volume := glmClient.NewVolume{}
	volume.Name = vol.Name
	volume.Description = vol.Description
	volume.FlavorID = vol.FlavorID
	volume.Capacity = vol.Capacity
	volume.LocationID = vol.LocationID
	ctx, r, err := GetREST(cli.Url, cli.UserName, cli.MembershipID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("context value: %v\n", ctx)
	log.Infof("r value: %v\n", r)
	result, _, err := r.VolumesApi.Add(ctx, volume)
	if err != nil && err.Error() == UndefinedResponseMsg {
		msg := fmt.Sprintf("create volume error %+v", err)
		log.Errorf(msg)
		return nil, UndefinedResponseError
	} else if err != nil {
		_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
		msg := fmt.Sprintf("Create volume error: %+v", errMsg.Message)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("Create volume response %v", result)
	var resp glmClient.Volume
	state := result.State
	for iteration := constants.MIN_RETRY_COUNT; iteration < constants.MAX_RETRY_COUNT; iteration++ {
		if state != "allocated" {
			time.Sleep(constants.SLEEP_TIME * time.Second)
			resp, _, err = r.VolumesApi.GetByID(ctx, result.ID)
			if err != nil {
				_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
				msg := fmt.Sprintf("get volume failed with error: %+v\n", errMsg.Message)
				log.Errorf(msg)
				return nil, errors.New(msg)
			}
		} else {
			break
		}
		state = resp.State
	}
	if state != "allocated" {
		msg := fmt.Sprintf("Create volume failed with volume state : %v", state)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	volumeResp := model.CreateResponse(resp, constants.CREATE)
	//covert create volume results into model.Volume and return it
	return volumeResp, nil
}

func (cli *Client) DeleteVolume(volumeID string) (*model.Volume, error) {
	log.Infof("DeleteVolume volume id: %s", volumeID)
	var errMsg errorMsg
	ctx, r, err := GetREST(cli.Url, cli.UserName, cli.MembershipID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("context value: %v\n", ctx)
	log.Infof("r value: %v\n", r)
	httpResponse, err := r.VolumesApi.Delete(ctx, volumeID)
	if err != nil && err.Error() == UndefinedResponseMsg {
		msg := fmt.Sprintf("delete volume error %+v", err)
		log.Errorf(msg)
		return nil, UndefinedResponseError
	} else if err != nil {
		_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
		msg := fmt.Sprintf("Delete volume failed with error: %+v", errMsg.Message)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	log.Infof("Delete volume response %v", httpResponse)
	if httpResponse.StatusCode != constants.STATUS_OK {
		msg := fmt.Sprintf("Delete volume failed with status code = %d", httpResponse.StatusCode)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	var resp glmClient.Volume
	resp.ID = volumeID
	resp.State = constants.VOLUME_STATE_DELETING
	state := resp.State
	for iteration := constants.MIN_RETRY_COUNT; iteration < constants.MAX_RETRY_COUNT; iteration++ {
		if state != constants.STATE_DELETED {
			time.Sleep(constants.SLEEP_TIME * time.Second)
			resp, _, err = r.VolumesApi.GetByID(ctx, resp.ID)
			if err != nil {
				_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
				msg := fmt.Sprintf("Get volume failed with error: %+v", errMsg.Message)
				log.Errorf(msg)
				return nil, errors.New(msg)
			}
		} else {
			break
		}
		state = resp.State
	}
	if state != constants.STATE_DELETED {
		msg := fmt.Sprintf("Delete volume failed with volume state : %v", state)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	volumeResp := model.CreateResponse(resp, constants.DELETE)
	log.Infof("Delete volume response structure %+v", volumeResp)
	//covert create volume results into model.Volume and return it
	return volumeResp, nil
}

func (cli *Client) GetVolume(volumeID string) (*model.Volume, error) {
	log.Infof("GetVolume volume id: %s", volumeID)
	var errMsg errorMsg
	ctx, r, err := GetREST(cli.Url, cli.UserName, cli.MembershipID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("context value: %v\n", ctx)
	log.Infof("r value: %v\n", r)
	result, _, err := r.VolumesApi.GetByID(ctx, volumeID)
	if err != nil && err.Error() == UndefinedResponseMsg {
		msg := fmt.Sprintf("get volume error %+v", err)
		log.Errorf(msg)
		return nil, UndefinedResponseError
	} else if err != nil {
		_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
		msg := fmt.Sprintf("Get volume failed with error: %+v", errMsg.Message)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("Get volume response %v", result)
	volumeResp := model.CreateResponse(result, constants.GET)
	log.Infof("Get volume response structure %+v", volumeResp)
	//covert get volume results into model.Volume and return it
	return volumeResp, nil
}

func (cli *Client) ListVolumes() (*[]model.Volume, error) {
	log.Infof("List Volume")
	var errMsg errorMsg
	ctx, r, err := GetREST(cli.Url, cli.UserName, cli.MembershipID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("context value: %v\n", ctx)
	log.Infof("r value: %v\n", r)
	result, _, err := r.VolumesApi.List(ctx)
	if err != nil && err.Error() == UndefinedResponseMsg {
		msg := fmt.Sprintf("list volume error %+v", err)
		log.Errorf(msg)
		return nil, UndefinedResponseError
	} else if err != nil {
		_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
		msg := fmt.Sprintf("List volumes failed with error: %+v", errMsg.Message)
		log.Errorf(msg)
		return nil, err
	}
	log.Infof("List volume response %v", result)
	volumesList := []model.Volume{}
	for _, item := range result {
		volumesList = append(volumesList, *model.CreateResponse(item, constants.LIST))
	}
	log.Infof("List volume response structure %+v", volumesList)
	return &volumesList, nil
}

func (cli *Client) CreateVolumeAttachment(
	volumeAttachment *model.VolumeAttachment) (*model.VolumeAttachment, error) {
	log.Infof("create volume attachment: %v", volumeAttachment)
	var errMsg errorMsg
	volAttachment := glmClient.NewVolumeAttachment{}
	volAttachment.Name = volumeAttachment.Name
	volAttachment.VolumeID = volumeAttachment.VolumeID
	protocol := glmClient.ProtocolParameters{}
	protocol.Protocol = constants.PROTOCOL_FUSE
	volAttachment.Protocol = protocol
	ctx, r, err := GetREST(cli.Url, cli.UserName, cli.MembershipID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugf("context value: %v\n", ctx)
	log.Infof("r value: %v\n", r)
	result, _, err := r.VolumeAttachmentsApi.Add(ctx, volAttachment)
	if err != nil && err.Error() == UndefinedResponseMsg {
		msg := fmt.Sprintf("create volume attachment failed with error %+v", err)
		log.Errorf(msg)
		return nil, UndefinedResponseError
	} else if err != nil {
		_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
		msg := fmt.Sprintf("attach volumes failed with error: %+v", errMsg.Message)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("Create volume attachment response %v", result)
	var resp glmClient.VolumeAttachment
	state := result.State
	for iteration := constants.MIN_RETRY_COUNT; iteration < constants.MAX_RETRY_COUNT; iteration++ {
		if state != constants.VOLUME_ATTACHMENT_STATE_READY {
			time.Sleep(constants.SLEEP_TIME * time.Second)
			resp, _, err = r.VolumeAttachmentsApi.GetByID(ctx, result.ID)
			if err != nil {
				_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
				msg := fmt.Sprintf("Get volume attachment failed with error: %+v", errMsg.Message)
				log.Errorf(msg)
				return nil, errors.New(msg)
			}
		} else {
			break
		}
		state = resp.State
	}
	if state != constants.VOLUME_ATTACHMENT_STATE_READY {
		msg := fmt.Sprintf("Create volume attachment failed with state : %v", state)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	volumeResp := model.CreateVolumeAttachmentResponse(resp, constants.CREATE)
	//covert create volume results into model.Volume and return it
	return volumeResp, nil
	//return &model.VolumeAttachment{}, errors.New("not implemented")
}

func (cli *Client) DeleteVolumeAttachment(attachmentId string) (*model.VolumeAttachment, error) {
	log.Infof("delete volume attachment id: %v", attachmentId)
	var errMsg errorMsg
	ctx, r, err := GetREST(cli.Url, cli.UserName, cli.MembershipID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugf("context value: %v\n", ctx)
	log.Infof("r value: %v\n", r)
	httpResponse, err := r.VolumeAttachmentsApi.Delete(ctx, attachmentId)
	if err != nil && err.Error() == UndefinedResponseMsg {
		msg := fmt.Sprintf("delete volume attachment failed with error %+v", err)
		log.Errorf(msg)
		return nil, UndefinedResponseError
	} else if err != nil {
		_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
		msg := fmt.Sprintf("Delete volume attachment failed with error: %+v", errMsg.Message)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	if httpResponse.StatusCode != constants.STATUS_OK {
		msg := fmt.Sprintf("Delete volume attachment failed with status code = %d", httpResponse.StatusCode)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	log.Infof("Delete volume attachment response %v", httpResponse)
	var resp glmClient.VolumeAttachment
	resp.ID = attachmentId
	resp.State = constants.STATE_DELETED
	volumeResp := model.CreateVolumeAttachmentResponse(resp, constants.CREATE)
	log.Infof("Delete volume response structure %+v", volumeResp)
	//covert create volume results into model.Volume and return it
	return volumeResp, nil
}

func (cli *Client) GetVolumeAttachment(attachmentId string) (*model.VolumeAttachment, error) {
	log.Infof("get volume attachment id: %v", attachmentId)
	var errMsg errorMsg
	ctx, r, err := GetREST(cli.Url, cli.UserName, cli.MembershipID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugf("context value: %v\n", ctx)
	log.Infof("r value: %v\n", r)
	result, _, err := r.VolumeAttachmentsApi.GetByID(ctx, attachmentId)
	if err != nil && err.Error() == UndefinedResponseMsg {
		msg := fmt.Sprintf("get volume attachment failed with error %+v", err)
		log.Errorf(msg)
		return nil, UndefinedResponseError
	} else if err != nil {
		_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
		msg := fmt.Sprintf("Get volume attachment failed with error: %+v", errMsg.Message)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("Get volume attachment response %v", result)
	getResp := model.CreateVolumeAttachmentResponse(result, constants.GET)
	log.Infof("Get volume attachment response structure %+v", getResp)
	//covert get volume results into model.Volume and return it
	return getResp, nil
}

func (cli *Client) ListVolumeAttachments() (*[]model.VolumeAttachment, error) {
	log.Infof("list volume attachments")
	var errMsg errorMsg
	ctx, r, err := GetREST(cli.Url, cli.UserName, cli.MembershipID)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("context value: %v\n", ctx)
	log.Infof("r value: %v\n", r)
	result, _, err := r.VolumeAttachmentsApi.List(ctx)
	if err != nil && err.Error() == UndefinedResponseMsg {
		msg := fmt.Sprintf("list volume attachment failed with error %+v", err)
		log.Errorf(msg)
		return nil, UndefinedResponseError
	} else if err != nil {
		_ = json.Unmarshal(err.(glmClient.GenericOpenAPIError).Body(), &errMsg)
		msg := fmt.Sprintf("List volume attachments failed with error: %+v", errMsg.Message)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("List volume attachment response %v", result)
	volumeAttachmentsList := []model.VolumeAttachment{}
	for _, item := range result {
		volumeAttachmentsList = append(volumeAttachmentsList, *model.CreateVolumeAttachmentResponse(item, constants.LIST))
	}
	log.Infof("List volume response structure %+v", volumeAttachmentsList)
	return &volumeAttachmentsList, nil
}

func (cli *Client) ListVolumeFlavors() (*[]model.VolumeFlavor, error) {
	path := ""
	rawQuery := ""
	client := restclient.NewRestClient()
	userInfo := restclient.UserInfo{
		UserName: cli.UserName,
		UserPwd:  cli.Password,
	}

	header := map[string]string{
		"Membership": cli.MembershipID,
	}
	sessionToken, err := model.GetSessionToken()
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	statusCode, responseBody, err := client.ExecuteRestRequest(constants.GET_FLAVOR, constants.REST_VOLUME_FLAVOR_URL,
		path, rawQuery, userInfo, cli.Url, "", 0, nil, header, sessionToken)
	if err != nil {
		msg := fmt.Sprintf("volume flavor list failed with error: %+v\n", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("responseBody %+v\n", string(responseBody))
	if strings.TrimSuffix(string(responseBody), "\n") == "Token is expired" {
		log.Errorf(string(responseBody))
		return nil, TokenExpiredError
	}

	if statusCode != restclient.StatusCodeOk {
		msg := fmt.Sprintf("Volume flavor list failed with status code: %+v\n", statusCode)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	var volumeFlavorResp []model.VolumeFlavor
	err = json.Unmarshal(responseBody, &volumeFlavorResp)
	if err != nil {
		msg := fmt.Sprintf("listVolumeFlavors: unmarshal: %v", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	return &volumeFlavorResp, nil
}

func (cli *Client) GetGlmCredentials() (map[string]string, error) {
	glmCredentials := make(map[string]string)
	if cli.Url == "" || cli.UserName == "" || cli.Password == "" || cli.MembershipID == "" {
		log.Errorln("glm credentials are missing")
		return nil, errors.New("glm credentials are missing")
	}
	glmCredentials["URL"] = cli.Url
	glmCredentials["USER_NAME"] = cli.UserName
	glmCredentials["PASSWORD"] = cli.Password
	glmCredentials["MEMBERSHIP_ID"] = cli.MembershipID
	return glmCredentials, nil
}
