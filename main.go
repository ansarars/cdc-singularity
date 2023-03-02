// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package main

import (
	"fmt"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	handlers "github.com/hpe-hcss/lh-cdc-singularity/handlers"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	"github.com/hpe-hcss/lh-cdc-singularity/utils"
	log "github.com/hpe-storage/common-host-libs/logger"
	"github.com/spf13/cobra"
	"github.com/sylabs/singularity/pkg/cmdline"
	pluginapi "github.com/sylabs/singularity/pkg/plugin"
	clicallback "github.com/sylabs/singularity/pkg/plugin/callback/cli"
)

// Plugin is the only variable which a plugin MUST export.
// This symbol is accessed by the plugin framework to initialize the plugin.
// Plugin variable is used by plugin.so during singularity plugin compilation.
// Since go-lint is throwing error message "Plugin is unused" hence ignoring
// this error in go-lint by using below statement

//nolint:all
var Plugin = pluginapi.Plugin{
	Manifest: pluginapi.Manifest{
		Name:        constants.PLUGIN_NAME,
		Author:      constants.AUTHOR,
		Version:     constants.PLUGIN_VERSION,
		Description: constants.PLUGIN_DESCRIPTION,
	},
	Callbacks: []pluginapi.Callback{
		(clicallback.Command)(callbackVolumeCmd),
		(clicallback.Command)(callbackVolumeAttachmentCmd),
		(clicallback.Command)(callbackVolumeFlavorCmd),
	},
}

func callbackVolumeCmd(manager *cmdline.CommandManager) {
	manager.RegisterCmd(&cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Use:                   volumeUsage,
		Short:                 "volume",
		Long:                  "Allows life-cycle management of a volume",
		Example:               "singularity volume list username=xyz password=xyz@hpe.com",
		Run:                   run,
		TraverseChildren:      true,
	})
}

func callbackVolumeAttachmentCmd(manager *cmdline.CommandManager) {
	manager.RegisterCmd(&cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Use:                   volumeAttachmentUsage,
		Short:                 "volume-attachment",
		Long:                  "Allows life-cycle management of a volume-attachment",
		Example:               "singularity volume-attachment list username=xyz password=xyz@hpe.com",
		Run:                   run,
		TraverseChildren:      true,
	})
}

func callbackVolumeFlavorCmd(manager *cmdline.CommandManager) {
	manager.RegisterCmd(&cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Use:                   volumeFlavorUsage,
		Short:                 "volume-flavor",
		Long:                  "Allows life-cycle management of a volume-flavor",
		Example:               "singularity volume-flavor list username=xyz password=xyz@hpe.com",
		Run:                   run,
		TraverseChildren:      true,
	})
}

func run(cmd *cobra.Command, args []string) {
	resourceType := cmd.Short
	if len(args) < constants.MIN_ARGS_LENGTH {
		log.Errorln("Error: args length is zero")
		fmt.Println("Error: args length is zero")
		return
	}
	operationType := args[0]
	if err := log.InitLogging(constants.LOG_FILE, nil, false); err != nil {
		log.Errorf("Error: InitLogging: %v\n", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	log.Infof("Processing command %v %v\n", resourceType, args)
	defer log.Infof("Processed command %v %v\n", resourceType, args)
	ch := &handlers.CmdHandlerFactoryImpl{}

	cmdHandler, err := ch.Create(resourceType, args)
	if err != nil {
		log.Errorf("Error: %v\n", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	glm := &model.GLMCredDetails{}
	glmCredDetails, err := glm.GetGLMCredDetails()
	if err != nil {
		log.Errorln(err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	argsMap, err := utils.MakeCommand(args[1:])
	if err != nil {
		log.Errorln(err)
		fmt.Printf("Error: Response is %v\n", err)
		return
	}
	//handler a pointer to interface having parse,validate, execute
	resp, err := cmdHandler.Handle(glmCredDetails, argsMap)
	if err != nil {
		log.Errorf("Error: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}
	if resp == nil {
		log.Errorf("Error: %v", resp)
		fmt.Printf("Error: %v", resp)
		return
	}
	log.Infof("Response %v", resp)
	var format interface{}
	format = constants.FORMAT_TABLE
	if val, ok := argsMap[constants.FORMAT_KEY]; ok {
		format = val
	}
	var formatter utils.OutputFormatter
	formatter.SetFormatterType(format)
	err = formatter.PrintOutput(resp, resourceType, operationType)
	if err != nil {
		log.Errorln(err)
		fmt.Printf("Error: PrintOutput error is %v\n", err)
		return
	}
	log.Infof("Command '%v %v' successfully performed\n", resourceType, args)
}
