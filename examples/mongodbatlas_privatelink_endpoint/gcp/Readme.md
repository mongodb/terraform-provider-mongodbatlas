# Example with GCP with legacy PSC architecture and MongoDB Atlas Private Endpoint

This project demonstrates the **legacy GCP architecture** for setting up GCP Private Service Connect with MongoDB Atlas. The legacy architecture requires dedicated resources for each Atlas node (a dedicated customer forwarding rule, service attachment, internal forwarding rule, and instance group per node). Unlike the new port-based architecture that uses a single set of resources to support up to 1000 nodes through port mapping, the legacy design requires one customer IP address per Atlas node.

## Architecture Comparison

| Feature | Legacy Architecture (this example) | New Port-Based Architecture |
|---------|-----------------------------------|---------------------------|
| Resources per Atlas node | Dedicated forwarding rule, service attachment, and instance group | Single set of resources for up to 1000 nodes |
| `port_mapping_enabled` | `false` (or omitted) | `true` |
| Customer IP addresses | One per Atlas node | One total |

For the **new GCP port-based architecture** (enabled with `port_mapping_enabled = true`), see the [`gcp-port-based`](../gcp-port-based/) example.

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

- MongoDB Atlas GCP Private Endpoint (legacy architecture)
- Google Compute Network, SubNetwork, Address and Forwarding Rule
- Google Private Service Connect (PSC)-MongoDB Private Link

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
