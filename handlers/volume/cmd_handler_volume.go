// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package volume

import (
	"errors"
	"fmt"
	client "github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	"github.com/hpe-hcss/lh-cdc-singularity/utils"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type VolumeHandler interface {
	MakeResource(string, map[string]interface{}) (*model.Volume, error)
	ValidateResource(string, *model.Volume) error
	Execute(*model.Volume, client.ClientInterface) (interface{}, error)
}

type CmdHandlerVolume struct {
	args  []string
	opMap map[string]VolumeHandler
}

var supportedVolumeOperations = []string{"create", "get", "delete", "list"}

func NewCmdHandlerVolume(args []string) *CmdHandlerVolume {
	log.Infof("NewCmdHandlerVolume : %v", args)
	ch := &CmdHandlerVolume{}
	opMap := map[string]VolumeHandler{
		constants.CREATE: &CreateVolumeHandler{},
		constants.DELETE: &DeleteVolumeHandler{},
		constants.GET:    &GetVolumeHandler{},
		constants.LIST:   &ListVolumeHandler{},
	}
	ch.args = args
	ch.opMap = opMap
	return ch
}

func (ch *CmdHandlerVolume) Handle(glmCredDetails map[string]string,
	argsMap map[string]interface{}) (interface{}, error) {
	log.Infof("Handle function")
	var resp interface{}
	operation := ch.args[0]
	if err := utils.ValidateOperations(operation, supportedVolumeOperations); err != nil {
		log.Errorln(err)
		return nil, err
	}

	handler := ch.opMap[operation]
	if handler != nil {
		//validating parameters
		volume, err := handler.MakeResource(operation, argsMap)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}

		err = handler.ValidateResource(operation, volume)
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
		resp, err = handler.Execute(volume, cli)
		if err == model.TokenError || err == client.UndefinedResponseError || err == client.TokenExpiredError {
			log.Errorf("Execute err %+v", err)
			err := cli.Login()
			if err != nil {
				msg := fmt.Sprintf("Session creation failed with error: %v", err)
				log.Errorf(msg)
				return nil, errors.New(msg)
			}
			defer cli.Logout()
			resp, err = handler.Execute(volume, cli)
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

func ValidateCreateVolumeRequest(volume *model.Volume) error {
	log.Infof("ValidateCreateVolumeRequest function")
	if volume.Name == "" {
		msg := fmt.Sprintf("invalid value of volume %s is provided", constants.VOLUME_NAME)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volume.FlavorName == "" {
		msg := fmt.Sprintf("invalid value of %s is provided", constants.FLAVORNAME)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volume.Capacity == 0 {
		msg := fmt.Sprintf("invalid value of %s is provided", constants.VOLUME_CAPACITY)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volume.LocationID == "" {
		msg := fmt.Sprintf("invalid value of %s is provided", constants.LOCATION_ID)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volume.Description == "" {
		msg := fmt.Sprintf("%s is not provided", constants.DESCRIPTION)
		log.Errorln(msg)
		return errors.New(msg)
	} else {
		return nil
	}
}

func ValidateCommonVolParams(volume *model.Volume) error {
	if volume.FlavorName != "" {
		msg := fmt.Sprintf("%s is not supported", constants.FLAVORNAME)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volume.Capacity != 0 {
		msg := fmt.Sprintf("%s is not supported", constants.VOLUME_CAPACITY)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volume.LocationID != "" {
		msg := fmt.Sprintf("%s is not supported", constants.LOCATION_ID)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volume.Description != "" {
		msg := fmt.Sprintf("%s is not supported", constants.DESCRIPTION)
		log.Errorln(msg)
		return errors.New(msg)
	} else if volume.Name != "" {
		msg := fmt.Sprintf("volume %s is not supported", constants.VOLUME_NAME)
		log.Errorln(msg)
		return errors.New(msg)
	} else {
		return nil
	}
}

func ValidateGetVolumeRequest(volume *model.Volume) error {
	log.Infof("ValidateGetVolumeRequest function")
	if volume.VolumeID == "" {
		log.Errorln("volume ID is not provided")
		return errors.New("volume ID is not provided")
	}
	return ValidateCommonVolParams(volume)
}

func ValidateListVolumeRequest(volume *model.Volume) error {
	log.Infof("ValidateListVolumeRequest function")
	if volume.VolumeID != "" {
		msg := fmt.Sprintf("%s is not required", constants.VOLUME_ID)
		log.Errorln(msg)
		return errors.New(msg)
	}
	return ValidateCommonVolParams(volume)
}

func ValidateVolumeRequest(operation string, volume *model.Volume) error {
	log.Infof("ValidateVolumeRequest function")
	opMap := map[string]func(volume *model.Volume) error{
		constants.CREATE: ValidateCreateVolumeRequest,
		constants.DELETE: ValidateGetVolumeRequest,
		constants.GET:    ValidateGetVolumeRequest,
		constants.LIST:   ValidateListVolumeRequest,
	}
	return opMap[operation](volume)
}
