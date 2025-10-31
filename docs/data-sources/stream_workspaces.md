---
subcategory: "Streams"
---

# Data Source: mongodbatlas_stream_workspaces

`mongodbatlas_stream_workspaces` describes the stream workspaces defined in a project.

~> **NOTE:** This data source is an alias for `mongodbatlas_stream_instances`. Use this data source for new configurations.

## Example Usage

```terraform
data "mongodbatlas_stream_workspaces" "test" {
    project_id = "<PROJECT_ID>"
}
```

## Migration from stream_instances

To migrate from `mongodbatlas_stream_instances` data source to `mongodbatlas_stream_workspaces`, use the following `moved` block:

```terraform
moved {
  from = data.mongodbatlas_stream_instances.example
  to   = data.mongodbatlas_stream_workspaces.example
}

data "mongodbatlas_stream_workspaces" "example" {
  project_id = "<PROJECT_ID>"
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.

* `page_num` - (Optional) Number of the page that displays the current set of the total objects that the response returns. Defaults to `1`.
* `items_per_page` - (Optional) Number of items that the response returns per page, up to a maximum of `500`. Defaults to `100`.


## Attributes Reference

In addition to all arguments above, it also exports the following attributes:

* `results` - A list where each element contains a Stream Workspace.
* `total_count` - Count of the total number of items in the result set. The count might be greater than the number of objects in the results array if the entire result set is paginated.

### Stream Workspace

* `project_id` - Unique 24-hexadecimal digit string that identifies your project.
* `workspace_name` - Label that identifies the stream workspace.
* `data_process_region` - Defines the cloud service provider and region where MongoDB Cloud performs stream processing. See [data process region](#data-process-region).
* `hostnames` - List that contains the hostnames assigned to the stream workspace.
* `stream_config` - Defines the configuration options for an Atlas Stream Processing Instance. See [stream config](#stream-config)

### Data Process Region

* `cloud_provider` - Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.
* `region` - Name of the cloud provider region hosting Atlas Stream Processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.
  
### Stream Config

* `tier` - Selected tier for the Stream Workspace. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.

To learn more, see: [MongoDB Atlas API - Stream Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) Documentation.
