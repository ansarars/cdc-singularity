// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package capacity_pool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	"github.com/hpe-hcss/lh-cdc-singularity/restclient"
	log "github.com/hpe-storage/common-host-libs/logger"
)

func GetCapacityPool(glmUsername string, glmPassword string, sessionToken string,
	membershipID string, glmUrl string, capacityPoolID string) (*model.CapacityPools, error) {
	path := ""
	rawQuery := ""
	client := restclient.NewRestClient()
	userInfo := restclient.UserInfo{
		UserName: glmUsername,
		UserPwd:  glmPassword,
	}

	header := map[string]string{
		"Membership": membershipID,
	}
	getCapacityPoolsUrl := constants.REST_CAPACITYPOOLS_URL + "/" + capacityPoolID
	statusCode, responseBody, err := client.ExecuteRestRequest(constants.GET_CAPACITYPOOLS,
		getCapacityPoolsUrl, path, rawQuery, userInfo, glmUrl, "", 0, nil,
		header, sessionToken)
	if err != nil {
		msg := fmt.Sprintf("get capacitypools failed with error: %+v\n", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("responseBody %+v\n", string(responseBody))
	if statusCode != restclient.StatusCodeOk {
		msg := fmt.Sprintf("get capacitypools failed with status code: %+v\n", statusCode)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	var capacityPoolsResp model.CapacityPools
	err = json.Unmarshal(responseBody, &capacityPoolsResp)
	if err != nil {
		msg := fmt.Sprintf("capacitypools get Unmarshalling failed with error %v\n", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("Get CapacityPool response %+v", capacityPoolsResp)
	return &capacityPoolsResp, nil
}
