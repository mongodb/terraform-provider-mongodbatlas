# Data Source: mongodbatlas_stream_workspace

`mongodbatlas_stream_workspace` describes a stream workspace.

## Example Usage

```terraform
data "mongodbatlas_stream_workspace" "example" {
    project_id = "<PROJECT_ID>"
    workspace_name = "<INSTANCE_NAME>"
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `workspace_name` - (Required) Human-readable label that identifies the stream workspace.

## Attributes Reference

* `data_process_region` - Defines the cloud service provider and region where MongoDB Cloud performs stream processing. See [data process region](#data-process-region).
* `hostnames` - List that contains the hostnames assigned to the stream workspace.
* `stream_config` - Defines the configuration options for an Atlas Stream Processing Workspace. See [stream config](#stream-config)


### Data Process Region

* `cloud_provider` - Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
* `region` - Name of the cloud provider region hosting Atlas Stream Processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.

### Stream Config

* `tier` - Selected tier for the Stream Workspace. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
* `defaultTier` - Selected defaultTier for the Stream Workspace. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
* `maxTierSize` - Selected maxTierSize for the Stream Workspace. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.

To learn more, see: [MongoDB Atlas API - Stream Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) Documentation.
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_stream_instance/atlas-streams-user-journey.md) also contains details on the overall support for Atlas Streams Processing in Terraform.
