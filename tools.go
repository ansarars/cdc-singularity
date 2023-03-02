// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

//go:build tools
// +build tools

// Package tools is used to manipulate the go package system into keeping
// packages in go.mod and go.sum that are used by the build process but are not
// included by any code in the repository. See
// https://marcofranssen.nl/manage-go-tools-via-go-modules for more information.
package tools

import (
	_ "github.com/golang/mock/mockgen"
)
