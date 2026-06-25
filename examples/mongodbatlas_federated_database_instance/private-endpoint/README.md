# Example - MongoDB Atlas Federated Database Instance with Private Endpoint - AWS

This example shows how to configure a [MongoDB Atlas Federated Database Instance](https://www.mongodb.com/docs/atlas/data-federation/adf-overview/overview/) with an AWS Private Endpoint

## Dependencies

* A MongoDB Atlas account
* An AWS account

## Usage

**1\. Set up credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
export AWS_ACCESS_KEY_ID="<AWS_ACCESS_KEY_ID>"
export AWS_SECRET_ACCESS_KEY="<AWS_SECRET_ACCESS_KEY>"
```

**2\. Set required variables.**

Create a `terraform.tfvars` file:

```hcl
project_id       = "<ATLAS_PROJECT_ID>"
vpce_service_name = "com.amazonaws.vpce.us-east-1.vpce-svc-<ID>"
```

You can find the `vpce_service_name` for your region in the [Atlas Data Federation private endpoint documentation](https://www.mongodb.com/docs/atlas/data-federation/tutorial/config-private-endpoint/?cloud-provider=aws&interface=atlas-ui#choose-a-cloud-provider-and-region.).

**3\. Review the Terraform plan.**

```bash
terraform plan
```

**4\. Apply.**

```bash
terraform apply
```

> **Note:** `private_endpoint_hostnames` is populated asynchronously by Atlas and will be empty on the first apply if the private endpoint is created within the same `apply`. After the initial apply completes, run `terraform apply -refresh-only` to update the state with the populated hostnames. If you have downstream resources that consume this value, run `terraform apply` afterwards as well.

**5\. Destroy the resources.**

```bash
terraform destroy
```
