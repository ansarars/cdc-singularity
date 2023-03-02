// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package handlers

import (
	"fmt"
	constants "github.com/hpe-hcss/lh-cdc-singularity/constants"
	volume "github.com/hpe-hcss/lh-cdc-singularity/handlers/volume"
	volumeAttachment "github.com/hpe-hcss/lh-cdc-singularity/handlers/volume_attachment"
	volumeFlavor "github.com/hpe-hcss/lh-cdc-singularity/handlers/volume_flavors"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type CmdHandlerFactory interface {
	Create(args []string) *CmdHandler
}

type CmdHandlerFactoryImpl struct {
}

type CmdHandler interface {
	// Handle -glmCredentials is a key value pairs of glm info parameters i.e. glmPortal and membershipId
	// saved in plugin.conf
	// -argsMap is a key value pairs of the arguments passed in command line.

	Handle(glmCredentials map[string]string, argsMap map[string]interface{}) (interface{}, error)
}

func (ch *CmdHandlerFactoryImpl) Create(resourceType string, args []string) (CmdHandler, error) {
	log.Infof("Create function")
	if resourceType == constants.VOLUME {
		return volume.NewCmdHandlerVolume(args), nil
	} else if resourceType == constants.VOLUME_ATTACHMENT {
		return volumeAttachment.NewCmdHandlerVolumeAttachment(args), nil
	} else if resourceType == constants.VOLUME_FLAVORS {
		return volumeFlavor.NewCmdHandlerVolumeFlavor(args), nil
	} else {
		return nil, fmt.Errorf("invalid resource type")
	}
}
