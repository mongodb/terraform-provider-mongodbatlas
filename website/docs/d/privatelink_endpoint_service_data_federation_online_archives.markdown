---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_privatelink_endpoint_service_data_federation_online_archives"
sidebar_current: "docs-mongodbatlas-data-source-privatelink-endpoint-service-data-federation-online-archives"
description: |-
    Provides a data source for a Private Endpoints Service Data Federation Online Archive.
---

# Data Source: mongodbatlas_privatelink_endpoint_service_data_federation_online_archives

`mongodbatlas_privatelink_endpoint_service_data_federation_online_archives` describes Private Endpoint Service resources for Data Federation and Online Archive.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_project" "atlas-project" {
  org_id = var.atlas_org_id
  name   = var.atlas_project_name
}

resource "mongodbatlas_privatelink_endpoint_service_data_federation_online_archive" "test" {
  project_id = mongodbatlas_project.atlas-project.id
  endpoint_id = "<PRIVATE-ENDPOINT-SERVICE-ID>"
  provider_name = "AWS"
  type = "DATA_LAKE"
  comment = "Test"
}

data "mongodbatlas_privatelink_endpoint_service_data_federation_online_archives" "test_data_source" {
  project_id = mongodbatlas_project.atlas-project.id
}



## Argument Reference

* `project_id` (Required) - Unique 24-hexadecimal digit string that identifies your project. 

## Attributes Reference
* `results` - A list where each represents a Private Endpoint Service


### Private Endpoint Service 
In addition to all arguments above, the following attributes are exported:

* `endpoint_id` - Unique 22-character alphanumeric string that identifies the private endpoint. See [Atlas Data Lake supports Amazon Web Services private endpoints using the AWS PrivateLink feature](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint:~:text=Atlas%20Data%20Lake%20supports%20Amazon%20Web%20Services%20private%20endpoints%20using%20the%20AWS%20PrivateLink%20feature.).
* `type` - Human-readable label that identifies the resource type associated with this private endpoint.
* `comment` - Human-readable string to associate with this private endpoint.
* `provider_name` - Human-readable label that identifies the cloud service provider. 

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint) Documentation for more information.

