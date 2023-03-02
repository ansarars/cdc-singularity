// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package volume

import (
	"errors"
	"fmt"
	client "github.com/hpe-hcss/lh-cdc-singularity/client"
	volumeFlavors "github.com/hpe-hcss/lh-cdc-singularity/handlers/volume_flavors"
	model "github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type CreateVolumeHandler struct{}

func (ch *CreateVolumeHandler) MakeResource(operation string, argsMap map[string]interface{}) (*model.Volume, error) {
	log.Infof("Create volume args:%v\n", argsMap)
	volume, err := model.NewVolume(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return volume, nil
}

func (ch *CreateVolumeHandler) ValidateResource(operation string, volume *model.Volume) error {
	log.Infof("Validate volume args:%v\n", volume)
	err := ValidateVolumeRequest(operation, volume)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (ch *CreateVolumeHandler) Execute(volume *model.Volume, cli client.ClientInterface) (interface{}, error) {
	log.Infof("create volume req:%v\n", volume)
	volumeFlavorList, err := cli.ListVolumeFlavors()
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	flavorID, err := volumeFlavors.GetVolumeFlavorID(volume.FlavorName, volumeFlavorList)
	if err != nil {
		msg := fmt.Sprintf("Volume flavor ID with flavor name %v not found.", volume.FlavorName)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	volume.FlavorID = flavorID
	resp, err := cli.CreateVolume(volume) // model.Volume
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("create volume response:%+v", resp)
	return resp, nil
}
