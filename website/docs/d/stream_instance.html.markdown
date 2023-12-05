---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: stream instance"
sidebar_current: "docs-mongodbatlas-datasource-stream-instance"
description: |-
Describes a Stream Instance.
---

# Data Source: mongodbatlas_stream_instance

`mongodbatlas_stream_instance` describes a stream instance.

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
