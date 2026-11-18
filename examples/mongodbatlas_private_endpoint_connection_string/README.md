# MongoDB Atlas Private Endpoint Connection String

Creates an AWS private endpoint for an Atlas project, links it to Atlas, and creates a customer-managed private endpoint connection string for an existing cluster. The example then reads the connection string with the singular and plural data sources (the plural one filtered to `SELF_MANAGED` records).

## Prerequisites

- An Atlas project that is opted in to private endpoint connection strings. Enable the project setting `isPrivateEndpointConnectionStringsEnabled` with the [Update Project Settings](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-updategroupsettings) API operation. When a project opts in, Atlas copies each existing generated private endpoint connection string into a read-only `IMPORTED` record with the same hostname. This example creates and manages a `SELF_MANAGED` connection string.
- An existing dedicated AWS cluster in the project that runs in the same AWS region as the private endpoint (`aws_region` and `atlas_region_name` must refer to the same region).
- An existing AWS VPC with subnets and security groups in that region.
- MongoDB Atlas credentials via environment variables:

  ```bash
  export MONGODB_ATLAS_CLIENT_ID="your-service-account-client-id"
  export MONGODB_ATLAS_CLIENT_SECRET="your-service-account-client-secret"
  ```

- AWS credentials via the standard AWS credential chain, for example:

  ```bash
  export AWS_ACCESS_KEY_ID="your-access-key-id"
  export AWS_SECRET_ACCESS_KEY="your-secret-access-key"
  ```

## Configuration

Supply the required input variables, for example in a `terraform.tfvars` file:

```hcl
project_id         = "0123456789abcdef01234567"
cluster_name       = "my-cluster"
vpc_id             = "vpc-0abc1234567890def"
subnet_ids         = ["subnet-0abc1234567890def", "subnet-0def1234567890abc"]
security_group_ids = ["sg-0abc1234567890def"]
```

## Usage

Run the example with:

```bash
terraform init
terraform apply
```

Create and update operations are asynchronous. Terraform waits until the connection string status is `AVAILABLE` (default 2 hours) and fails the apply if Atlas reports `FAILED`.

Destroy the configuration when you finish testing:

```bash
terraform destroy
```

For more information see the [MongoDB Atlas API documentation](https://mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createGroupClusterPrivateEndpointConnectionString).
