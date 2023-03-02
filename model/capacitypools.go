// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package model

type CapacityPools struct {
	// Volume flavor unique ID
	ID string `json:"ID,omitempty"`
	// Typical user-visible name for a volume flavor
	Name          string   `json:"Name,omitempty"`
	ClusterName   string   `json:"arrayCapacityPoolID,omitempty"`
	VolumeFlavors []string `json:"volumeFlavors,omitempty"`
}
