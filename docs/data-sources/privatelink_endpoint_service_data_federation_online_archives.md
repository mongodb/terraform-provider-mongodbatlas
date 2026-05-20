---
subcategory: "Data Federation"
---

# Data Source: mongodbatlas_privatelink_endpoint_service_data_federation_online_archives

`mongodbatlas_privatelink_endpoint_service_data_federation_online_archives` describes Private Endpoint Service resources for Data Federation and Online Archive.

## Example Usage

```terraform
resource "mongodbatlas_project" "atlas-project" {
  org_id = var.atlas_org_id
  name   = var.atlas_project_name
}

resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
  project_id = mongodbatlas_project.atlas-project.id
  endpoint_id = "vpce-046cf43c79424d4c9"
  provider_name = "AWS"
  comment = "Test"
  region        = "US_EAST_1"
  customer_endpoint_dns_name = "vpce-046cf43c79424d4c9-nmls2y9k.vpce-svc-0824460b72e1a420e.us-east-1.vpce.amazonaws.com"
}

data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archives" "test_data_source" {
  project_id = mongodbatlas_project.atlas-project.id
}
```


## Argument Reference

* `project_id` (Required) - Unique 24-hexadecimal digit string that identifies your project, also known as `groupId` in the official documentation.

## Attributes Reference
* `results` - A list where each represents a Private Endpoint Service


### Private Endpoint Service 
In addition to all arguments above, the following attributes are exported:

* `endpoint_id` - Unique string that identifies the private endpoint. For AWS, this is a 22-character alphanumeric string (e.g. `vpce-xxxxxxxxxxxxxxxxx`). For Azure, this is the full resource ID of the private endpoint (e.g. `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateEndpoints/{privateEndpointName}`). See [Atlas Data Federation supports AWS and Azure private endpoints using the PrivateLink feature](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint).
* `type` - Human-readable label that identifies the resource type associated with this private endpoint.
* `comment` - Human-readable string to associate with this private endpoint.
* `provider_name` - Human-readable label that identifies the cloud service provider. Atlas Data Federation supports `AWS` and `AZURE`.
* `region` - Human-readable region label for the customer's endpoint. For `AWS`, if defined, you must also specify a value for `customer_endpoint_dns_name`.
* `customer_endpoint_dns_name` - (Optional) Human-readable DNS name to identify the customer's endpoint. If defined, you must also specify a value for `region`.
* `customer_endpoint_ip_address` - IP address used to connect to the Azure private endpoint.
* `azure_link_id` - Link ID that identifies the Azure private endpoint connection.
* `error_message` - Error message describing a failure approving the private endpoint request.
* `status` - Status of the private endpoint connection request.

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint) Documentation for more information.

