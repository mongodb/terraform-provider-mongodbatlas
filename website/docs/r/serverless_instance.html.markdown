---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: serverless instance"
sidebar_current: "docs-mongodbatlas-resource-serverless-instance"
description: |-
Provides a Serverless Instance resource.
---

# Resource: mongodbatlas_serverless_instance

`mongodbatlas_serverless_instance` provides a Serverless Instance resource. This allows serverless instances to be created.

> **NOTE:**  Serverless instances do not support some Atlas features at this time.
For a full list of unsupported features, see [Serverless Instance Limitations](https://docs.atlas.mongodb.com/reference/serverless-instance-limitations/).

## Example Usage

### Basic
```terraform
resource "mongodbatlas_serverless_instance" "test" {
  project_id   = "<PROJECT_ID>"
  name         = "<SERVERLESS_INSTANCE_NAME>"

  provider_settings_backing_provider_name = "AWS"
  provider_settings_provider_name = "SERVERLESS"
  provider_settings_region_name = "US_EAST_1"
}
```

## Argument Reference

* `name` - (Required) Human-readable label that identifies the serverless instance.
* `project_id` - (Required) The ID of the organization or project you want to create the serverless instance within.
* `provider_settings_backing_provider_name` - (Required) Cloud service provider on which MongoDB Cloud provisioned the serverless instance.
* `provider_settings_provider_name` - (Required) Cloud service provider that applies to the provisioned the serverless instance.
* `provider_settings_region_name` - (Required) 	
  Human-readable label that identifies the physical location of your MongoDB serverless instance. The region you choose can affect network latency for clients accessing your databases.
* `continuous_backup_enabled` - (Optional) Flag that indicates whether the serverless instance uses [Serverless Continuous Backup](https://www.mongodb.com/docs/atlas/configure-serverless-backup). If this parameter is false or not used, the serverless instance uses [Basic Backup](https://www.mongodb.com/docs/atlas/configure-serverless-backup).  
* `termination_protection_enabled` - Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique 24-hexadecimal digit string that identifies the serverless instance.
* `connection_strings_standard_srv` - Public `mongodb+srv://` connection string that you can use to connect to this serverless instance.
* `create_date` - Timestamp that indicates when MongoDB Cloud created the serverless instance. The timestamp displays in the ISO 8601 date and time format in UTC.
* `mongo_db_version` - Version of MongoDB that the serverless instance runs, in `<major version>`.`<minor version>` format.
* `state_name` - Stage of deployment of this serverless instance when the resource made its request.
* `connection_strings_private_endpoint_srv` - List of Serverless Private Endpoint Connections

## Import

Serverless Instance can be imported using the group ID and serverless instance name, in the format `GROUP_ID-SERVERLESS_INSTANCE_NAME`, e.g.

```
$ terraform import mongodbatlas_serverless_instance.my_serverless_instance 1112222b3bf99403840e8934-My Serverless Instance
```

For more information see: [MongoDB Atlas API - Serverless Instance](https://docs.atlas.mongodb.com/reference/api/serverless-instances/) Documentation.
