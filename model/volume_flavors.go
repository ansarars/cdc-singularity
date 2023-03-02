// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	log "github.com/hpe-storage/common-host-libs/logger"
)

const (
	operationFlavorList = "list"
)

type VolumeFlavor struct {
	// Volume flavor unique ID
	ID string `json:"ID,omitempty"`
	// Typical user-visible name for a volume flavor
	Name string `json:"Name,omitempty"`
}

var supportedListFlavorArgs = []string{"format", "username", "password"}

func MakeVolumeFlavor(operationType string, args map[string]interface{}) (*VolumeFlavor, error) {
	log.Infof("MakeVolumeFlavor args %+v", args)
	var requiredArgs []string
	if operationType == constants.LIST {
		requiredArgs = supportedListFlavorArgs
	}
	err := ValidateArguments(args, requiredArgs)
	if err != nil {
		msg := fmt.Sprintf("%v volume flavor failed with error: %v", operationType, err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	volumeFlavor := &VolumeFlavor{}
	jsonString, err := json.Marshal(args)
	if err != nil {
		msg := fmt.Sprintf("Volume flavor marshalling failed with error: %+v", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	// convert json to struct
	if err := json.Unmarshal(jsonString, volumeFlavor); err != nil {
		msg := fmt.Sprintf("Volume flavor unmarshalling failed with error: %+v", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	return volumeFlavor, nil
}

func (volFlavor *VolumeFlavor) CovertToTable(operationType string) *DisplayContent {
	log.Infof("CovertToTable function\n")
	display := &DisplayContent{}
	if operationType == operationFlavorList {
		display.Init(constants.LIST_FLAVOR_TABLE_COLUMNS, constants.TABLE_ROWS)
		display.Rows[0] = make([]string, constants.LIST_FLAVOR_TABLE_COLUMNS)
		display.Header[0] = "NAME"
		display.Rows[0][0] = volFlavor.Name
		display.Header[1] = "ID"
		display.Rows[0][1] = volFlavor.ID
	}
	return display
}
