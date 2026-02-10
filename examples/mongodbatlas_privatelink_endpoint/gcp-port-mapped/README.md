# Example with GCP with Port-Mapped Architecture and MongoDB Atlas Private Endpoint

This project demonstrates the **port-mapped architecture** for setting up GCP Private Service Connect with MongoDB Atlas. Unlike the GCP legacy private endpoint architecture that requires dedicated resources for each Atlas node, the port-mapped architecture design uses a single set of resources to support up to 150 nodes, enabling direct targeting of specific nodes using only one customer IP address.

## Architecture Comparison

| Feature | GCP Legacy Private Endpoint Architecture | Port-Mapped Architecture (this example) |
|---------|-------------------|---------------------------|
| Resources per Atlas node | Dedicated forwarding rule, service attachment, and instance group | Single set of resources for up to 150 nodes |
| `port_mapping_enabled` | `false` (or omitted) | `true` |
| Customer IP addresses | One per Atlas node | One total |

## Architecture Overview

The port-mapped architecture uses:
- **1 Google Compute Address** (supports up to 150 nodes)
- **1 Google Compute Forwarding Rule** (supports up to 150 nodes)
- `port_mapping_enabled = true` on the `mongodbatlas_privatelink_endpoint` resource

## Terraform Configuration

The key configuration is the `port_mapping_enabled = true` setting:

```hcl
resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id           = var.project_id
  provider_name        = "GCP"
  region               = var.gcp_region
  port_mapping_enabled = true  # Enables port-mapped architecture
  # ...
}
```
- Use `endpoint_service_id` (forwarding rule name) and `private_endpoint_ip_address` (IP address) in `mongodbatlas_privatelink_endpoint_service`
- The `endpoints` list is not used for the port-mapped architecture


## Dependencies

* Terraform v0.13+
* Google Cloud account
* MongoDB Atlas account

```
Terraform v0.13.0
+ provider registry.terraform.io/hashicorp/google
+ provider registry.terraform.io/terraform-providers/mongodbatlas
```

## Usage

**1\. Ensure your Google credentials are set up.**

1. Install the GCloud SDK by following the steps from the [official GCP documentation](https://cloud.google.com/sdk/docs/install).
2. Run the command `gcloud init` and authenticate with GCP.
3. Once authenticated you will need to select a project to use. After you select a project a success message will appear, see the example below. You are then ready to proceed.
```
⇒  gcloud init
You are logged in as: [user@example.com].

Pick cloud project to use:
 [1] project1
 [2] project2
...

Please enter numeric choice or text value (must exactly match list item): 1

Your Google Cloud SDK is configured and ready to use!

```
**2\. TFVARS**

Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

An existing cluster on the project can optionally be linked via the `cluster_name` variable.
If included, the gcp connection string to the cluster will be output.

**3\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project deploys:

- MongoDB Atlas GCP Private Endpoint
- Google Compute Network, SubNetwork, Address and Forwarding Rule
- Google Private Service Connect-MongoDB Private Link

**4\. Execute the Terraform apply.**

Now execute the plan to provision the GCP resources.

``` bash
$ terraform apply
```

**5\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary GCP and Atlas charges.

``` bash
$ terraform destroy
```

## References

- [Set Up a Private Endpoint for GCP (MongoDB Atlas Documentation)](https://www.mongodb.com/docs/atlas/security-private-endpoint/?cloud-provider=gcp)
- [Migration Guide: GCP Private Service Connect to Port-Mapped Architecture](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/gcp-privatelink-port-mapping-migration)