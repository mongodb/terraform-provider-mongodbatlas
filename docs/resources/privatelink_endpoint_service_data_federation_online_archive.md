# Resource: mongodbatlas_privatelink_endpoint_service_data_federation_online_archive

`mongodbatlas_privatelink_endpoint_service_data_federation_online_archive` provides a Private Endpoint Service resource for Data Federation and Online Archive. The resource allows you to create and manage a private endpoint for Federated Database Instances and Online Archives to the specified project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **NOTE:** Updates are limited to the `comment` argument.

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
```

## Argument Reference

* `project_id` (Required) - Unique 24-hexadecimal digit string that identifies your project. 
* `endpoint_id` (Required) - Unique 22-character alphanumeric string that identifies the private endpoint. See [Atlas Data Lake supports Amazon Web Services private endpoints using the AWS PrivateLink feature](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint:~:text=Atlas%20Data%20Lake%20supports%20Amazon%20Web%20Services%20private%20endpoints%20using%20the%20AWS%20PrivateLink%20feature).
* `provider_name` (Required) - Human-readable label that identifies the cloud service provider. 
* `timeouts`- (Optional) The duration of time to wait for Private Endpoint Service to be created or deleted. The timeout value is definded by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Private Endpoint create & delete is `2h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).
* `region` -  Human-readable label to identify the region of VPC endpoint.  Requires the **Atlas region name**, see the reference list for [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
* `customer_endpoint_dns_name` - (Optional) Human-readable label to identify VPC endpoint DNS name.
* `comment` - (Optional) Human-readable string to associate with this private endpoint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - Human-readable label that identifies the resource type associated with this private endpoint.

## Import

Private Endpoint Service resource for Data Federation and Online Archive can be imported using project ID, endpoint ID, in the format `project_id`--`endpoint_id`, e.g.

```
$ terraform import mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.example 1112222b3bf99403840e8934--vpce-3bf78b0ddee411ba1
```

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint) Documentation for more information.

