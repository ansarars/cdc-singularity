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

func ListCapacityPool(glmUsername string, glmPassword string, sessionToken string,
	membershipID string, glmUrl string) (*[]model.CapacityPools, error) {
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
	statusCode, responseBody, err := client.ExecuteRestRequest(constants.GET_CAPACITYPOOLS,
		constants.REST_CAPACITYPOOLS_URL, path, rawQuery, userInfo, glmUrl, "", 0, nil,
		header, sessionToken)
	if err != nil {
		msg := fmt.Sprintf("list capacitypools failed with error: %+v\n", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("responseBody %+v\n", string(responseBody))
	if statusCode != restclient.StatusCodeOk {
		msg := fmt.Sprintf("list capacitypools failed with status code: %+v\n", statusCode)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}

	var capacityPoolsListResp []model.CapacityPools
	err = json.Unmarshal(responseBody, &capacityPoolsListResp)
	if err != nil {
		msg := fmt.Sprintf("capacitypools list Unmarshalling failed with error %v\n", err)
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	log.Infof("List CapacityPools response %+v", capacityPoolsListResp)
	return &capacityPoolsListResp, nil
}
