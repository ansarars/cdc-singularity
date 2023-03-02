Singularity plugin
==================

This directory contains an example CLI plugin for singularity. It
demonstrates how to add a command and flags.

Building
--------

In order to build the plugin you need to check if singularity is installed
on the setup:

    $ singularity version
    3.9.5

If singularity "3.9.5" is not installed then follow steps given in URL 
`https://sylabs.io/guides/latest/user-guide/quick_start.html`

Create config file `/etc/hpe-data-fabric/singularity/plugin.conf`. Example
of config file content is given below:

    [glm_credentials]
    glmPortal="http://172.30.215.27:3002"
    glmUsername=aW1yYW4uYW5zYXJpQGhwZS5jb20=
    glmPassword=dGVtcF9xdWFrZV85ODc2
    membershipId=D23C0865-01C2-4401-8280-E3397CBB35B5



Obtain a copy of the source code by running:

    git clone https://github.com/hpe-hcss/lh-cdc-singularity.git
    cd df-singularity-plugin

Still from within that directory, run:

	singularity plugin compile .

This will produce a file `df-singularity-plugin.sif`.

Installing
----------

Once you have compiled the plugin into a SIF file, you can install it
into the correct singularity directory using the command:

	$ singularity plugin install df-singularity-plugin.sif

Singularity will automatically load the plugin code from now on.

Other commands
--------------

You can query the list of installed plugins:

    $ singularity plugin list
    ENABLED  NAME
        yes  hpe-gl-singularity-plugin

Disable an installed plugin:

    $ singularity plugin disable hpe-gl-singularity-plugin

Enable a disabled plugin:

    $ singularity plugin enable hpe-gl-singularity-plugin

Uninstall an installed plugin:

    $ singularity plugin uninstall hpe-gl-singularity-plugin

And inspect a SIF file before installing:

    $ singularity plugin inspect df-singularity-plugin.sif
    Name: hpe-gl-singularity-plugin
    Description: CLI plugin to interface Singularity with Data Fabric
    Author: HPE Team
    Version: 0.0.1

Volume Usage:
------------
Following command would display the usage of the volume command:

    $ singularity volume --help

    Usage:
    singularity [global options...] volume <create|delete|get|list> <name> <capacity> <location_id> [description] <flavor_name> <volume_id> [format] <username> <password>


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

    Description:Allows life-cycle management of a volume

    Options:
    -h, --help   help for volume

    Examples:singularity volume list username=xyz password=xyz@hpe.com



volume operation examples:
-------------------------
1.Create volume command for mapr:

    Command for table output:
    singularity volume create name=volume_1 capacity=1 location_id=1ad98170-993e-4bfc-8b84-e689ea9a429b flavor_name=0344e238-5a04-4310-a7b2-a969b5c7bc03 description="my first volume" username=xyz@hpe.com password=xyz_9876
    
    response:
    NAME      ID                                    FLAVOR_ID                             CAPACITY  LOCATION_ID                           STATUS  STATE
    volume_1  f0907af6-5d60-4459-9077-7dc5fa9d97cb  0344e238-5a04-4310-a7b2-a969b5c7bc03  1048576   1ad98170-993e-4bfc-8b84-e689ea9a429b  ok      allocated

    Command for json output:
    singularity volume create name=volume_test capacity=12 location_id=1ad98170-993e-4bfc-8b84-e689ea9a429b flavor_id=b90a5f2d-de57-46b9-9b71-e9f9e4f25550 description="my second volume" format=json username=xyz@hpe.com password=xyz_9876
  
    Response:
    [{"CAPACITY":"12582912","FLAVOR_ID":"b90a5f2d-de57-46b9-9b71-e9f9e4f25550","ID":"02c5fe15-e35e-4b08-b925-62a318c00334","LOCATION_ID":"1ad98170-993e-4bfc-8b84-e689ea9a429b","NAME":"volume_test","STATE":"new","STATUS":""}]

2.Get volume by id:

    Command for table output:
    singularity volume get volume_id=cf2fa9bf-aee1-4924-97cc-c023ed91c524 username=xyz@hpe.com password=xyz_9876
    
    Response:
    NAME      ID                                    FLAVOR_ID                             CAPACITY  LOCATION_ID                           STATUS  STATE
    volume_1  cf2fa9bf-aee1-4924-97cc-c023ed91c524  0344e238-5a04-4310-a7b2-a969b5c7bc03  1048576   1ad98170-993e-4bfc-8b84-e689ea9a429b  ok      allocated

    Command for json output:
    singularity volume get volume_id=02c5fe15-e35e-4b08-b925-62a318c00334 format=json username=xyz@hpe.com password=xyz_9876

    Response:
    [{"CAPACITY":"12582912","FLAVOR_ID":"b90a5f2d-de57-46b9-9b71-e9f9e4f25550","ID":"02c5fe15-e35e-4b08-b925-62a318c00334","LOCATION_ID":"1ad98170-993e-4bfc-8b84-e689ea9a429b","NAME":"volume_test","STATE":"allocated","STATUS":"ok"}]

