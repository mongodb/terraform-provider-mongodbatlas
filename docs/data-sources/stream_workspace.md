---
subcategory: "Streams"
---

# Data Source: mongodbatlas_stream_workspace

`mongodbatlas_stream_workspace` describes a stream workspace that contains configurations for stream processing.

~> **NOTE:** This data source is an alias for `mongodbatlas_stream_instance`. Use this data source for new configurations.

## Example Usage

```terraform
data "mongodbatlas_stream_workspace" "example" {
    project_id = "<PROJECT_ID>"
    workspace_name = "<WORKSPACE_NAME>"
}
```

## Migration from stream_instance

To migrate from `mongodbatlas_stream_instance` data source to `mongodbatlas_stream_workspace`, use the following `moved` block:

```terraform
moved {
  from = data.mongodbatlas_stream_instance.example
  to   = data.mongodbatlas_stream_workspace.example
}

data "mongodbatlas_stream_workspace" "example" {
  project_id = "<PROJECT_ID>"
  workspace_name = "<WORKSPACE_NAME>"  # Changed from instance_name
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `workspace_name` - (Required) Label that identifies the stream workspace.

## Attributes Reference

* `data_process_region` - Defines the cloud service provider and region where MongoDB Cloud performs stream processing. See [data process region](#data-process-region).
* `hostnames` - List that contains the hostnames assigned to the stream workspace.
* `stream_config` - Defines the configuration options for an Atlas Stream Processing Instance. See [stream config](#stream-config)


### Data Process Region

* `cloud_provider` - Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.
* `region` - Name of the cloud provider region hosting Atlas Stream Processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.

### Stream Config

* `tier` - Selected tier for the Stream Instance. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.

To learn more, see: [MongoDB Atlas API - Stream Workspace](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) Documentation.
