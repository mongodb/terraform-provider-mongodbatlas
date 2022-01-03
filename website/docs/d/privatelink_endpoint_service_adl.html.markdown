---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: privatelink_endpoint_service_adl"
sidebar_current: "docs-mongodbatlas-datasource-privatelink-endpoint-service-adl"
description: |-
Describes an Atlas Data Lake and Online Archive PrivateLink endpoint
---


# privatelink_endpoint_service_adl

`privatelink_endpoint_service_adl` Provides an Atlas Data Lake (ADL) and Online Archive PrivateLink endpoint resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

### Basic
```terraform
resource "mongodbatlas_privatelink_endpoint_service_adl" "adl_test" {
  project_id   = "<PROJECT_ID>"
  endpoint_id  = "<ENDPOINT_ID>"
  comment      = "Comment for PrivateLink endpoint ADL"
  type		 = "DATA_LAKE"
  provider_name	 = "AWS"
}

data "mongodbatlas_privatelink_endpoint_service_adl" "test" {
  project_id            = mongodbatlas_privatelink_endpoint_service_adl.adl_test.project_id
  private_link_id       = mongodbatlas_privatelink_endpoint_service_adl.adl_test.endpoint_id
}
```


## Argument Reference

* `project_id` - (Required) Unique 24-digit hexadecimal string that identifies the project.
* `endpoint_id` - (Required) Unique 22-character alphanumeric string that identifies the private endpoint. Atlas supports AWS private endpoints using the [|aws| PrivateLink](https://aws.amazon.com/privatelink/) feature.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - Human-readable label that identifies the type of resource to associate with this private endpoint.
* `provider_name` - Human-readable label that identifies the cloud provider for this endpoint.
* `comment` - Human-readable string to associate with this private endpoint.

For more information see: [MongoDB Atlas API - DataLake](https://docs.mongodb.com/datalake/reference/api/datalakes-api/)  and [MongoDB Atlas API - Online Archive](https://docs.atlas.mongodb.com/reference/api/online-archive/) Documentation.
