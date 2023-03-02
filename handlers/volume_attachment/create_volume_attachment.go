// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package ss

import (
	client "github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type CreateVolumeAttachmentHandler struct{}

func (ch *CreateVolumeAttachmentHandler) MakeResource(operation string,
	argsMap map[string]interface{}) (*model.VolumeAttachment, error) {
	log.Infof("Create volume attachment args:%v\n", argsMap)
	volumeAttachment, err := model.MakeVolumeAttachment(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return volumeAttachment, nil
}

func (ch *CreateVolumeAttachmentHandler) ValidateResource(operation string,
	volumeAttachment *model.VolumeAttachment) error {
	log.Infof("Validate volume attachment args:%v\n", volumeAttachment)
	err := ValidateVolumeAttachmentRequest(operation, volumeAttachment)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (ch *CreateVolumeAttachmentHandler) Execute(volumeAttachment *model.VolumeAttachment,
	cli client.ClientInterface) (interface{}, error) {
	log.Infof("create volume attachment:%v\n", volumeAttachment)
	resp, err := cli.CreateVolumeAttachment(volumeAttachment) // model.Volume
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("create volume attachment response:%+v", resp)
	return resp, nil
}
