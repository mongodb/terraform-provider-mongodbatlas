---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: stream instance"
sidebar_current: "docs-mongodbatlas-datasource-stream-instance"
description: |-
    Describes a Stream Instance.
---

# Data Source: mongodbatlas_stream_instance

`mongodbatlas_stream_instance` describes a stream instance.

-> **NOTE:** Atlas Stream Processing is currently in preview, you'll need to set the environment variable `MONGODB_ATLAS_ENABLE_BETA=true` to use this data source. To learn more about stream processing and existing limitations, see the [Atlas Stream Processing Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/overview/#atlas-stream-processing-overview).

## Example Usage

```terraform
data "mongodbatlas_stream_instance" "example" {
    project_id = "<PROJECT_ID>"
    instance_name = "<INSTANCE_NAME>"
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `instance_name` - (Required) Human-readable label that identifies the stream instance.

## Attributes Reference

* `data_process_region` - Defines the cloud service provider and region where MongoDB Cloud performs stream processing. See [data process region](#data-process-region).
* `hostnames` - List that contains the hostnames assigned to the stream instance.


### Data Process Region

* `cloud_provider` - Label that identifies the cloud service provider where MongoDB Cloud performs stream processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.
* `region` - Name of the cloud provider region hosting Atlas Stream Processing. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) describes the valid values.

To learn more, see: [MongoDB Atlas API - Stream Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamInstance) Documentation.
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/atlas-streams/README.md) also contains details on the overall support for Atlas Streams Processing in Terraform.