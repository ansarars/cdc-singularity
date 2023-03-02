// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package utils

import (
	"errors"
	"fmt"
	constants "github.com/hpe-hcss/lh-cdc-singularity/constants"
	log "github.com/hpe-storage/common-host-libs/logger"
	"strings"
)

func GetCredentials(argsMap map[string]interface{}) (string, string, error) {
	glmCredentialsList := [2]string{constants.PASSWORD, constants.USERNAME}
	for _, key := range glmCredentialsList {
		if val, ok := argsMap[key]; ok {
			if val == "" {
				err := fmt.Errorf("argument '%s' value is missing. For usage, execute 'singularity volume' "+
					"command", key)
				return "", "", err
			}
		} else {
			err := fmt.Errorf("argument '%s' is missing. For usage, execute 'singularity volume' command", key)
			return "", "", err
		}
	}
	password := fmt.Sprint(argsMap[constants.PASSWORD])
	username := fmt.Sprint(argsMap[constants.USERNAME])
	return username, password, nil
}

func MakeCommand(args []string) (map[string]interface{}, error) {
	argsMap := map[string]interface{}{}
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		tokens := strings.SplitN(arg, "=", 2)
		if len(tokens) < 2 {
			msg := fmt.Sprintf("value for key %s is not provided", tokens[constants.ARG_KEY_INDEX])
			log.Errorf(msg)
			return nil, errors.New(msg)
		}
		argsMap[tokens[constants.ARG_KEY_INDEX]] = tokens[constants.ARG_VALUE_INDEX]
	}
	return argsMap, nil
}

func ValidateOperations(operation string, supportedOperations []string) error {
	found := false
	for _, value := range supportedOperations {
		if operation == value {
			found = true
			break
		}
	}
	if !found {
		msg := fmt.Sprintf("Invalid operation %s, valid operations are %v", operation, supportedOperations)
		return errors.New(msg)
	}
	return nil
}
