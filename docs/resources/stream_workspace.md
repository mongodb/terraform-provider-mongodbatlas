---
subcategory: "Streams"
---

# Resource: mongodbatlas_stream_workspace

`mongodbatlas_stream_workspace` provides a Stream Workspace resource. The resource lets you create, edit, and delete stream workspaces in a project.

~> **NOTE:** This resource is an alias for `mongodbatlas_stream_instance`. Use this resource for new configurations.

## Example Usage

```terraform
resource "mongodbatlas_stream_workspace" "test" {
    project_id = var.project_id
	workspace_name = "WorkspaceName"
	data_process_region = {
		region = "VIRGINIA_USA"
		cloud_provider = "AWS"
  }
}
```

### Further Examples
- [Atlas Stream Workspace](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_stream_workspace)

## Migration from stream_instance

To migrate from `mongodbatlas_stream_instance` to `mongodbatlas_stream_workspace`, use the following `moved` block:

```terraform
moved {
  from = mongodbatlas_stream_instance.example
  to   = mongodbatlas_stream_workspace.example
}

resource "mongodbatlas_stream_workspace" "example" {
  project_id = var.project_id
  workspace_name = "WorkspaceName"  # Changed from instance_name
  data_process_region = {
    region = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `workspace_name` - (Required) Label that identifies the stream workspace.
* `data_process_region` - (Required) Cloud service provider and region where MongoDB Cloud performs stream processing. See [data process region](#data-process-region).
* `stream_config` - (Optional) Configuration options for an Atlas Stream Processing Instance. See [stream config](#stream-config)


### Data Process Region

* `cloud_provider` - (Required) Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.
* `region` - (Required) Name of the cloud provider region hosting Atlas Stream Processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.

### Stream Config

* `max_tier_size` - (Optional) Max tier size for the Stream Workspace. Configures Memory / VCPU allowances.
* `tier` - (Optional) Selected tier for the Stream Workspace. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/creategroupstreamworkspace) describes the valid values.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hostnames` - List that contains the hostnames assigned to the stream workspace.

## Import

You can import stream workspace resource using the project ID and workspace name, in the format `PROJECT_ID-WORKSPACE_NAME`. For example:

```
$ terraform import mongodbatlas_stream_workspace.test 650972848269185c55f40ca1-WorkspaceName
```
