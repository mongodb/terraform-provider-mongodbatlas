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

**NOTE:**  `mongodbatlas_serverless_instance` and `mongodbatlas_privatelink_endpoint_service_serverless` resources have a circular dependency in some respects.\
That is, the `serverless_instance` must exist before the `privatelink_endpoint_service` can be created,\
and the `privatelink_endpoint_service` must exist before the `serverless_instance` gets its respective `connection_strings_private_endpoint_srv` values.

Because of this, the `serverless_instance` data source has particular value as a source of the `connection_strings_private_endpoint_srv`.\
When using the data_source in-tandem with the afforementioned resources, we can create and retrieve the `connection_strings_private_endpoint_srv` in a single `terraform apply`.

Follow this example to [setup private connection to a serverless instance using aws vpc](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/aws-privatelink-endpoint/serverless-instance) and get the connection strings in a single `terraform apply`

## Argument Reference

* `name` - (Required) Human-readable label that identifies the serverless instance.
* `project_id` - (Required) The ID of the organization or project you want to create the serverless instance within.
* `provider_settings_backing_provider_name` - (Required) Cloud service provider on which MongoDB Cloud provisioned the serverless instance.
* `provider_settings_provider_name` - (Required) Cloud service provider that applies to the provisioned the serverless instance.
* `provider_settings_region_name` - (Required) 	
  Human-readable label that identifies the physical location of your MongoDB serverless instance. The region you choose can affect network latency for clients accessing your databases.
* `continuous_backup_enabled` - (Optional) Flag that indicates whether the serverless instance uses [Serverless Continuous Backup](https://www.mongodb.com/docs/atlas/configure-serverless-backup). If this parameter is false or not used, the serverless instance uses [Basic Backup](https://www.mongodb.com/docs/atlas/configure-serverless-backup).  
* `termination_protection_enabled` - Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
* `auto_indexing` - (Optional) Flag that indicates whether the serverless instance uses [Serverless Auto Indexing](https://www.mongodb.com/docs/atlas/performance-advisor/auto-index-serverless/). This parameter defaults to true.
* `tags` - (Optional) Set that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster. See [below](#tags).

### Tags

 ```terraform
 tags {
        key   = "Key 1"
        value = "Value 1"
  }
 tags {
        key   = "Key 2"
        value = "Value 2"
  }
```

Key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.

* `key` - (Required) Constant that defines the set of the tag.
* `value` - (Required) Variable that belongs to the set of the tag.

To learn more, see [Resource Tags](https://dochub.mongodb.org/core/add-cluster-tag-atlas).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique 24-hexadecimal digit string that identifies the serverless instance.
* `connection_strings_standard_srv` - Public `mongodb+srv://` connection string that you can use to connect to this serverless instance.
* `create_date` - Timestamp that indicates when MongoDB Cloud created the serverless instance. The timestamp displays in the ISO 8601 date and time format in UTC.
* `mongo_db_version` - Version of MongoDB that the serverless instance runs, in `<major version>`.`<minor version>` format.
* `state_name` - Stage of deployment of this serverless instance when the resource made its request.
* `connection_strings_private_endpoint_srv` - List of Serverless Private Endpoint Connections

**NOTE**

## Import

Serverless Instance can be imported using the group ID and serverless instance name, in the format `GROUP_ID-SERVERLESS_INSTANCE_NAME`, e.g.

```
$ terraform import mongodbatlas_serverless_instance.my_serverless_instance 1112222b3bf99403840e8934-My Serverless Instance
```

For more information see: [MongoDB Atlas API - Serverless Instance](https://docs.atlas.mongodb.com/reference/api/serverless-instances/) Documentation.
