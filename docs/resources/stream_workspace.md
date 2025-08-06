# Resource: mongodbatlas_stream_workspace

`mongodbatlas_stream_workspace` provides a Stream Workspace resource. The resource lets you create, edit, and delete stream workspaces in a project.

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

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `workspace_name` - (Required) Human-readable label that identifies the stream workspace.
* `data_process_region` - (Required) Cloud service provider and region where MongoDB Cloud performs stream processing. See [data process region](#data-process-region).
* `stream_config` - (Optional) Configuration options for an Atlas Stream Processing Workspace. See [stream config](#stream-config)


### Data Process Region

* `cloud_provider` - (Required) Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
* `region` - (Required) Name of the cloud provider region hosting Atlas Stream Processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.

### Stream Config

* `tier` - (Required) Selected tier for the Stream Workspace. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
* `defaultTier` - Selected defaultTier for the Stream Workspace. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
* `maxTierSize` - Selected maxTierSize for the Stream Workspace. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hostnames` - List that contains the hostnames assigned to the stream workspace.

## Import

You can import stream workspace resource using the project ID and workspace name, in the format `PROJECT_ID-INSTANCE_NAME`. For example:

```
$ terraform import mongodbatlas_stream_workspace.test 650972848269185c55f40ca1-WorkspaceName
```

To learn more, see: [MongoDB Atlas API - Stream Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) Documentation. 
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_stream_instance/atlas-streams-user-journey.md) also contains details on the overall support for Atlas Streams Processing in Terraform.
