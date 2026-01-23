---
subcategory: "Private Endpoint Services"
---

# Resource: mongodbatlas_privatelink_endpoint

`mongodbatlas_privatelink_endpoint` provides a Private Endpoint resource. This represents a [Private Endpoint Service](https://www.mongodb.com/docs/atlas/security-private-endpoint/#private-endpoint-concepts) that can be created in an Atlas project.

~> **IMPORTANT:** This resource creates a Private Endpoint *Service* in MongoDB Atlas. The endpoint itself is created in your cloud provider using the information returned by this resource. The complementary resource `mongodbatlas_privatelink_endpoint_service` is used to link your cloud provider's endpoint to the Atlas service.

The [private link Terraform module](https://registry.terraform.io/modules/terraform-mongodbatlas-modules/private-endpoint/mongodbatlas/latest) makes use of this resource and simplifies its use.

~> **IMPORTANT:**You must have one of the following roles to successfully handle the resource: <br> - Organization Owner <br> - Project Owner

~> **IMPORTANT:** Before configuring a private endpoint for a new region in your cluster,
ensure that you review the [Multi-Region Private Endpoints](https://www.mongodb.com/docs/atlas/troubleshoot-private-endpoints/#multi-region-private-endpoints) troubleshooting documentation.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

-> **NOTE:** A network container is created for a private endpoint to reside in if one does not yet exist in the project.  

~> **IMPORTANT:** For GCP Private Service Connect, MongoDB encourages customers to use the port-based architecture by setting `port_mapping_enabled = true`. The port-based architecture simplifies setup by requiring only 1 endpoint instead of multiple endpoints required by the legacy architecture, and uses a single set of resources to support up to 1000 nodes. For migration guidance, see the [GCP Private Link Port Mapping Migration Guide](../guides/gcp-privatelink-port-mapping-migration.md).

## Example Usage

```terraform
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = "<PROJECT-ID>"
  provider_name = "AWS/AZURE"
  region        = "US_EAST_1"

  timeouts {
    create = "30m"
    delete = "20m"
  }
}
```

### Further Examples
- [AWS PrivateLink Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.2.0/examples/mongodbatlas_privatelink_endpoint/aws)
- [Azure PrivateLink Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.2.0/examples/mongodbatlas_privatelink_endpoint/azure)
- [GCP Private Service Connect Endpoint (Legacy Architecture)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.2.0/examples/mongodbatlas_privatelink_endpoint/gcp)
- [GCP Private Service Connect Endpoint (Port-Based Architecture)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.2.0/examples/mongodbatlas_privatelink_endpoint/gcp-port-based)

## Argument Reference

* `project_id` - Required 	Unique identifier for the project.
* `provider_name` - (Required) Name of the cloud provider for which you want to create the private endpoint service. Atlas accepts `AWS`, `AZURE` or `GCP`.
* `region` - (Required) Cloud provider region in which you want to create the private endpoint connection.
Accepted values are: [AWS regions](https://docs.atlas.mongodb.com/reference/amazon-aws/#amazon-aws), [AZURE regions](https://docs.atlas.mongodb.com/reference/microsoft-azure/#microsoft-azure) and [GCP regions](https://docs.atlas.mongodb.com/reference/google-gcp/#std-label-google-gcp)
* `timeouts`- (Optional) The duration of time to wait for Private Endpoint to be created or deleted. The timeout value is defined by a signed sequence of decimal numbers with a time unit suffix such as: `1h45m`, `300s`, `10m`, etc. The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Private Endpoint create & delete is `1h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).
* `delete_on_create_timeout`- (Optional) Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.
* `port_mapping_enabled` - (Optional) Flag that indicates whether this endpoint service uses GCP port-mapping. When set to `true`, enables the new port-based architecture, which requires only 1 endpoint. Defaults to `false`. Only applicable for GCP provider.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `private_link_id` - Unique identifier of the AWS PrivateLink connection.
* `error_message` - Error message pertaining to the AWS PrivateLink connection. Returns null if there are no errors.
AWS: 
  * `endpoint_service_name` - Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
  * `interface_endpoints` - Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
AZURE:
  * `private_endpoints` - All private endpoints that you have added to this Azure Private Link Service.
  * `private_link_service_name` - Name of the Azure Private Link Service that Atlas manages.
GCP: 
  * `endpoint_group_names` - For GCP legacy architecture (when `port_mapping_enabled` is not set or `false` on the endpoint resource): A list of the endpoint group names associated with the private endpoint service. For GCP port-based architecture (when `port_mapping_enabled = true` on the endpoint resource): A list of private endpoint names associated with the private endpoint service.
  * `region_name` - GCP region for the Private Service Connect endpoint service.
  * `service_attachment_names` - For GCP legacy architecture (when `port_mapping_enabled` is not set or `false` on the endpoint resource): A list of service attachments connected to the private endpoint service (one per Atlas node). For GCP port-based architecture (when `port_mapping_enabled = true` on the endpoint resource): A list of one service attachment connected to the private endpoint service. Returns an empty list while Atlas creates the service attachments.
* `status` - Status of the AWS PrivateLink connection or Status of the Azure/GCP Private Link Service. Atlas returns one of the following values:
  AWS:
    * `AVAILABLE` 	Atlas is creating the network load balancer and VPC endpoint service.
    * `WAITING_FOR_USER` The Atlas network load balancer and VPC endpoint service are created and ready to receive connection requests. When you receive this status, create an interface endpoint to continue configuring the AWS PrivateLink connection.
    * `FAILED` 	A system failure has occurred.
    * `DELETING` 	The AWS PrivateLink connection is being deleted.
  AZURE:
    * `AVAILABLE` 	Atlas created the load balancer and the Private Link Service.
    * `INITIATING` 	Atlas is creating the load balancer and the Private Link Service.
    * `FAILED` 	Atlas failed to create the load balancer and the Private Link service.
    * `DELETING` 	Atlas is deleting the Private Link service.
  GCP:
    * `AVAILABLE` 	Atlas created the load balancer and the GCP Private Service Connect service.
    * `INITIATING` 	Atlas is creating the load balancer and the GCP Private Service Connect service.
    * `FAILED`  	Atlas failed to create the load balancer and the GCP Private Service Connect service.
    * `DELETING` 	Atlas is deleting the GCP Private Service Connect service.

## Import
Private Endpoint Service can be imported using project ID, private link ID, provider name and region, in the format `{project_id}-{private_link_id}-{provider_name}-{region}`, e.g.

```
$ terraform import mongodbatlas_privatelink_endpoint.test 1112222b3bf99403840e8934-3242342343112-AWS-us-east-1
```

See detailed information for arguments and attributes: [MongoDB API Private Endpoint Service](https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-create-one/)