3.Delete volume by id:

    Command for table output:
    singularity volume delete volume_id=cf2fa9bf-aee1-4924-97cc-c023ed91c524 username=xyz@hpe.com password=xyz_9876
    
    Response:
    ID                                    STATE
    cf2fa9bf-aee1-4924-97cc-c023ed91c524  deleted

    Command for json output:
    singularity volume delete volume_id=02c5fe15-e35e-4b08-b925-62a318c00334 format=json username=xyz@hpe.com password=xyz_9876

    Response:
    [{"ID":"02c5fe15-e35e-4b08-b925-62a318c00334","STATE":"deleted"}]

4.List volumes:

    Command for table output:
    singularity volume list username=xyz@hpe.com password=xyz_9876

    Response:
    NAME       ID                                    MOUNT_PATH
    my_volume  ca10d15d-4d07-4ace-9205-45b7a0a1d354  /mapr/my_mapr_cluster/cdc-vol-AV.ca10d15d4d074ace920545b7a0a1
    volume3    40e31358-f9c0-44af-b376-e2e8d5e2834f

    Comamnd for json output:
    singularity volume list format=json username=xyz@hpe.com password=xyz_9876

    Response:
    [{"ID":"ca10d15d-4d07-4ace-9205-45b7a0a1d354","NAME":"my_volume", "MOUNT_PATH": "/mapr/my_mapr_cluster/cdc-vol-AV.ca10d15d4d074ace920545b7a0a1"},{"ID":"40e31358-f9c0-44af-b376-e2e8d5e2834f",
    "NAME":"volume3", "MOUNT_PATH": ""}]

Volume Attachment Usage:
------------------------
Following command would display the usage of the volume-attachment command:

    $ singularity volume-attachment --help

    Usage:
    singularity [global options...] volume-attachment <create|delete|get|list> <name> <volume_id> <attachment_id> [format] <username> <password>

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

    Description:Allows life-cycle management of a volume-attachment

    Options:
    -h, --help   help for volume-attachment

    Examples:singularity volume-attachment list username=xyz password=xyz@hpe.com


volume-attachment operation examples:
-------------------------------------
1.Create volume attachment command for mapr:

    Command for table output:
    singularity volume-attachment create name=myattachment volume_id=97b19b91-f22e-4b0a-a7bb-77a3bddf4454 username=xyz@hpe.com password=xyz_9876
    
    Response:
    NAME           ID                                    VOLUME_ID                             STATE
    myattachment2  169c43bc-bdef-4353-877a-fa54737b87f8  97b19b91-f22e-4b0a-a7bb-77a3bddf4454  new
    
    Command for json output:
    singularity volume-attachment create name=myattachment3 volume_id=72dc1b94-cc42-4ebc-b283-9857f1553736 format=json username=xyz@hpe.com password=xyz_9876

    Response:
    [{"ID":"ef3a9b0a-0dc3-451a-9ee9-9a04172ca4fa","NAME":"myattachment3","STATE":"new","VOLUME_ID":"72dc1b94-cc42-4ebc-b283-9857f1553736"}]

