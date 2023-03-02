// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package volume_flavors

import (
	"errors"
	"github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type ListVolumeFlavorHandler struct{}

func (ch *ListVolumeFlavorHandler) MakeResource(operation string,
	argsMap map[string]interface{}) (*model.VolumeFlavor, error) {
	log.Infof("list volume flavors args:%v\n", argsMap)
	volumeFlavor, err := model.MakeVolumeFlavor(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return volumeFlavor, nil
}

func (ch *ListVolumeFlavorHandler) ValidateResource(operation string,
	volumeFlavor *model.VolumeFlavor) error {
	log.Infof("Validate list volume flavors args:%v\n", volumeFlavor)
	if err := validateVolumeFlavorRequest(operation, volumeFlavor); err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (ch *ListVolumeFlavorHandler) Execute(volumeFlavor *model.VolumeFlavor,
	cli client.ClientInterface) (interface{}, error) {
	log.Infof("list volume flavors :%v\n", volumeFlavor.ID)
	resp, err := cli.ListVolumeFlavors()
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("list volume flavors response:%+v", resp)
	return resp, nil
}

func GetVolumeFlavorID(flavorName string, volumeFlavorList *[]model.VolumeFlavor) (string, error) {
	for _, item := range *volumeFlavorList {
		if item.Name == flavorName {
			return item.ID, nil
		}
	}
	return "", errors.New("Volume flavor ID not found\n")
}
