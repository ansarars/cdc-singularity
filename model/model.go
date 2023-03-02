// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package model

type CreateVolumeRequest struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}
