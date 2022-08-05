---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: privatelink_endpoints_service_adl"
sidebar_current: "docs-mongodbatlas-datasource-privatelink-endpoints-service-adl"
description: |-
Describes the list of all Atlas Data Lake and Online Archive PrivateLink endpoints.
---

# Data Source: privatelink_endpoints_service_adl

`privatelink_endpoints_service_adl` Describes the list of all Atlas Data Lake (ADL) and Online Archive PrivateLink endpoints resource.

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

data "mongodbatlas_privatelink_endpoints_service_adl" "test" {
  project_id            = mongodbatlas_privatelink_endpoint_service_adl.adl_test.project_id
}
```

## Argument Reference

* `project_id`    - (Required) The unique ID for the project.

# Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `links` - The links array includes one or more links to sub-resources or related resources. The relations between URLs are explained in the [Web Linking Specification](http://tools.ietf.org/html/rfc5988).
* `results` - Each element in the `result` array is one private endpoint.
* `total_count` - This value is the count of the total number of items in the result set. `total_count may be greater than the number of objects in the results array if the entire result set is paginated.

### links
Each object in the `links` array represents an online archive with the following attributes:
* `self` - The URL endpoint for this resource.

### results

Each object in the `results` array represents an online archive with the following attributes:

* `endpoint_id` - (Required) Unique 22-character alphanumeric string that identifies the private endpoint. Atlas supports AWS private endpoints using the [|aws| PrivateLink](https://aws.amazon.com/privatelink/) feature.
* `type` - Human-readable label that identifies the type of resource to associate with this private endpoint.
* `provider_name` - Human-readable label that identifies the cloud provider for this endpoint.
* `comment` - Human-readable string to associate with this private endpoint.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/online-archive-get-all-for-cluster/) Documentation for more information.