2.Get volume attachment by id:

    Command for table output:
    singularity volume-attachment get attachment_id=7875a96f-6582-410c-8689-047bcfc23745 username=xyz@hpe.com password=xyz_9876

    Response:
    NAME          ID                                    VOLUME_ID                             STATE  STORAGE_ID       USER_NAME  TICKET                                                                                                                                                                                                                                                                                                            TICKET_EXPIRY_TIME
    myattachment  7875a96f-6582-410c-8689-047bcfc23745  fb20da8e-4dbb-46fb-93f1-3a68ec83c70a  ready  my_mapr_cluster  mn         my_mapr_cluster pdNAjOsSTZI5xNo/vof75PD9VbiTfk/HcQoOs6jJe5kZSpxlJGWFPGczFu2bj3bigjg6juXFObFMYe+2GwCdukPjB/snPyXxCiciODfWeLQsm7RmdLSaMK74ONF0F/Iqyw5Qpj7r5C26sKL2V9PTZk0OaBhHZ6mawK+AWqMBATLgSzH0RON0sPUIpKJsB5DJJp7i071PdVxRnd02uXSaEhiJgyR3+HPyb1FnqcS3U0OTquaHJA4vhemUclFl7li+tdCv8IG1ZN3axC7le83cFuInaslz39g=  Wed Jul 13 22:34:05 UTC 2022
    
    Command for json output:
    singularity volume-attachment get attachment_id=7875a96f-6582-410c-8689-047bcfc23745 format=json username=xyz@hpe.com password=xyz_9876

    Response:
    [{"ID":"7875a96f-6582-410c-8689-047bcfc23745","NAME":"myattachment","STATE":"ready","STORAGE_ID":"my_mapr_cluster","TICKET":"my_mapr_cluster pdNAjOsSTZI5xNo/vof75PD9VbiTfk/HcQoOs6jJe5kZSpxlJGWFPGczFu2bj3bigjg6juXFObFMYe+2GwCdukPjB/snPyXxCiciODfWeLQsm7RmdLSaMK74ONF0F/Iqyw5Qpj7r5C26sKL2V9PTZk0OaBhHZ6mawK+AWqMBATLgSzH0RON0sPUIpKJsB5DJJp7i071PdVxRnd02uXSaEhiJgyR3+HPyb1FnqcS3U0OTquaHJA4vhemUclFl7li+tdCv8IG1ZN3axC7le83cFuInaslz39g=","TICKET_EXPIRY_TIME":"Wed Jul 13 22:34:05 UTC 2022","USER_NAME":"mn","VOLUME_ID":"fb20da8e-4dbb-46fb-93f1-3a68ec83c70a"}]

3.Delete volume attachment by id:
    
    Command for table output:
    singularity volume-attachment delete attachment_id=fd7f9508-11b1-49f0-b347-b1d7fbb4f038 username=xyz@hpe.com password=xyz_9876

    Response:
    ID                                    STATE
    fd7f9508-11b1-49f0-b347-b1d7fbb4f038  deleted

    Command for json output:
    singularity volume-attachment delete attachment_id=169c43bc-bdef-4353-877a-fa54737b87f8 format=json username=xyz@hpe.com password=xyz_9876

    Response:
    [{"ID":"169c43bc-bdef-4353-877a-fa54737b87f8","STATE":"deleted"}]

4.List volume attachments:

    Command for table output:
    singularity volume-attachment list username=xyz@hpe.com password=xyz_9876

    Response:
    NAME           ID                                    VOLUME_ID                             STATE
    myattachment   7875a96f-6582-410c-8689-047bcfc23745  fb20da8e-4dbb-46fb-93f1-3a68ec83c70a  ready
    myattachment2  169c43bc-bdef-4353-877a-fa54737b87f8  97b19b91-f22e-4b0a-a7bb-77a3bddf4454  ready

    Command for json output:
    singularity volume-attachment list format=json username=xyz@hpe.com password=xyz_9876

    Response:
    [{"ID":"7875a96f-6582-410c-8689-047bcfc23745","NAME":"myattachment","STATE":"ready","VOLUME_ID":
    "fb20da8e-4dbb-46fb-93f1-3a68ec83c70a"},{"ID":"169c43bc-bdef-4353-877a-fa54737b87f8","NAME":"myattachment2",
    "STATE":"ready","VOLUME_ID":"97b19b91-f22e-4b0a-a7bb-77a3bddf4454"}]

Volume Flavor Usage:
------------------------
Following command would display the usage of the volume-flavor command:

    $ singularity volume-flavor --help

    Usage:
    singularity [global options...] volume-flavor <list> <username> <password>

    Options
    - list
        Specifies volume flavor list operation.
    - username
        specifies GLM username
    - password
        specifies GLM password

    Description:Allows life-cycle management of a volume-flavor

    Options:
    -h, --help   help for volume-flavor

    Examples:singularity volume-flavor list username=xyz password=xyz@hpe.com



volume-flavor operation examples:
-------------------------------------
1.List volume flavors:

    Command for table output:
    singularity volume-flavor list username=xyz@hpe.com password=xyz_9876

    Response:
    NAME                                ID
    Default                             bce767ff-2d9e-41dd-b453-ee9b2505fc5f
    HiPerformance Filesystem Share aaa  1234e238-5a04-4310-a7b2-a969b5c7bc07

    Command for json output:
    singularity volume-flavor list format=json username=xyz@hpe.com password=xyz_9876

    Response:
    [{"ID":"bce767ff-2d9e-41dd-b453-ee9b2505fc5f","NAME":"Default"},{"ID":"1234e238-5a04-4310-a7b2-a969b5c7bc07","NAME":"HiPerformance Filesystem Share aaa"}]
