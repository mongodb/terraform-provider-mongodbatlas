# Resource: privatelink_endpoint_service_serverless

`privatelink_endpoint_service_serverless` Provides a Serverless PrivateLink Endpoint Service resource.
This is the second of two resources required to configure PrivateLink for Serverless, the first is [mongodbatlas_privatelink_endpoint_serverless](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/privatelink_endpoint_serverless).

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.
-> **NOTE:** Create waits for all serverless instances on the project to IDLE in order for their operations to complete. This ensures the latest connection strings can be retrieved following creation of this resource. Default timeout is 2hrs.

## Example Usage

## Example with AWS
```terraform

resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
	project_id   = "<PROJECT_ID>"
	instance_name = mongodbatlas_serverless_instance.test.name
	provider_name = "AWS"
}
	  

resource "aws_vpc_endpoint" "ptfe_service" {
  vpc_id             = "vpc-7fc0a543"
  service_name       = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = ["subnet-de0406d2"]
  security_group_ids = ["sg-3f238186"]
}

resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
	project_id   = "<PROJECT_ID>"
	instance_name = mongodbatlas_serverless_instance.test.name
	endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
	cloud_provider_endpoint_id = aws_vpc_endpoint.ptfe_service.id
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
resource  "mongodbatlas_privatelink_endpoint_serverless" "test" {
  project_id    = var.project_id
  provider_name = "AZURE"
}

resource "azurerm_private_endpoint" "test" {
  name                = "endpoint-test"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = var.resource_group_name
  subnet_id           = azurerm_subnet.test.id
  private_service_connection {
    name                           = mongodbatlas_privatelink_endpoint_serverless.test.private_link_service_name
    private_connection_resource_id = mongodbatlas_privatelink_endpoint_serverless.test.private_link_service_resource_id
    is_manual_connection           = true
    request_message                = "Azure Private Link test"
  }

}

resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
  project_id                  = mongodbatlas_privatelink_endpoint_serverless.test.project_id
  instance_name               = mongodbatlas_serverless_instance.test.name
  endpoint_id                 = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
  cloud_provider_endpoint_id  = azurerm_private_endpoint.test.id 
  private_endpoint_ip_address = azurerm_private_endpoint.test.private_service_connection.0.private_ip_address
  provider_name               = "AZURE"
  comment                     = "test"
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
* `instance_name` - (Required) Human-readable label that identifies the serverless instance.
* `endpoint_id` - (Required) Unique 24-hexadecimal digit string that identifies the private endpoint.
* `cloud_provider_endpoint_id` - (Optional) Unique string that identifies the private endpoint's network interface.
* `private_endpoint_ip_address` - (Optional) IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
* `provider_name` - (Required) Cloud provider for which you want to create a private endpoint. Atlas accepts `AWS`, `AZURE`.
* `comment` - (Optional) Human-readable string to associate with this private endpoint.
* `timeouts`- (Optional) The duration of time to wait for Private Endpoint Service to be created or deleted. The timeout value is defined by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Private Endpoint create & delete is `2h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `endpoint_service_name` - Unique string that identifies the PrivateLink endpoint service.
* `private_link_service_resource_id` - Root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages.
* `private_endpoint_ip_address` - IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
* `cloud_provider_endpoint_id` - Unique string that identifies the private endpoint's network interface.
* `comment` - Human-readable string to associate with this private endpoint.
* `error_message` - Human-readable error message that indicates the error condition associated with establishing the private endpoint connection.
* `status` - Human-readable label that indicates the current operating status of the private endpoint. Values include: RESERVATION_REQUESTED, RESERVED, INITIATING, AVAILABLE, FAILED, DELETING.

## Import

Serverless privatelink endpoint can be imported using project ID and endpoint ID, in the format `project_id`--`endpoint_id`, e.g.

```
$ terraform import mongodbatlas_privatelink_endpoint_service_serverless.test 1112222b3bf99403840e8934--serverless_name--vpce-jjg5e24qp93513h03
```

For more information see: [MongoDB Atlas API - Serverless Private Endpoints](https://www.mongodb.com/docs/atlas/reference/api/serverless-private-endpoints-get-one/).
