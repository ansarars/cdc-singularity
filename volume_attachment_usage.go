// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package main

const volumeAttachmentUsage = `volume-attachment <create|delete|get|list> <name> <volume_id> <attachment_id> [format] <username> <password>

Options
- create
    Specifies volume attachment creation operation.
- delete
    Specifies volume attachment deletion operation.
- list
    Specifies volume attachment list operation.
- get 
    Specifies volume attachment get operation.
- name
    Specifies volume attachment name with type string. Required for volume attachment create operation only.
- volume_id
    Specifies volume_id to be attached with type string. Required for volume attachment create operation only.
- attachment_id
    Specifies volume attachment ID with type string. Required for volume attachment <get|delete> operations only.
- format
    specifies the format of <create|get|delete|list> response. format having two values "json" or "table". 
    if format is not mentioned in the commandline then default format value will be "table".
- username
	specifies GLM username
- password
	specifies GLM password
`
