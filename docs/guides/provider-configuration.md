---
page_title: "Guide: Provider Configuration"
---

# Provider Configuration

This guide covers authentication and configuration options for the MongoDB Atlas Provider.

## Authentication Methods

The MongoDB Atlas provider supports the following authentication methods:

1. [**Service Account (SA)** - Recommended](#service-account-recommended)
2. [**Programmatic Access Key (PAK)**](#programmatic-access-key)

Credentials can be provided through (in priority order):

- AWS Secrets Manager
- Provider attributes
- Environment variables

The provider uses the first available credentials source.

### Service Account (Recommended)

SAs simplify authentication by eliminating the need to create new Atlas-specific user identities and permission credentials. See [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/) for more information.

To use SA authentication, create an SA in your [MongoDB Atlas organization](https://www.mongodb.com/docs/atlas/configure-api-access/#grant-programmatic-access-to-an-organization) and set the credentials, for example:

```terraform
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

**Note:** SAs can't be used with `mongodbatlas_event_trigger` resources because its API doesn't support it yet.

**Troubleshooting Service Accounts**

If you encounter a rate limit error when using Service Accounts, you may see:

```
│ Error: error initializing provider: oauth2: cannot fetch token: 429 Too Many Requests
│ Response: {"detail":"Resource /api/oauth/token is limited to 10 requests every 1 minutes.","error":429,"errorCode":"RATE_LIMITED","parameters":["/api/oauth/token",10,1],"reason":"Too Many Requests"}
```

Atlas enforces rate limiting for each combination of IP address and SA client, see [MongoDB Atlas Service Account Limits](https://www.mongodb.com/docs/manual/reference/limits/#mongodb-atlas-service-account-limits) for more information. Each Terraform operation generates a new token that is used for the duration of that operation. These limits work well for individual development environments. For CI pipelines or enterprise environments with shared infrastructure, consider optimizing your configuration using one of these approaches:

- Contact [MongoDB Support](https://support.mongodb.com/) to request a rate limit increase for your organization
- Create separate Service Accounts for different environments or CI/CD pipelines, as each SA client has its own rate limit quota
- Distribute Terraform executions across different IP addresses, since rate limits apply per IP and SA client combination
- Add retry logic to your automation workflows to handle temporary rate limit errors gracefully.

### Programmatic Access Key

Generate a PAK with the appropriate [role](https://docs.atlas.mongodb.com/reference/user-roles/). See the [MongoDB Atlas documentation](https://www.mongodb.com/docs/atlas/configure-api-access-org/) for detailed instructions.

**Role recommendation:** If unsure which role to grant, use an organization API key with the Organization Owner role to ensure sufficient access as in the following example:

```terraform
provider "mongodbatlas" {
  public_key  = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}
```

~> **Migrating from PAK to SA:** Update your provider attributes or environment variables to use SA credentials instead of PAK credentials, then run `terraform plan` to verify everything works correctly.

## AWS Secrets Manager

The provider supports retrieving credentials from AWS Secrets Manager. See [AWS Secrets Manager documentation](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html) for more details.

### Setup Instructions

1. **Create secrets in AWS Secrets Manager**

   For SA, create a secret with the following key-value pairs:
   - `client_id`: your-client-id
   - `client_secret`: your-client-secret

   For PAK, create a secret with the following key-value pairs:
   - `public_key`: your-public-key
   - `private_key`: your-private-key

2. **Create an IAM Role** with:
   - Permission for `sts:AssumeRole`
   - Attached AWS managed policy `SecretsManagerReadWrite`

3. **Configure AWS credentials** (using AWS CLI or environment variables)

4. **Assume the role** to obtain STS credentials

   ```shell
   aws sts assume-role --role-arn <ROLE_ARN> --role-session-name newSession
   ```

5. **Configure provider with AWS Secrets Manager**

   Using provider attributes:

   ```terraform
   provider "mongodbatlas" {
      aws_access_key_id     = var.aws_access_key_id
      aws_secret_access_key = var.aws_secret_access_key
      aws_session_token     = var.aws_session_token
      assume_role           = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts"
      secret_name           = "mongodbsecret"
      region                = "us-east-2"
   }
   ```

   Alternatively, you can use environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`, `ASSUME_ROLE_ARN`, `SECRET_NAME`, `AWS_REGION`).

### Cross-Account and Cross-Region Access

For cross-account secrets, use the fully qualified ARN for `secret_name`. For cross-region or cross-account access, the `sts_endpoint` parameter is required, for example:

```terraform
provider "mongodbatlas" {
   aws_access_key_id     = var.aws_access_key_id
   aws_secret_access_key = var.aws_secret_access_key
   aws_session_token     = var.aws_session_token
   assume_role    = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts"
   secret_name    = "arn:aws:secretsmanager:us-east-1:<AWS_ACCOUNT_ID>:secret:test789-TO06Hy"
   region         = "us-east-2"
   sts_endpoint   = "https://sts.us-east-2.amazonaws.com/"
}
```

## Provider Configuration Reference

### Provider Arguments

* `client_id` - (Optional) SA Client ID (env: `MONGODB_ATLAS_CLIENT_ID`).
* `client_secret` - (Optional) SA Client Secret (env: `MONGODB_ATLAS_CLIENT_SECRET`).
* `access_token` - (Optional) SA Access Token (env: `MONGODB_ATLAS_ACCESS_TOKEN`). Instead of using Client ID and Client Secret, you can generate and use an SA token directly. See [Generate Service Account Token](https://www.mongodb.com/docs/atlas/api/service-accounts/generate-oauth2-token/#std-label-generate-oauth2-token-atlas) for details. Note: tokens have expiration times.
* `public_key` - (Optional) PAK Public Key (env: `MONGODB_ATLAS_PUBLIC_API_KEY`).
* `private_key` - (Optional) PAK Private Key (env: `MONGODB_ATLAS_PRIVATE_API_KEY`).
* `base_url` - (Optional) MongoDB Atlas Base URL (env: `MONGODB_ATLAS_BASE_URL`). For advanced use cases, you can configure custom API endpoints.
* `realm_base_url` - (Optional) MongoDB Realm Base URL (env: `MONGODB_REALM_BASE_URL`).
* `is_mongodbgov_cloud` - (Optional) Set to `true` to use MongoDB Atlas for Government, a dedicated deployment option for government agencies and contractors requiring FedRAMP compliance. When enabled, the provider uses government-specific API endpoints. Ensure credentials are created in the government environment. See [Atlas for Government Considerations](https://www.mongodb.com/docs/atlas/government/api/#atlas-for-government-considerations) for feature limitations and requirements.
  ```terraform
  provider "mongodbatlas" {
    client_id           = var.mongodbatlas_client_id
    client_secret       = var.mongodbatlas_client_secret
    is_mongodbgov_cloud = true
  }
  ```
* `assume_role` - (Optional) AWS IAM role configuration for accessing secrets in AWS Secrets Manager. Role ARN env: `ASSUME_ROLE_ARN`. See [AWS Secrets Manager](#aws-secrets-manager) section for details.
* `secret_name` - (Optional) Name of the secret in AWS Secrets Manager (env: `SECRET_NAME`).
* `region` - (Optional) AWS region where the secret is stored (env: `AWS_REGION`).
* `aws_access_key_id` - (Optional) AWS Access Key ID (env: `AWS_ACCESS_KEY_ID`).
* `aws_secret_access_key` - (Optional) AWS Secret Access Key (env: `AWS_SECRET_ACCESS_KEY`).
* `aws_session_token` - (Optional) AWS Session Token (env: `AWS_SESSION_TOKEN`).
* `sts_endpoint` - (Optional) AWS STS endpoint (env: `STS_ENDPOINT`).

## Credential Priority

When multiple credentials are provided in the same source, the provider uses this priority order:

1. Access Token
2. Service Account (SA)
3. Programmatic Access Key (PAK)

The provider displays a warning when multiple credentials are detected.

## Supported OS and Architectures

As per [HashiCorp's recommendations](https://developer.hashicorp.com/terraform/registry/providers/os-arch), the MongoDB Atlas Provider fully supports the following operating system / architecture combinations:

- Darwin / AMD64
- Darwin / ARMv8
- Linux / AMD64
- Linux / ARMv8 (AArch64/ARM64)
- Linux / ARMv6
- Windows / AMD64

We ship binaries but do not prioritize fixes for the following operating system / architecture combinations:

- Linux / 386
- Windows / 386
- FreeBSD / 386
- FreeBSD / AMD64

## Additional Resources

- [MongoDB Atlas API Documentation](https://www.mongodb.com/docs/atlas/api/)
- [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/)
- [Configure API Access](https://www.mongodb.com/docs/atlas/configure-api-access/)
- [Atlas for Government](https://www.mongodb.com/docs/atlas/government/)
- [Terraform Provider Documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs)
