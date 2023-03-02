// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package volume

import (
	client "github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type DeleteVolumeHandler struct{}

func (ch *DeleteVolumeHandler) MakeResource(operation string, argsMap map[string]interface{}) (*model.Volume, error) {
	log.Infof("Delete volume args:%v\n", argsMap)
	volume, err := model.NewVolume(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return volume, nil
}

func (ch *DeleteVolumeHandler) ValidateResource(operation string, volume *model.Volume) error {
	log.Infof("Validate delete volume args:%v\n", volume)
	err := ValidateVolumeRequest(operation, volume)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (ch *DeleteVolumeHandler) Execute(volume *model.Volume, cli client.ClientInterface) (interface{}, error) {
	log.Infof("delete volume:%v\n", volume.VolumeID)
	resp, err := cli.DeleteVolume(volume.VolumeID) // model.Volume
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("delete volume response:%+v", resp)
	return resp, nil
}
