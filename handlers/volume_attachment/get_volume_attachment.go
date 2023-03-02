// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package ss

import (
	"github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type GetVolumeAttachmentHandler struct{}

func (ch *GetVolumeAttachmentHandler) MakeResource(operation string,
	argsMap map[string]interface{}) (*model.VolumeAttachment, error) {
	log.Infof("Get volume attachment args:%v\n", argsMap)
	volumeAttachment, err := model.MakeVolumeAttachment(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return volumeAttachment, nil
}

func (ch *GetVolumeAttachmentHandler) ValidateResource(operation string,
	volumeAttachment *model.VolumeAttachment) error {
	log.Infof("Validate get volume attachment args:%v\n", volumeAttachment)
	err := ValidateVolumeAttachmentRequest(operation, volumeAttachment)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (ch *GetVolumeAttachmentHandler) Execute(volumeAttachment *model.VolumeAttachment,
	cli client.ClientInterface) (interface{}, error) {
	log.Infof("get volume attachment :%v\n", volumeAttachment.AttachmentID)
	resp, err := cli.GetVolumeAttachment(volumeAttachment.AttachmentID) // model.Volume
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("get volume attachment response:%+v", resp)
	return resp, nil
}
