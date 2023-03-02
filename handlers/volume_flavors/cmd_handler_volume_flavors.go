// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package volume_flavors

import (
	"errors"
	"fmt"
	"github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	"github.com/hpe-hcss/lh-cdc-singularity/utils"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type VolumeFlavorHandler interface {
	MakeResource(string, map[string]interface{}) (*model.VolumeFlavor, error)
	ValidateResource(string, *model.VolumeFlavor) error
	Execute(*model.VolumeFlavor, client.ClientInterface) (interface{}, error)
}

type CmdHandlerVolumeFlavor struct {
	args  []string
	opMap map[string]VolumeFlavorHandler
}

var supportedVolFlavorOperations = []string{"list"}

func NewCmdHandlerVolumeFlavor(args []string) *CmdHandlerVolumeFlavor {
	log.Infof("NewCmdHandlerVolumeFlavor : %v", args)
	return &CmdHandlerVolumeFlavor{
		opMap: map[string]VolumeFlavorHandler{
			"list": &ListVolumeFlavorHandler{},
		},
		args: args,
	}
}

func (ch *CmdHandlerVolumeFlavor) Handle(glmCredDetails map[string]string,
	argsMap map[string]interface{}) (interface{}, error) {
	log.Infof("Handle function")
	var resp interface{}
	operation := ch.args[0]
	if err := utils.ValidateOperations(operation, supportedVolFlavorOperations); err != nil {
		log.Errorln(err)
		return nil, err
	}

	handler := ch.opMap[operation]
	if handler == nil {
		msg := fmt.Sprintf("unsupported sub-command: %s", operation)
		log.Errorln(msg)
		return nil, errors.New(msg)
	}
	//validating parameters
	volumeFlavor, err := handler.MakeResource(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	if err := handler.ValidateResource(operation, volumeFlavor); err != nil {
		log.Errorln(err)
		return nil, err
	}

	glmUserName, glmPassword, err := utils.GetCredentials(argsMap)
	if err != nil {
		return nil, err
	}
	cli := client.NewClient(glmCredDetails[constants.GLM_PORTAL], glmUserName,
		glmPassword, glmCredDetails[constants.MEMBERSHIP_ID])
	resp, err = handler.Execute(volumeFlavor, cli)
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
		resp, err = handler.Execute(volumeFlavor, cli)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}
	} else if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return resp, nil
}

func validateListVolumeFlavorRequest(volumeFlavor *model.VolumeFlavor) error {
	log.Infof("ValidateListVolumeFlavorRequest function")
	if volumeFlavor.ID != "" || volumeFlavor.Name != "" {
		msg := fmt.Sprintf("%s and/or %s is not supported", constants.FLAVOR_ID, constants.FLAVOR_NAME)
		log.Errorln(msg)
		return errors.New(msg)
	} else {
		return nil
	}
}

func validateVolumeFlavorRequest(operation string, volumeFlavor *model.VolumeFlavor) error {
	log.Infof("ValidateVolumeAttachmentRequest function")
	opMap := map[string]func(volumeFlavor *model.VolumeFlavor) error{
		constants.LIST: validateListVolumeFlavorRequest,
	}
	return opMap[operation](volumeFlavor)
}
