---
page_title: "Guide: Provider Configuration"
---

# Provider Configuration

This guide covers authentication and configuration options for the MongoDB Atlas Provider.

## Authentication Methods

The MongoDB Atlas provider supports the following authentication methods:

1. **Service Account (SA)** - Recommended
2. **Programmatic Access Key (PAK)**

Credentials can be provided through (in priority order):
- AWS Secrets Manager
- Provider attributes
- Environment variables

The provider uses the first available credentials source.

### Service Account (Recommended)

SAs simplify authentication by eliminating the need to create new Atlas-specific user identities and permission credentials. See [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/) and [MongoDB Atlas Service Account Limits](https://www.mongodb.com/docs/manual/reference/limits/#mongodb-atlas-service-account-limits) for more information.

Create an SA in your [MongoDB Atlas organization](https://www.mongodb.com/docs/atlas/configure-api-access/#grant-programmatic-access-to-an-organization) and set the credentials, for example:

```terraform
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

**Note:** SAs can't be used with `mongodbatlas_event_trigger` resources as its API doesn't support it yet.

### Programmatic Access Key

Generate a PAK with the appropriate [role](https://docs.atlas.mongodb.com/reference/user-roles/). See [MongoDB Atlas documentation](https://www.mongodb.com/docs/atlas/configure-api-access-org/) for instructions.

**Role recommendation:** If unsure which role to grant, use an organization API key with the Organization Owner role to ensure sufficient access.

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

3. **Configure AWS credentials**
   ```shell
   export AWS_ACCESS_KEY_ID='<AWS_ACCESS_KEY_ID>'
   export AWS_SECRET_ACCESS_KEY='<AWS_SECRET_ACCESS_KEY>'
   ```

4. **Assume the role**
   ```shell
   aws sts assume-role --role-arn <ROLE_ARN> --role-session-name newSession
   ```

5. **Store STS credentials**
   ```shell
   export AWS_ACCESS_KEY_ID='<STS_ACCESS_KEY_ID>'
   export AWS_SECRET_ACCESS_KEY='<STS_SECRET_ACCESS_KEY>'
   export AWS_SESSION_TOKEN='<STS_SESSION_TOKEN>'
   ```

6. **Set provider configuration via environment variables**
   ```shell
   export ASSUME_ROLE_ARN='arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts'
   export SECRET_NAME='mongodbsecret'
   export AWS_REGION='us-east-2'
   ```

Alternatively, you can configure these settings using provider attributes instead of environment variables.

### Cross-Account and Cross-Region Access

For cross-account secrets, use the fully qualified ARN for `secret_name`. For cross-region or cross-account access, the `sts_endpoint` parameter is required.

Example with environment variables:
```shell
export ASSUME_ROLE_ARN='arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts'
export SECRET_NAME='arn:aws:secretsmanager:us-east-1:<AWS_ACCOUNT_ID>:secret:test789-TO06Hy'
export AWS_REGION='us-east-2'
export STS_ENDPOINT='https://sts.us-east-2.amazonaws.com/'
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

## Security Best Practices

- Never hard-code credentials in Terraform configuration files
- Use environment variables or a secrets management system
- Regularly rotate credentials
- Apply the principle of least privilege when assigning roles
- Use Terraform's `sensitive` attribute for credential variables

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
