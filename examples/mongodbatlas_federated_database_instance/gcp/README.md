# Example - MongoDB Atlas Federated Database Instance with Google Cloud provider configuration

This project aims to provide an example of using [MongoDB Atlas Federated Database Instance](https://www.mongodb.com/docs/atlas/data-federation/adf-overview/overview/).

## Dependencies

* A MongoDB Atlas account
* A Google Cloud account

## Usage

**1\. Set up credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**2\. Set required variables.**

Create a `terraform.tfvars` file:

```hcl
project_id              = "<ATLAS_PROJECT_ID>"
federated_instance_name = "<FEDERATED_INSTANCE_NAME>"
```

**3\. Review the Terraform plan.**

Run the following command and review the plan you created.

```bash
$ terraform plan
```

This project currently supports the following deployments:

- MongoDB Atlas Cloud Provider Access Setup for Google Cloud
- MongoDB Atlas Cloud Provider Access Authorization for Google Cloud
- MongoDB Atlas Federated Database Instance with Google Cloud provider configuration

**4\. Run the Terraform apply command to apply the plan.**

Now run the plan to provision the resources.

```bash
$ terraform apply
```

**5\. Authorize the Atlas service account.**

Atlas creates a Google Cloud service account for the federated database instance, which should be authorized to allow Atlas to access Google Cloud Storage. Read its email address from the `gcp_service_account` output and grant it access to the buckets you want to query, for example the `roles/storage.objectViewer` role.

**6\. Destroy the resources.**

Once you are finished with our testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

```bash
$ terraform destroy
```
