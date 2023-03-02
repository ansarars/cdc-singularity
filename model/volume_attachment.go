// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	glmClient "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	log "github.com/hpe-storage/common-host-libs/logger"
)

const (
	operationVolumeAttachmentCreate = "create"
	operationVolumeAttachmentGet    = "get"
	operationVolumeAttachmentList   = "list"
	operationVolumeAttachmentDelete = "delete"
)

type FSConfig struct {
	UserName         string
	StorageID        string
	Ticket           string
	TicketExpiryTime string
	Permissions      []string
}

type VolumeAttachment struct {
	Name         string                `json:"Name"`
	VolumeID     string                `json:"volume_id"`
	AttachmentID string                `json:"attachment_id"`
	State        glmClient.VaStateEnum `json:"State,omitempty"`
	FSConfig     *glmClient.VafsConfig `json:"FSConfig,omitempty"`
}

var supportedCreateAttachmentArgs = []string{"name", "volume_id", "format", "username", "password"}

var supportedDeleteAttachmentArgs = []string{"attachment_id", "format", "username", "password"}

var supportedGetAttachmentArgs = []string{"attachment_id", "format", "username", "password"}

var supportedListAttachmentArgs = []string{"format", "username", "password"}

func MakeVolumeAttachment(operationType string, args map[string]interface{}) (*VolumeAttachment, error) {
	var requiredArgs []string
	if operationType == constants.CREATE {
		requiredArgs = supportedCreateAttachmentArgs
	} else if operationType == constants.DELETE {
		requiredArgs = supportedDeleteAttachmentArgs
	} else if operationType == constants.GET {
		requiredArgs = supportedGetAttachmentArgs
	} else if operationType == constants.LIST {
		requiredArgs = supportedListAttachmentArgs
	}
	err := ValidateArguments(args, requiredArgs)
	if err != nil {
		msg := fmt.Sprintf("%v volume attachment failed with error: %v", operationType, err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	volumeAttachment := &VolumeAttachment{}
	jsonString, _ := json.Marshal(args)
	// convert json to struct
	err = json.Unmarshal(jsonString, volumeAttachment)
	if err != nil {
		msg := fmt.Sprintf("Volume attachment unmarshalling failed with error: %+v", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	return volumeAttachment, nil
}

func (attachment *VolumeAttachment) VolumeAttachmentCreateContent(display *DisplayContent) {
	display.Init(constants.CREATE_ATTACHMENT_TABLE_COL, constants.TABLE_ROWS)
	display.Rows[0] = make([]string, constants.CREATE_ATTACHMENT_TABLE_COL)
	display.Header[0] = "NAME"
	display.Rows[0][0] = attachment.Name
	display.Header[1] = "ID"
	display.Rows[0][1] = attachment.AttachmentID
	display.Header[2] = "VOLUME_ID"
	display.Rows[0][2] = attachment.VolumeID
	display.Header[3] = "STATE"
	display.Rows[0][3] = string(attachment.State)
}

func (attachment *VolumeAttachment) VolumeAttachmentGetContent(display *DisplayContent) {
	display.Init(constants.GET_ATTACHMENT_TABLE_COL, constants.TABLE_ROWS)
	display.Rows[0] = make([]string, constants.GET_ATTACHMENT_TABLE_COL)
	display.Header[0] = "NAME"
	display.Rows[0][0] = attachment.Name
	display.Header[1] = "ID"
	display.Rows[0][1] = attachment.AttachmentID
	display.Header[2] = "VOLUME_ID"
	display.Rows[0][2] = attachment.VolumeID
	display.Header[3] = "STATE"
	display.Rows[0][3] = string(attachment.State)
	display.Header[4] = "STORAGE_ID"
	display.Rows[0][4] = attachment.FSConfig.StorageID
	display.Header[5] = "USER_NAME"
	display.Rows[0][5] = attachment.FSConfig.UserName
	display.Header[6] = "TICKET"
	display.Rows[0][6] = attachment.FSConfig.Ticket
	display.Header[7] = "TICKET_EXPIRY_TIME"
	display.Rows[0][7] = attachment.FSConfig.TicketExpiryTime
}

func (attachment *VolumeAttachment) VolumeAttachmentListContent(display *DisplayContent) {
	display.Init(constants.LIST_ATTACHMENT_TABLE_COLUMNS, constants.TABLE_ROWS)
	display.Rows[0] = make([]string, constants.LIST_ATTACHMENT_TABLE_COLUMNS)
	display.Header[0] = "NAME"
	display.Rows[0][0] = attachment.Name
	display.Header[1] = "ID"
	display.Rows[0][1] = attachment.AttachmentID
	display.Header[2] = "VOLUME_ID"
	display.Rows[0][2] = attachment.VolumeID
	display.Header[3] = "STATE"
	display.Rows[0][3] = string(attachment.State)
}

func (attachment *VolumeAttachment) VolumeAttachmentDeleteContent(display *DisplayContent) {
	display.Init(constants.DELETE_ATTACHMENT_TABLE_COLUMNS, constants.TABLE_ROWS)
	display.Rows[0] = make([]string, constants.DELETE_ATTACHMENT_TABLE_COLUMNS)
	display.Header[0] = "ID"
	display.Rows[0][0] = attachment.AttachmentID
	display.Header[1] = "STATE"
	display.Rows[0][1] = string(attachment.State)
}

func (attachment *VolumeAttachment) CovertToTable(operationType string) *DisplayContent {
	log.Infof("CovertToTable function\n")
	display := &DisplayContent{}
	if operationType == operationVolumeAttachmentCreate {
		attachment.VolumeAttachmentCreateContent(display)
	} else if operationType == operationVolumeAttachmentGet {
		attachment.VolumeAttachmentGetContent(display)
	} else if operationType == operationVolumeAttachmentList {
		attachment.VolumeAttachmentListContent(display)
	} else if operationType == operationVolumeAttachmentDelete {
		attachment.VolumeAttachmentDeleteContent(display)
	}
	return display
}

func CreateVolumeAttachmentResponse(resp glmClient.VolumeAttachment, operationType string) *VolumeAttachment {
	log.Infof("Response %+v\n", resp)
	volAttachment := &VolumeAttachment{}
	if operationType == "create" || operationType == "get" {
		volAttachment.Name = resp.Name
		volAttachment.AttachmentID = resp.ID
		volAttachment.VolumeID = resp.VolumeID
		volAttachment.State = resp.State
		volAttachment.FSConfig = resp.FSConfig
	} else if operationType == "list" {
		volAttachment.Name = resp.Name
		volAttachment.AttachmentID = resp.ID
		volAttachment.VolumeID = resp.VolumeID
		volAttachment.State = resp.State
	} else if operationType == "delete" {
		volAttachment.VolumeID = resp.ID
		volAttachment.Name = resp.Name
		volAttachment.State = resp.State
	}
	return volAttachment
}
