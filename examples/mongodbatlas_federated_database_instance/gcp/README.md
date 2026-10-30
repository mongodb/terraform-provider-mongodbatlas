# Example - MongoDB Atlas Federated Database Instance with Google Cloud Storage

This project aims to provide an example of using [MongoDB Atlas Federated Database Instance](https://www.mongodb.com/docs/atlas/data-federation/adf-overview/overview/).

## Dependencies

* Terraform MongoDB Atlas Provider v2.20.0
* A MongoDB Atlas account
* A Google Cloud account with permission to create Cloud Storage buckets

```
+ provider registry.terraform.io/mongodb/mongodbatlas v2.20.0
```

## Usage

**1\. Set up credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

Alternatively, set `atlas_client_id` and `atlas_client_secret` in the `terraform.tfvars` file described in the next step.

The Google Cloud provider uses [Application Default Credentials](https://cloud.google.com/docs/authentication/provide-credentials-adc):

```bash
gcloud auth application-default login
```

**2\. Set required variables.**

Create a `terraform.tfvars` file:

```hcl
project_id      = "<ATLAS_PROJECT_ID>"
gcp_project_id  = "<GCP_PROJECT_ID>"
gcs_bucket_name = "<GCS_BUCKET_NAME>"

# Optional, if you are not using the environment variables from the previous step.
# atlas_client_id     = "<ATLAS_CLIENT_ID>"
# atlas_client_secret = "<ATLAS_CLIENT_SECRET>"
```

The bucket name must be globally unique. `federated_instance_name`, `gcs_bucket_location` and `data_process_region` are optional, and are documented with their defaults in `variables.tf`.

**3\. Review the Terraform plan.**

Run the following command and review the plan you created.

```bash
$ terraform plan
```

This project currently supports the following deployments:

- A Google Cloud Storage bucket
- An IAM binding granting the Atlas service account read access to the bucket
- MongoDB Atlas Cloud Provider Access Setup for Google Cloud
- MongoDB Atlas Cloud Provider Access Authorization for Google Cloud
- MongoDB Atlas Federated Database Instance querying the bucket

**4\. Run the Terraform apply command to apply the plan.**

Now run the plan to provision the resources.

```bash
$ terraform apply
```

**5\. Destroy the resources.**

Once you are finished with our testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

```bash
$ terraform destroy
```
