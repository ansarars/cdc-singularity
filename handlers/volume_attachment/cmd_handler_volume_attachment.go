// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package ss

import (
	"errors"
	"fmt"
	"github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	"github.com/hpe-hcss/lh-cdc-singularity/utils"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type VolumeAttachmentHandler interface {
	MakeResource(string, map[string]interface{}) (*model.VolumeAttachment, error)
	ValidateResource(string, *model.VolumeAttachment) error
	Execute(*model.VolumeAttachment, client.ClientInterface) (interface{}, error)
}

type CmdHandlerVolumeAttachment struct {
	args  []string
	opMap map[string]VolumeAttachmentHandler
}

var supportedVolAttachmentOperations = []string{"create", "get", "delete", "list"}

func NewCmdHandlerVolumeAttachment(args []string) *CmdHandlerVolumeAttachment {
	log.Infof("NewCmdHandlerVolumeAttachment : %v", args)
	ch := &CmdHandlerVolumeAttachment{}
	opMap := map[string]VolumeAttachmentHandler{
		constants.CREATE: &CreateVolumeAttachmentHandler{},
		constants.DELETE: &DeleteVolumeAttachmentHandler{},
		constants.GET:    &GetVolumeAttachmentHandler{},
		constants.LIST:   &ListVolumeAttachmentHandler{},
	}
	ch.args = args
	ch.opMap = opMap
	return ch
}

func (ch *CmdHandlerVolumeAttachment) Handle(glmCredDetails map[string]string,
	argsMap map[string]interface{}) (interface{}, error) {
	log.Infof("Handle function")
	var resp interface{}
	operation := ch.args[0]
	if err := utils.ValidateOperations(operation, supportedVolAttachmentOperations); err != nil {
		log.Errorln(err)
		return nil, err
	}

	handler := ch.opMap[operation]
	if handler != nil {
		//validating parameters
		volumeAttachment, err := handler.MakeResource(operation, argsMap)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		err = handler.ValidateResource(operation, volumeAttachment)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}
		glmUserName, glmPassword, err := utils.GetCredentials(argsMap)
		if err != nil {
			return nil, err
		}
		cli := client.NewClient(glmCredDetails[constants.GLM_PORTAL], glmUserName,
			glmPassword, glmCredDetails[constants.MEMBERSHIP_ID])
		resp, err = handler.Execute(volumeAttachment, cli)
		if err == model.TokenError || err == client.UndefinedResponseError || err == client.TokenExpiredError {
			log.Errorf("Execute err %+v", err)
			//loginRequired = true
			//cli.Login()
			err := cli.Login()
			if err != nil {
				msg := fmt.Sprintf("Session creation failed with error: %v", err)
				log.Errorf(msg)
				return nil, errors.New(msg)
			}
			defer cli.Logout()
			//resp, err = handler.Execute(volume, glmCredDetails, glmUserName, glmPassword, loginRequired)
			resp, err = handler.Execute(volumeAttachment, cli)
			if err != nil {
				log.Errorln(err)
				return nil, err
			}
		} else if err != nil {
			log.Errorln(err)
			return nil, err
		}
	} else {
		msg := fmt.Sprintf("unsupported sub-command: %s", operation)
		log.Errorln(msg)
		return nil, errors.New(msg)
	}
	return resp, nil
}

func ValidateCreateVolumeAttachmentRequest(volumeAttachment *model.VolumeAttachment) error {
	log.Infof("ValidateCreateVolumeAttachmentRequest function")
	if volumeAttachment.Name == "" {
		msg := fmt.Sprintf("invalid value of volume attachment %s is provided", constants.ATTACHMENT_NAME)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volumeAttachment.VolumeID == "" {
		msg := fmt.Sprintf("invalid value of volume attachment %s is provided", constants.VOLUME_ID)
		log.Errorln(msg)
		return errors.New(msg)
	} else {
		return nil
	}
}

func ValidateGetVolumeAttachmentRequest(volumeAttachment *model.VolumeAttachment) error {
	log.Infof("ValidateGetVolumeAttachmentRequest function")
	if volumeAttachment.AttachmentID == "" {
		msg := fmt.Sprintf("volume %s is not provided", constants.ATTACHMENT_ID)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volumeAttachment.Name != "" {
		msg := fmt.Sprintf("attachment %s is not supported", constants.ATTACHMENT_NAME)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volumeAttachment.VolumeID != "" {
		msg := fmt.Sprintf("%s is not supported", constants.VOLUME_ID)
		log.Errorln(msg)
		return errors.New(msg)
	} else {
		return nil
	}
}

func ValidateListVolumeAttachmentRequest(volumeAttachment *model.VolumeAttachment) error {
	log.Infof("ValidateGetVolumeAttachmentRequest function")
	if volumeAttachment.AttachmentID != "" {
		err := fmt.Errorf("volume %s is not supported", constants.ATTACHMENT_ID)
		return err
	} else if volumeAttachment.Name != "" {
		err := fmt.Errorf("volume attachment %s is not supported", constants.ATTACHMENT_NAME)
		return err
	} else if volumeAttachment.VolumeID != "" {
		err := fmt.Errorf("%s is not supported", constants.VOLUME_ID)
		return err
	} else {
		return nil
	}
}

func ValidateVolumeAttachmentRequest(operation string, volumeAttachment *model.VolumeAttachment) error {
	log.Infof("ValidateVolumeAttachmentRequest function")
	opMap := map[string]func(volumeAttachment *model.VolumeAttachment) error{
		constants.CREATE: ValidateCreateVolumeAttachmentRequest,
		constants.DELETE: ValidateGetVolumeAttachmentRequest,
		constants.GET:    ValidateGetVolumeAttachmentRequest,
		constants.LIST:   ValidateListVolumeAttachmentRequest,
	}
	return opMap[operation](volumeAttachment)
}
