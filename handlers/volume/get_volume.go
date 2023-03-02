// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package volume

import (
	"errors"
	"fmt"
	client "github.com/hpe-hcss/lh-cdc-singularity/client"
	constants "github.com/hpe-hcss/lh-cdc-singularity/constants"
	capacity_pool "github.com/hpe-hcss/lh-cdc-singularity/handlers/capacity_pool"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
	"io/ioutil"
	"strings"
)

type GetVolumeHandler struct{}

func (ch *GetVolumeHandler) MakeResource(operation string, argsMap map[string]interface{}) (*model.Volume, error) {
	log.Infof("Get volume args:%v\n", argsMap)
	volume, err := model.NewVolume(operation, argsMap)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return volume, nil
}

func (ch *GetVolumeHandler) ValidateResource(operation string, volume *model.Volume) error {
	log.Infof("Validate get volume args:%v\n", volume)
	err := ValidateVolumeRequest(operation, volume)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (ch *GetVolumeHandler) Execute(volume *model.Volume, cli client.ClientInterface) (interface{}, error) {
	log.Infof("get volume:%v\n", volume.VolumeID)
	resp, err := cli.GetVolume(volume.VolumeID) // model.Volume
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infof("get volume response:%+v", resp)
	return resp, nil
}

func GetMountPath() (string, error) {
	cfg, err := ioutil.ReadFile(constants.MAPR_FUSE_CONF_FILE)
	if err != nil {
		log.Errorf("File %v read error", constants.MAPR_FUSE_CONF_FILE)
	}
	fileContent := string(cfg)
	for _, item := range strings.Split(fileContent, "\n") {
		if strings.Contains(item, "fuse.mount.point") {
			return strings.Split(item, "=")[1], nil
		}
	}
	return string(cfg), nil
}

func GetMountedVolumeName(volume model.Volume) string {
	volumeID := strings.ReplaceAll(volume.VolumeID, "-", "")
	mountedVolPath := "/cdc-vol-AV." + volumeID[0:len(volumeID)-4]
	return mountedVolPath
}

func GetVolumeResponseWithMountPath(resp *[]model.Volume, glmUsername string, glmPassword string, glmSessionToken string,
	membershipID string, glmURL string) (*[]model.Volume, error) {
	capacityPools, err := capacity_pool.ListCapacityPool(glmUsername, glmPassword, glmSessionToken, membershipID, glmURL) // model.Volume
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	mountPath := ""
	volumesList := []model.Volume{}
	for _, volume := range *resp {
		for _, capacityPool := range *capacityPools {
			capacityPoolFound := false
			for _, flavorId := range capacityPool.VolumeFlavors {
				if volume.FlavorID == flavorId {
					resp, err := capacity_pool.GetCapacityPool(glmUsername, glmPassword, glmSessionToken, membershipID,
						glmURL, capacityPool.ID) // model.Volume
					if err != nil {
						msg := fmt.Sprintf("get capacitypool failed with error: %v", err)
						log.Errorf(msg)
						return nil, errors.New(msg)
					}
					log.Infof("CapacityPool cluster name: %v", resp.ClusterName)
					mountedVolName := GetMountedVolumeName(volume)
					if volume.State == "visible" {
						if mountPath == "" {
							mountPath, err = GetMountPath()
							if err != nil {
								msg := fmt.Sprintf("get volume mount path failed with error: %v", err)
								log.Errorf(msg)
								return nil, errors.New(msg)
							}
						}
						volume.MountPath = mountPath + "/" + resp.ClusterName + mountedVolName
					}
					log.Infof("Volume structure after getting path aaa: %v", volume)
					volumesList = append(volumesList, volume)
					capacityPoolFound = true
					break
				}
			}
			if capacityPoolFound {
				break
			}
		}
	}
	log.Infof("list volume resp %+v", volumesList)
	return &volumesList, nil
}
