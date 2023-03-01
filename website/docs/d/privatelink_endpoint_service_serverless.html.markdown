---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: privatelink_endpoint_service_serverless"
sidebar_current: "docs-mongodbatlas-datasource-privatelink-endpoint-service-serverless"
description: |-
Describes a Serverless PrivateLink Endpoint Service
---


# Data Source: privatelink_endpoint_service_serverless

`privatelink_endpoint_service_serverless` Provides a Serverless PrivateLink Endpoint Service resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

## Example with AWS
```terraform

data "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
  project_id   = "<PROJECT_ID>"
  instance_name = mongodbatlas_serverless_instance.test.name
  endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
}

resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
	project_id   = "<PROJECT_ID>"
	instance_name = mongodbatlas_serverless_instance.test.name
	provider_name = "AWS"
}
	  
	  
resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
	project_id   = "<PROJECT_ID>"
	instance_name = "test-db"
	endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
	provider_name = "AWS"
	comment = "New serverless endpoint"
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

## Example with AZURE
```terraform

data "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
  project_id   = "<PROJECT_ID>"
  instance_name = mongodbatlas_serverless_instance.test.name
  endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
}

resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
	project_id   = "<PROJECT_ID>"
	instance_name = mongodbatlas_serverless_instance.test.name
	provider_name = "AZURE"
}
	  
	  
resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
	project_id   = "<PROJECT_ID>"
	instance_name = "test-db"
	endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
	provider_name = "AZURE"
	comment = "New serverless endpoint"
}

resource "mongodbatlas_serverless_instance" "test" {
	project_id   = "<PROJECT_ID>"
	name         = "test-db"
	provider_settings_backing_provider_name = "AZURE"
	provider_settings_provider_name = "SERVERLESS"
	provider_settings_region_name = "US_EAST"
	continuous_backup_enabled = true
}
```

### Available complete examples
- [Setup private connection to a MongoDB Atlas Serverless Instance with AWS VPC](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/aws-privatelink-endpoint/serverless-instance)

## Argument Reference

* `project_id` - (Required) Unique 24-digit hexadecimal string that identifies the project.
* `instance_name` - (Required) Human-readable label that identifies the serverless instance 
* `endpoint_id` - (Required) Unique 22-character alphanumeric string that identifies the private endpoint. Atlas supports AWS private endpoints using the [AWS PrivateLink](https://aws.amazon.com/privatelink/) feature.
* `cloud_provider_endpoint_id` - Unique string that identifies the private endpoint's network interface.
* `comment` - Human-readable string to associate with this private endpoint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `endpoint_service_name` - Unique string that identifies the PrivateLink endpoint service. MongoDB Cloud returns null while it creates the endpoint service.
* `private_link_service_resource_id` - Root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages.
* `private_endpoint_ip_address` - IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
* `status` - Human-readable label that indicates the current operating status of the private endpoint. Values include: RESERVATION_REQUESTED, RESERVED, INITIATING, AVAILABLE, FAILED, DELETING.

For more information see: [MongoDB Atlas API - Serverless Private Endpoints](https://www.mongodb.com/docs/atlas/reference/api/serverless-private-endpoints-get-one/).
