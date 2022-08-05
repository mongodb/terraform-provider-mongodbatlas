---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: privatelink_endpoint_service_adl"
sidebar_current: "docs-mongodbatlas-resource-privatelink-endpoint-service-adl"
description: |-
Provides an Atlas Data Lake and Online Archive PrivateLink endpoint.
---


# Resource: privatelink_endpoint_service_adl

`privatelink_endpoint_service_adl` Provides an Atlas Data Lake (ADL) and Online Archive PrivateLink endpoint resource.   The same configuration will provide a PrivateLink connection for either Atlas Data Lake or Online Archive.

## Example Usage

### Basic
```terraform
resource "mongodbatlas_privatelink_endpoint_service_adl" "adl_test" {
  project_id   = "<PROJECT_ID>"
  endpoint_id  = "<ENDPOINT_ID>"
  comment      = "comments for private link endpoint adl"
  type		 = "DATA_LAKE"
  provider_name	 = "AWS"
}
```


## Argument Reference

* `project_id` - (Required) Unique 24-digit hexadecimal string that identifies the project.
* `endpoint_id` - (Required) Unique 22-character alphanumeric string that identifies the private endpoint. Atlas supports AWS private endpoints using the [|aws| PrivateLink](https://aws.amazon.com/privatelink/) feature.
* `type` - (Required) Human-readable label that identifies the type of resource to associate with this private endpoint. Atlas supports `DATA_LAKE` only. If empty, defaults to `DATA_LAKE`.
* `provider_name` - (Required) Human-readable label that identifies the cloud provider for this endpoint. Atlas supports AWS only. If empty, defaults to AWS.
* `comment` - Human-readable string to associate with this private endpoint.

## Import

ADL privatelink endpoint can be imported using project ID and endpoint ID, in the format `project_id`--`endpoint_id`, e.g.

```
$ terraform import privatelink_endpoint_service_adl.test 1112222b3bf99403840e8934--vpce-jjg5e24qp93513h03
```

For more information see: [MongoDB Atlas API - DataLake](https://docs.mongodb.com/datalake/reference/api/datalakes-api/)  and [MongoDB Atlas API - Online Archive](https://docs.atlas.mongodb.com/reference/api/online-archive/) Documentation.
