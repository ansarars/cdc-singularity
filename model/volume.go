// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	glmClient "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	constants "github.com/hpe-hcss/lh-cdc-singularity/constants"
	log "github.com/hpe-storage/common-host-libs/logger"
	"strconv"
)

const (
	operationVolumeCreate = "create"
	operationVolumeGet    = "get"
	operationVolumeList   = "list"
	operationVolumeDelete = "delete"
)

type Volume struct {
	Name        string `json:"name"`
	VolumeID    string `json:"volume_id"`
	Description string `json:"description,omitempty"`
	// Adds a new volume to the project.  This object requires the LocationID and is used when a new volume is created independently from the host creation therefore requiring a specified location.
	FlavorID string `json:"flavor_id"`
	// The size of the volume in GiB
	FlavorName string `json:"flavor_name"`
	Capacity   int64  `json:"capacity"`
	// The location of the volume (and the storage array) LocationID is one of those listed by the LocationInfo array returned as part of the get /available-resources call. Any volumes must be in the same location as their attached Host.
	LocationID string                 `json:"location_id"`
	State      glmClient.VolumeState  `json:"State,omitempty"`
	Status     glmClient.VolumeStatus `json:"Status,omitempty"`
	MountPath  string                 `json:"mount_path,omitempty"`
}

var supportedCreateVolArgs = []string{"name", "capacity", "location_id", "description", "flavor_name",
	"format", "username", "password"}

var supportedDeleteVolArgs = []string{"volume_id", "format", "username", "password"}

var supportedGetVolArgs = []string{"volume_id", "format", "username", "password"}

var supportedListVolArgs = []string{"format", "username", "password"}

func ValidateArguments(args map[string]interface{}, requiredArgs []string) error {
	for key := range args {
		found := false
		for _, value := range requiredArgs {
			if key == value {
				found = true
				break
			}
		}
		if !found {
			msg := fmt.Sprintf("Invalid argument %s", key)
			return errors.New(msg)
		}
	}
	requiredArgs = RemoveString(requiredArgs, constants.FORMAT_KEY)
	for _, key := range requiredArgs {
		if _, ok := args[key]; !ok {
			msg := fmt.Sprintf("Argument '%s' is missing. For usage, execute 'singularity volume' command", key)
			return errors.New(msg)
		}
	}
	return nil
}

func NewVolume(operationType string, args map[string]interface{}) (*Volume, error) {
	log.Infof("NewVolume args %+v", args)
	var requiredArgs []string
	if operationType == constants.CREATE {
		requiredArgs = supportedCreateVolArgs
	} else if operationType == constants.DELETE {
		requiredArgs = supportedDeleteVolArgs
	} else if operationType == constants.GET {
		requiredArgs = supportedGetVolArgs
	} else if operationType == constants.LIST {
		requiredArgs = supportedListVolArgs
	}
	err := ValidateArguments(args, requiredArgs)
	if err != nil {
		msg := fmt.Sprintf("%v volume failed with error: %v", operationType, err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	volume := &Volume{}
	if val, ok := args[constants.VOLUME_CAPACITY]; ok {
		capacity := fmt.Sprintf("%v", val)
		capacityInt64, _ := strconv.ParseInt(capacity, constants.BASE, constants.BIT_SIZE)
		args[constants.VOLUME_CAPACITY] = capacityInt64
	}

	jsonString, _ := json.Marshal(args)
	// convert json to struct
	err = json.Unmarshal(jsonString, volume)
	if err != nil {
		msg := fmt.Sprintf("Volume unmarshalling failed with error: %+v", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	return volume, nil
}

type DisplayContent struct {
	Header []string
	Rows   [][]string
}

func (content *DisplayContent) Init(numColumns int, numRows int) {
	content.Header = make([]string, numColumns)
	content.Rows = make([][]string, numRows)
}

func (vol *Volume) CovertToTable(operationType string) *DisplayContent {
	log.Infof("CovertToTable function\n")
	display := &DisplayContent{}
	if operationType == operationVolumeCreate || operationType == operationVolumeGet {
		display.Init(constants.CREATE_VOL_TABLE_COLUMNS, constants.TABLE_ROWS)
		display.Rows[0] = make([]string, constants.CREATE_VOL_TABLE_COLUMNS)
		display.Header[0] = "NAME"
		display.Rows[0][0] = vol.Name
		display.Header[1] = "ID"
		display.Rows[0][1] = vol.VolumeID
		display.Header[2] = "FLAVOR_ID"
		display.Rows[0][2] = vol.FlavorID
		display.Header[3] = "CAPACITY"
		display.Rows[0][3] = strconv.FormatInt(vol.Capacity, 10)
		display.Header[4] = "STATUS"
		display.Rows[0][4] = string(vol.Status)
		display.Header[5] = "STATE"
		display.Rows[0][5] = string(vol.State)
	} else if operationType == operationVolumeList {
		display.Init(constants.LIST_TABLE_COLUMNS, constants.TABLE_ROWS)
		display.Rows[0] = make([]string, constants.LIST_TABLE_COLUMNS)
		display.Header[0] = "NAME"
		display.Rows[0][0] = vol.Name
		display.Header[1] = "ID"
		display.Rows[0][1] = vol.VolumeID
		display.Header[2] = "MOUNT_PATH"
		display.Rows[0][2] = vol.MountPath
	} else if operationType == operationVolumeDelete {
		display.Init(constants.DELETE_TABLE_COLUMNS, constants.TABLE_ROWS)
		display.Rows[0] = make([]string, constants.DELETE_TABLE_COLUMNS)
		display.Header[0] = "ID"
		display.Rows[0][0] = vol.VolumeID
		display.Header[1] = "STATE"
		display.Rows[0][1] = string(vol.State)
	}
	return display
}

func CreateResponse(resp glmClient.Volume, operationType string) *Volume {
	log.Infof("CreateResponse %+v\n", resp)
	vol := &Volume{}
	if operationType == "create" || operationType == "get" {
		vol.Name = resp.Name
		vol.VolumeID = resp.ID
		vol.FlavorID = resp.FlavorID
		vol.Capacity = resp.Capacity
		vol.LocationID = resp.LocationID
		vol.State = resp.State
		vol.Status = resp.Status
	} else if operationType == "list" {
		vol.Name = resp.Name
		vol.VolumeID = resp.ID
		vol.FlavorID = resp.FlavorID
		vol.State = resp.State
	} else if operationType == "delete" {
		vol.VolumeID = resp.ID
		vol.State = resp.State
	}
	return vol
}

func RemoveString(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
