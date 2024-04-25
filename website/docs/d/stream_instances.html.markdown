---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: stream instances"
sidebar_current: "docs-mongodbatlas-datasource-stream-instances"
description: |-
    Describes stream instances of a project.
---

# Data Source: mongodbatlas_stream_instances

`mongodbatlas_stream_instances` describes the stream instances defined in a project.

## Example Usage

```terraform
data "mongodbatlas_stream_instances" "test" {
    project_id = "<PROJECT_ID>"
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.

* `page_num` - (Optional) Number of the page that displays the current set of the total objects that the response returns. Defaults to `1`.
* `items_per_page` - (Optional) Number of items that the response returns per page, up to a maximum of `500`. Defaults to `100`.


## Attributes Reference

In addition to all arguments above, it also exports the following attributes:

* `results` - A list where each element contains a Stream Instance.
* `total_count` - Count of the total number of items in the result set. The count might be greater than the number of objects in the results array if the entire result set is paginated.

### Stream Instance

* `project_id` - Unique 24-hexadecimal digit string that identifies your project.
* `instance_name` - Human-readable label that identifies the stream instance.
* `data_process_region` - Defines the cloud service provider and region where MongoDB Cloud performs stream processing. See [data process region](#data-process-region).
* `hostnames` - List that contains the hostnames assigned to the stream instance.
* `stream_config` - Defines the configuration options for an Atlas Stream Processing Instance. See [stream config](#stream-config)

### Data Process Region

* `cloud_provider` - Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
* `region` - Name of the cloud provider region hosting Atlas Stream Processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
  
### Stream Config

* `tier` - Selected tier for the Stream Instance. Configures Memory / VCPU allowances. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.

To learn more, see: [MongoDB Atlas API - Stream Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) Documentation.
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_stream_instance/atlas-streams-user-journey.md) also contains details on the overall support for Atlas Streams Processing in Terraform.
