// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	model "github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
	"github.com/rodaine/table"
	"reflect"
	"strings"
)

type OutputFormatter struct {
	formatter FormatterInterface
}

type TableOutputFormatter struct {
}

type JSONOutputFormatter struct {
}

type FormatterInterface interface {
	PrintOutput(content *model.DisplayContent) error
}

func stringToInterface(row []string) []interface{} {
	rowList := make([]interface{}, len(row))
	for i, item := range row {
		rowList[i] = reflect.ValueOf(item).Interface()
	}
	return rowList
}

func (f *TableOutputFormatter) PrintOutput(displaycontent *model.DisplayContent) error {
	if displaycontent == nil {
		return errors.New("invalid display content")
	}
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
	headerList := make([]interface{}, len(displaycontent.Header))
	for i, item := range displaycontent.Header {
		//header := reflect.ValueOf(displaycontent.Header).Interface()
		headerList[i] = reflect.ValueOf(item).Interface()
	}
	tbl := table.New(headerList...)
	rowLength := len(displaycontent.Rows)
	for i := 0; i < rowLength; i++ {
		rowList := stringToInterface(displaycontent.Rows[i])
		tbl.AddRow(rowList...)
	}
	//tbl := table.New("NAME", "ID")
	//row := reflect.ValueOf(displaycontent.Rows[0]).Interface()
	tbl.Print()
	return nil
}

func (f *JSONOutputFormatter) PrintOutput(displaycontent *model.DisplayContent) error {
	if displaycontent == nil {
		return errors.New("invalid display content")
	}
	log.Infof("json volume info %v", displaycontent)
	result := []map[string]interface{}{}
	for row := 0; row < len(displaycontent.Rows); row++ {
		mp := make(map[string]interface{})
		for column := 0; column < len(displaycontent.Header); column++ {
			mp[displaycontent.Header[column]] = displaycontent.Rows[row][column]
		}
		result = append(result, mp)
	}
	jsonStr, _ := json.Marshal(result)
	fmt.Printf("%v\n", string(jsonStr))
	return nil
}

type VolumeFormatter struct {
	operationType string
}

func NewVolumeFormatter(operationType string) *VolumeFormatter {
	return &VolumeFormatter{operationType: operationType}
}

func (v *VolumeFormatter) Format(resp interface{}) *model.DisplayContent {
	var displayContent *model.DisplayContent
	if v.operationType == constants.LIST {
		displayContent = &model.DisplayContent{}
		resources := resp.(*[]model.Volume)
		for index, item := range *resources {
			t := item.CovertToTable(constants.LIST)
			if index == 0 {
				displayContent.Header = t.Header
			}
			displayContent.Rows = append(displayContent.Rows, t.Rows[0])
		}
	} else if v.operationType == constants.GET || v.operationType == constants.CREATE {
		resource := resp.(*model.Volume)
		displayContent = resource.CovertToTable(constants.CREATE)
	} else if v.operationType == constants.DELETE {
		resource := resp.(*model.Volume)
		displayContent = resource.CovertToTable(constants.DELETE)
	}
	return displayContent
}

type VolumeAttachmentFormatter struct {
	operationType string
}

func NewVolumeAttachmentFormatter(operationType string) *VolumeAttachmentFormatter {
	return &VolumeAttachmentFormatter{operationType: operationType}
}

func (v *VolumeAttachmentFormatter) Format(resp interface{}) *model.DisplayContent {
	var displayContent *model.DisplayContent
	if v.operationType == constants.CREATE {
		resource := resp.(*model.VolumeAttachment)
		displayContent = resource.CovertToTable(constants.CREATE)
	} else if v.operationType == constants.GET {
		resource := resp.(*model.VolumeAttachment)
		displayContent = resource.CovertToTable(constants.GET)
	} else if v.operationType == constants.DELETE {
		resource := resp.(*model.VolumeAttachment)
		displayContent = resource.CovertToTable(constants.DELETE)
	} else if v.operationType == constants.LIST {
		resources := resp.(*[]model.VolumeAttachment)
		displayContent = &model.DisplayContent{}
		for index, item := range *resources {
			t := item.CovertToTable(constants.LIST)
			if index == 0 {
				displayContent.Header = t.Header
			}
			displayContent.Rows = append(displayContent.Rows, t.Rows[0])
		}
	}
	return displayContent
}

type VolumeFlavorFormatter struct {
	operationType string
}

func NewVolumeFlavorFormatter(operationType string) *VolumeFlavorFormatter {
	return &VolumeFlavorFormatter{operationType: operationType}
}

func (v *VolumeFlavorFormatter) Format(resp interface{}) *model.DisplayContent {
	var displayContent *model.DisplayContent
	if v.operationType == constants.LIST {
		displayContent = &model.DisplayContent{}
		resources := resp.(*[]model.VolumeFlavor)
		for index, item := range *resources {
			t := item.CovertToTable(constants.LIST)
			if index == 0 {
				displayContent.Header = t.Header
			}
			displayContent.Rows = append(displayContent.Rows, t.Rows[0])
		}
	}
	return displayContent
}

func (f *OutputFormatter) PrintOutput(resp interface{}, resourceType string, operationType string) error {
	var displayContent *model.DisplayContent
	if resourceType == constants.VOLUME {
		formatter := NewVolumeFormatter(operationType)
		displayContent = formatter.Format(resp)
	} else if resourceType == constants.VOLUME_ATTACHMENT {
		formatter := NewVolumeAttachmentFormatter(operationType)
		displayContent = formatter.Format(resp)
	} else if resourceType == constants.VOLUME_FLAVORS {
		formatter := NewVolumeFlavorFormatter(operationType)
		displayContent = formatter.Format(resp)
	}
	return f.formatter.PrintOutput(displayContent)
}

func (f *OutputFormatter) SetFormatterType(format interface{}) {
	if format == "table" {
		f.formatter = &TableOutputFormatter{}
	} else if format == "json" {
		f.formatter = &JSONOutputFormatter{}
	}
}
