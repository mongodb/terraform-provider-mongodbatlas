---
page_title: "Provider Configuration"
---

# Provider Configuration

This guide provides information about configuring the MongoDB Atlas Provider, including authentication methods, environment configuration, and special deployment scenarios.

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

**Migrating from PAK to SA:** To migrate from PAK to SA, simply update your provider attributes or environment variables to use SA credentials instead of PAK credentials, then run `terraform plan` to verify everything works correctly.

## AWS Secrets Manager

AWS Secrets Manager helps manage, retrieve, and rotate credentials throughout their lifecycles. See [AWS Secrets Manager documentation](https://docs.aws.amazon.com/secretsmanager/latest/userguide/intro.html) for more details.

### Setup Instructions

1. **Create secrets in AWS Secrets Manager**

   For SA:
   ```json
   {
     "client_id": "your-client-id",
     "client_secret": "your-client-secret"
   }
   ```

   For PAK:
   ```json
   {
     "public_key": "your-public-key",
     "private_key": "your-private-key"
   }
   ```

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
   export AWS_SESSION_TOKEN="<STS_SESSION_TOKEN>"
   ```

6. **Configure the provider**
   ```terraform
   provider "mongodbatlas" {
     assume_role {
       role_arn = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts"
     }
     secret_name = "mongodbsecret"
     region      = "us-east-2"
   }
   ```

### Cross-Account and Cross-Region Access

For cross-account secrets, use the fully qualified ARN for `secret_name`. For cross-region or cross-account access, the `sts_endpoint` parameter is required.

Example:
```terraform
provider "mongodbatlas" {
  assume_role {
    role_arn = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts"
  }
  secret_name  = "arn:aws:secretsmanager:us-east-1:<AWS_ACCOUNT_ID>:secret:test789-TO06Hy"
  region       = "us-east-2"
  sts_endpoint = "https://sts.us-east-2.amazonaws.com/"
}
```

## Provider Configuration Reference

### Provider Arguments

* `client_id` - (Optional) SA Client ID. Can also be set with the `MONGODB_ATLAS_CLIENT_ID` environment variable.
* `client_secret` - (Optional) SA Client Secret. Can also be set with the `MONGODB_ATLAS_CLIENT_SECRET` environment variable.
* `access_token` - (Optional) SA Access Token. Can also be set with the `MONGODB_ATLAS_ACCESS_TOKEN` environment variable. Note: tokens have expiration times.
* `public_key` - (Optional) PAK Public Key. Can also be set with the `MONGODB_ATLAS_PUBLIC_API_KEY` environment variable.
* `private_key` - (Optional) PAK Private Key. Can also be set with the `MONGODB_ATLAS_PRIVATE_API_KEY` environment variable.
* `base_url` - (Optional) MongoDB Atlas Base URL. Can also be set with the `MONGODB_ATLAS_BASE_URL` environment variable. For advanced use cases, you can configure custom API endpoints:
  ```terraform
  provider "mongodbatlas" {
    client_id     = var.mongodbatlas_client_id
    client_secret = var.mongodbatlas_client_secret
    base_url      = "https://custom-atlas-api.example.com"
  }
  ```
* `realm_base_url` - (Optional) MongoDB Realm Base URL. Can also be set with the `MONGODB_REALM_BASE_URL` environment variable. For advanced use cases, you can configure custom Realm API endpoints:
  ```terraform
  provider "mongodbatlas" {
    client_id       = var.mongodbatlas_client_id
    client_secret   = var.mongodbatlas_client_secret
    realm_base_url  = "https://custom-realm-api.example.com"
  }
  ```
* `is_mongodbgov_cloud` - (Optional) Set to `true` to use MongoDB Atlas for Government, a dedicated deployment option for government agencies and contractors requiring FedRAMP compliance. When enabled, the provider uses government-specific API endpoints. Ensure credentials are created in the government environment. See [Atlas for Government Considerations](https://www.mongodb.com/docs/atlas/government/api/#atlas-for-government-considerations) for feature limitations and requirements.
  ```terraform
  provider "mongodbatlas" {
    client_id           = var.mongodbatlas_client_id
    client_secret       = var.mongodbatlas_client_secret
    is_mongodbgov_cloud = true
  }
  ```
* `assume_role` - (Optional) AWS IAM role configuration for accessing secrets in AWS Secrets Manager. See [AWS Secrets Manager](#aws-secrets-manager) section for details.
* `secret_name` - (Optional) Name of the secret in AWS Secrets Manager.
* `region` - (Optional) AWS region where the secret is stored.
* `aws_access_key_id` - (Optional) AWS Access Key ID. Can also be set with the `AWS_ACCESS_KEY_ID` environment variable.
* `aws_secret_access_key` - (Optional) AWS Secret Access Key. Can also be set with the `AWS_SECRET_ACCESS_KEY` environment variable.
* `aws_session_token` - (Optional) AWS Session Token. Can also be set with the `AWS_SESSION_TOKEN` environment variable.
* `sts_endpoint` - (Optional) AWS STS endpoint. Can also be set with the `STS_ENDPOINT` environment variable.

### Complete Provider Configuration Example

```terraform
provider "mongodbatlas" {
  # Authentication (choose one method)
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret

  # OR use access token
  # access_token = var.mongodbatlas_access_token

  # OR use PAK
  # public_key  = var.mongodbatlas_public_key
  # private_key = var.mongodbatlas_private_key

  # Optional: MongoDB Atlas for Government
  # is_mongodbgov_cloud = true

  # Optional: Custom endpoints
  # base_url      = "https://custom-atlas-api.example.com"
  # realm_base_url = "https://custom-realm-api.example.com"

  # Optional: AWS Secrets Manager configuration
  # assume_role {
  #   role_arn = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts"
  # }
  # secret_name = "mongodbsecret"
  # region      = "us-east-2"
}
```

## Credential Priority and Warnings

If multiple credentials are provided in the same source, the provider displays a warning and uses credentials in this priority order: Access Token → SA → PAK.

Example warning messages:
- "Access Token will be used although Service Account is also set"
- "Service Account will be used although API Key is also set"

## Security Best Practices

- Never hard-code credentials in Terraform configuration files
- Use environment variables or a secrets management system
- Regularly rotate credentials
- Apply the principle of least privilege when assigning roles
- Use Terraform's `sensitive` attribute for credential variables
- Consider Terraform Cloud or Enterprise for secure variable storage

## Troubleshooting

### Authentication Issues

1. **Verify credentials are correctly set**
   ```shell
   # Check environment variables
   echo $MONGODB_ATLAS_CLIENT_ID
   echo $MONGODB_ATLAS_CLIENT_SECRET
   ```

2. **Check provider configuration**
   ```shell
   # Enable debug logging
   export TF_LOG=DEBUG
   terraform plan
   ```

3. **Verify permissions**
   - Ensure your SA or PAK has appropriate organization/project roles
   - Check IP access list configuration if using PAK

4. **Common error messages and solutions**
   - `401 Unauthorized`: Check credentials are correct and not expired
   - `403 Forbidden`: Verify account has necessary permissions
   - `IP not on access list`: Add your IP to the API access list (PAK only)

### MongoDB Atlas for Government Issues

- Ensure `is_mongodbgov_cloud = true` is set in provider configuration
- Verify credentials are from the Government environment, not commercial Atlas
- Check that requested resources are available in Atlas for Government

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