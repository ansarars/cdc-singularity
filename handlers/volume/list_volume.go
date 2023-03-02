// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package volume

import (
	client "github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type ListVolumeHandler struct{}

func (ch *ListVolumeHandler) MakeResource(operation string, argsMap map[string]interface{}) (*model.Volume, error) {
	log.Infof("List volume args:%v\n", argsMap)
	volume, err := model.NewVolume(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return volume, nil
}

func (ch *ListVolumeHandler) ValidateResource(operation string, volume *model.Volume) error {
	log.Infof("Validate list volume args:%v\n", volume)
	err := ValidateVolumeRequest(operation, volume)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (ch *ListVolumeHandler) Execute(volume *model.Volume, cli client.ClientInterface) (interface{}, error) {
	log.Infof("list volume:%v\n", volume)
	resp, err := cli.ListVolumes()
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	sessionToken, err := model.GetSessionToken()
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	glmCredentials, err := cli.GetGlmCredentials()
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	newResp, err := GetVolumeResponseWithMountPath(resp, glmCredentials["USER_NAME"], glmCredentials["PASSWORD"],
		sessionToken, glmCredentials["MEMBERSHIP_ID"], glmCredentials["URL"])
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("list volume response:%+v", newResp)
	return newResp, nil
}
