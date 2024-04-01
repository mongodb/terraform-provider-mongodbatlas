---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: serverless instances"
sidebar_current: "docs-mongodbatlas-datasource-serverless-instances"
description: |-
Describes a Serverless Instances.
---

# Data Source: mongodbatlas_serverless_instances

`mongodbatlas_serverless_instances` describe all serverless instances. This represents serverless instances that have been created for the specified group id.

> **NOTE:**  Serverless instances do not support some Atlas features at this time.
For a full list of unsupported features, see [Serverless Instance Limitations](https://docs.atlas.mongodb.com/reference/serverless-instance-limitations/).

> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.


## Example Usage

```terraform
data "mongodbatlas_serverless_instances" "data_serverless" {
  project_id = "<PROJECT_ID"
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the [project](https://docs.atlas.mongodb.com/organizations-projects/#std-label-projects) that contains the specified cluster.

## Attributes Reference
* `results` - A list where each result represents a Serverless Instance. To learn more, see [Serverless Instance](#serverless-instance).

### Serverless Instance

* `name` - (Required) Human-readable label that identifies your serverless instance.
* `connection_strings_standard_srv` - Public `mongodb+srv://` connection string that you can use to connect to this serverless instance.
* `created_date` - Timestamp that indicates when MongoDB Cloud created the serverless instance. The timestamp displays in the ISO 8601 date and time format in UTC.
* `id` - Unique 24-hexadecimal digit string that identifies the serverless instance.
* `mongo_db_version` - Version of MongoDB that the serverless instance runs, in `<major version>`.`<minor version>` format.
* `provider_settings_backing_provider_name` - Cloud service provider on which MongoDB Cloud provisioned the serverless instance.
* `provider_settings_provider_name` - Cloud service provider that applies to the provisioned the serverless instance.
* `provider_settings_region_name` - Human-readable label that identifies the physical location of your MongoDB serverless instance. The region you choose can affect network latency for clients accessing your databases.
* `state_name` - Stage of deployment of this serverless instance when the resource made its request.
* `continuous_backup_enabled` - Flag that indicates whether the serverless instance uses Serverless Continuous Backup.
* `termination_protection_enabled` - Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
* `auto_indexing` - Flag that indicates whether the serverless instance uses [Serverless Auto Indexing](https://www.mongodb.com/docs/atlas/performance-advisor/auto-index-serverless/).
* `tags` - Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#tags).

### Tags

Key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.

* `key` - Constant that defines the set of the tag.
* `value` - Variable that belongs to the set of the tag.

To learn more, see [Resource Tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas).



For more information see: [MongoDB Atlas API - Serverless Instance](https://docs.atlas.mongodb.com/reference/api/serverless-instances/) Documentation.
