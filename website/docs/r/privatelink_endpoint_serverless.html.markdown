---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: privatelink_endpoint_serverless"
sidebar_current: "docs-mongodbatlas-datasource-privatelink-endpoint-serverless"
description: |-
Describes a Serverless PrivateLink Endpoint
---


# Resource: privatelink_endpoint_serverless

`privatelink_endpoint_serverless` Provides a Serverless PrivateLink Endpoint resource.
This is the first of two resources required to configure PrivateLink for Serverless, the second is [mongodbatlas_privatelink_endpoint_service_serverless](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/privatelink_endpoint_service_serverless).

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

### AWS Example
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
* `instance_name` - (Required) Human-readable label that identifies the serverless instance.
* `provider_name` - (Required) Cloud provider name; AWS is currently supported

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `endpoint_id` - Unique 24-hexadecimal digit string that identifies the private endpoint.
* `endpoint_service_name` - Unique string that identifies the PrivateLink endpoint service.
* `private_link_service_resource_id` - Root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages.
* `cloud_provider_endpoint_id` - Unique string that identifies the private endpoint's network interface.
* `comment` - Human-readable string to associate with this private endpoint.
* `status` - Human-readable label that indicates the current operating status of the private endpoint. Values include: RESERVATION_REQUESTED, RESERVED, INITIATING, AVAILABLE, FAILED, DELETING.
* `timeouts`- (Optional) The duration of time to wait for Private Endpoint Service to be created or deleted. The timeout value is defined by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Private Endpoint create & delete is `2h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).

## Import

Serverless privatelink endpoint can be imported using project ID and endpoint ID, in the format `project_id`--`endpoint_id`, e.g.

```
$ terraform import mongodbatlas_privatelink_endpoint_serverless.test 1112222b3bf99403840e8934--serverless_name--vpce-jjg5e24qp93513h03
```

For more information see: [MongoDB Atlas API - Serverless Private Endpoints](https://www.mongodb.com/docs/atlas/reference/api/serverless-private-endpoints-get-one/).
