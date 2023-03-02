// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package ss

import (
	"github.com/hpe-hcss/lh-cdc-singularity/client"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type ListVolumeAttachmentHandler struct{}

func (ch *ListVolumeAttachmentHandler) MakeResource(operation string,
	argsMap map[string]interface{}) (*model.VolumeAttachment, error) {
	log.Infof("list volume attachment args:%v\n", argsMap)
	volumeAttachment, err := model.MakeVolumeAttachment(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return volumeAttachment, nil
}

func (ch *ListVolumeAttachmentHandler) ValidateResource(operation string,
	volumeAttachment *model.VolumeAttachment) error {
	log.Infof("Validate list volume attachments args:%v\n", volumeAttachment)
	err := ValidateVolumeAttachmentRequest(operation, volumeAttachment)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (ch *ListVolumeAttachmentHandler) Execute(volumeAttachment *model.VolumeAttachment,
	cli client.ClientInterface) (interface{}, error) {
	log.Infof("list volume attachments :%v\n", volumeAttachment.AttachmentID)
	resp, err := cli.ListVolumeAttachments() // model.Volume
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("list volume attachments response:%+v", resp)
	return resp, nil
}
