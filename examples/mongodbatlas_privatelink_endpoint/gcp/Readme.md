# Example with GCP with legacy PSC architecture and MongoDB Atlas Private Endpoint

This project demonstrates the **legacy PSC architecture** for setting up GCP Private Service Connect with MongoDB Atlas, which requires 50 endpoints.


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
This project currently does the below deployments:

- MongoDB Atlas GCP Private Endpoint (legacy architecture with 50 endpoints)
- Google resource Compute Network, SubNetwork, Address and Forwarding Rule
- Google Private Service Connect (PSC)-MongoDB Private Link

## Architecture Options

This example demonstrates the **legacy PSC architecture** which requires 50 endpoints.

For the **new PSC port-based architecture** (which requires only 1 endpoint and is enabled with `port_mapping_enabled = true`), see the [`gcp-port-based`](../gcp-port-based/) example.

### Architecture Comparison

| Feature | Legacy Architecture (this example) | New Port-Based Architecture |
|---------|-----------------------------------|---------------------------|
| Endpoints Required | 50 | 1 |
| `port_mapping_enabled` | `false` (or omitted) | `true` |

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
