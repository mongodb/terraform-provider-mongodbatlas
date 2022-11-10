---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: privatelink_endpoint_serverless"
sidebar_current: "docs-mongodbatlas-datasource-privatelink-endpoint-serverless"
description: |-
Describes a Serverless PrivateLink Endpoint
---


# Data Source: privatelink_endpoint_serverless

`privatelink_endpoint_serverless` Provides a Serverless PrivateLink Endpoint resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

### Basic
```terraform

resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
	project_id   = "<PROJECT_ID>"
	instance_name = mongodbatlas_serverless_instance.test.name
	provider_name = "AWS"
}
	  
resource "mongodbatlas_serverless_instance" "test" {
	project_id   = "<PROJECT_ID>"
	name         = "test-db"
	provider_settings_backing_provider_name = "AWS"
	provider_settings_provider_name = "SERVERLESS"
	provider_settings_region_name = "US_EAST_1"
	continuous_backup_enabled = true
}
```


## Argument Reference

* `project_id` - (Required) Unique 24-digit hexadecimal string that identifies the project.
* `instance_name` - (Required) Human-readable label that identifies the serverless instance associated with the tenant endpoint
* `provider_name` - (Required) Cloud provider name; AWS is currently supported

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `endpoint_id` - Unique 22-character alphanumeric string that identifies the private endpoint. Atlas supports AWS private endpoints using the [|aws| PrivateLink](https://aws.amazon.com/privatelink/) feature.
* `endpoint_service_name` - Unique string that identifies the PrivateLink endpoint service. MongoDB Cloud returns null while it creates the endpoint service.
* `cloud_provider_endpoint_id` - Unique string that identifies the private endpoint's network interface.
* `comment` - Human-readable string to associate with this private endpoint.
* `status` - Human-readable label that indicates the current operating status of the private endpoint. Values include: RESERVATION_REQUESTED, RESERVED, INITIATING, AVAILABLE, FAILED, DELETING.

For more information see: [MongoDB Atlas API - Serverless Private Endpoints](https://www.mongodb.com/docs/atlas/reference/api/serverless-private-endpoints-get-one/)  and [MongoDB Atlas API - Online Archive](https://docs.atlas.mongodb.com/reference/api/online-archive/) Documentation.
