// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package main

const volumeUsage = `volume <create|delete|get|list> <name> <capacity> <location_id> [description] <flavor_id> <volume_id> [format] <username> <password>


Options
- create
    Specifies volume creation operation.
- delete
    Specifies volume deletion operation.
- list
    Specifies volume list operation.
- get 
    Specifies volume get operation.
- name
    Specifies the name of the volume with type string. Required for <create> operation only.
- capacity
    Specifies volume capacity in GiB with type int. Required for <create> operation only.
- location_id
    Specifies volume location with type string. Required for <create> operation only.
- description
    description specifies volume description with type string. Optional for <create> operation. Not required
    for <get|delete|list> operations.
- flavor_name
    Specifies storage flavor name with type string. Required for <create> operation only.
- volume_id
    Specifies volume ID with type string. Required for <get|delete> operations only.
- format
    specifies the format of <create|get|delete|list> response. format having two values "json" or "table". 
    If format is not mentioned in the commandline then default format value will be "table".
- username
	specifies GLM username
- password
	specifies GLM password
`
