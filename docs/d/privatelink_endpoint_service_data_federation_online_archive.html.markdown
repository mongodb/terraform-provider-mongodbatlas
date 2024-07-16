# Data Source: mongodbatlas_privatelink_endpoint_service_data_federation_online_archive

`mongodbatlas_privatelink_endpoint_service_data_federation_online_archive` describes a Private Endpoint Service resource for Data Federation and Online Archive.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

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

data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test_data_source" {
  project_id = mongodbatlas_project.atlas-project.id
  endpoint_id = mongodbatlas_privatelink_endpoint_service_data_federation_online_archive.test.endpoint_id
}
```


## Argument Reference

* `project_id` (Required) - Unique 24-hexadecimal digit string that identifies your project. 
* `endpoint_id` (Required) - Unique 22-character alphanumeric string that identifies the private endpoint. See [Atlas Data Lake supports Amazon Web Services private endpoints using the AWS PrivateLink feature](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint:~:text=Atlas%20Data%20Lake%20supports%20Amazon%20Web%20Services%20private%20endpoints%20using%20the%20AWS%20PrivateLink%20feature).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `type` - Human-readable label that identifies the resource type associated with this private endpoint.
* `comment` - Human-readable string to associate with this private endpoint.
* `provider_name` - Human-readable label that identifies the cloud service provider. 
* `region` -  Human-readable label to identify the region of VPC endpoint.  Requires the **Atlas region name**, see the reference list for [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
* `customer_endpoint_dns_name` - (Optional) Human-readable label to identify VPC endpoint DNS name.

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint) Documentation for more information.

