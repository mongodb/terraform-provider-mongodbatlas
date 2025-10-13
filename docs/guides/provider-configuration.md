---
page_title: "Provider Configuration"
---

# Provider Configuration

This guide provides comprehensive information about configuring the MongoDB Atlas Provider, including authentication methods, environment configuration, and special deployment scenarios.

## Authentication Methods

The MongoDB Atlas provider supports the following authentication methods, listed in order of preference:

1. **Service Account (SA)** - Recommended method
2. **Service Account Token** - Alternative SA method with direct token usage
3. **Programmatic Access Key (PAK)** - Legacy method

Credentials can be provided through:
- AWS Secrets Manager (highest priority)
- Provider attributes
- Environment variables

The provider will use the first available credentials source in the order listed above.

### Service Account (Recommended)

Service Accounts simplify authentication by eliminating the need to create new Atlas-specific user identities and permission credentials. See [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/) and [MongoDB Atlas Service Account Limits](https://www.mongodb.com/docs/manual/reference/limits/#mongodb-atlas-service-account-limits) for more information.

Using provider attributes:
```terraform
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

Using environment variables:
```shell
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**Note:** Service Accounts can't be used with `mongodbatlas_event_trigger` resources as its API doesn't support it yet.

### Service Account Token

Instead of using Client ID and Client Secret, you can generate and use an SA Token directly. See [Generate Service Account Token](https://www.mongodb.com/docs/atlas/api/service-accounts/generate-oauth2-token/#std-label-generate-oauth2-token-atlas) for details. Note that tokens have an expiration time.

Using provider attributes:
```terraform
provider "mongodbatlas" {
  access_token = var.mongodbatlas_access_token
}
```

Using environment variables:
```shell
export MONGODB_ATLAS_ACCESS_TOKEN="<ATLAS_ACCESS_TOKEN>"
```

**Important:** The MongoDB Terraform provider currently does not support additional Token OAuth features like scopes.

### Programmatic Access Key (PAK)

PAK is the legacy authentication method. You need to generate a Programmatic Access Key with the appropriate [role](https://docs.atlas.mongodb.com/reference/user-roles/). See [MongoDB Atlas documentation](https://www.mongodb.com/docs/atlas/configure-api-access-org/) for instructions on creating and managing your keys.

**Role recommendation:** If unsure of which role level to grant your key, we suggest creating an organization API Key with an Organization Owner role to ensure sufficient access for all actions.

Using provider attributes:
```terraform
provider "mongodbatlas" {
  public_key  = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}
```

Using environment variables:
```shell
export MONGODB_ATLAS_PUBLIC_API_KEY="<ATLAS_PUBLIC_API_KEY>"
export MONGODB_ATLAS_PRIVATE_API_KEY="<ATLAS_PRIVATE_API_KEY>"
```

## AWS Secrets Manager

AWS Secrets Manager helps manage, retrieve, and rotate credentials throughout their lifecycles. See [AWS Secrets Manager documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/what-is-systems-manager.html) for more details.

### Setup Instructions

1. **Create secrets in AWS Secrets Manager**

   For Service Account:
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

- For cross-account secrets, use fully qualified ARN for `secret_name`
- For cross-region or cross-account access, `sts_endpoint` parameter is required

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

## MongoDB Atlas for Government

MongoDB Atlas for Government is a dedicated deployment option for government agencies and contractors requiring FedRAMP compliance. To use the provider with MongoDB Atlas for Government, add the `is_mongodbgov_cloud` parameter to your provider configuration.

### Configuration

```terraform
provider "mongodbatlas" {
  client_id           = var.mongodbatlas_client_id
  client_secret       = var.mongodbatlas_client_secret
  is_mongodbgov_cloud = true
}
```

### Important Considerations

- MongoDB Atlas for Government uses different API endpoints than standard MongoDB Atlas
- Ensure your credentials are created in the MongoDB Atlas for Government environment
- Not all features available in standard MongoDB Atlas may be available in the Government environment

See [Atlas for Government Considerations](https://www.mongodb.com/docs/atlas/government/api/#atlas-for-government-considerations) for detailed information about limitations and requirements.

## Custom API Endpoints

For advanced use cases, you can configure custom API endpoints:

```terraform
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
  base_url      = "https://custom-atlas-api.example.com"
  realm_base_url = "https://custom-realm-api.example.com"
}
```

## Migration from PAK to Service Account

If you're currently using Programmatic Access Keys and want to migrate to Service Accounts:

### Environment Variables Migration

Replace PAK environment variables:
```shell
# Remove these
unset MONGODB_ATLAS_PUBLIC_API_KEY
unset MONGODB_ATLAS_PRIVATE_API_KEY

# Add these
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

### Provider Attributes Migration

Update your provider configuration:
```terraform
# Old configuration (remove)
provider "mongodbatlas" {
  public_key  = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}

# New configuration (add)
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

### AWS Secrets Manager Migration

Update your secret in AWS Secrets Manager:
- Remove `public_key` and `private_key`
- Add `client_id` and `client_secret`

After making these changes, run `terraform plan` to verify everything is working correctly.

## Provider Configuration Reference

### Provider Arguments

* `client_id` - (Optional) Service Account Client ID. Can also be set with the `MONGODB_ATLAS_CLIENT_ID` environment variable.
* `client_secret` - (Optional) Service Account Client Secret. Can also be set with the `MONGODB_ATLAS_CLIENT_SECRET` environment variable.
* `access_token` - (Optional) Service Account Access Token. Can also be set with the `MONGODB_ATLAS_ACCESS_TOKEN` environment variable. Note: tokens have expiration times.
* `public_key` - (Optional) MongoDB Atlas Programmatic Access Key Public Key. Can also be set with the `MONGODB_ATLAS_PUBLIC_API_KEY` environment variable.
* `private_key` - (Optional) MongoDB Atlas Programmatic Access Key Private Key. Can also be set with the `MONGODB_ATLAS_PRIVATE_API_KEY` environment variable.
* `base_url` - (Optional) MongoDB Atlas Base URL. Can also be set with the `MONGODB_ATLAS_BASE_URL` environment variable.
* `realm_base_url` - (Optional) MongoDB Realm Base URL. Can also be set with the `MONGODB_REALM_BASE_URL` environment variable.
* `is_mongodbgov_cloud` - (Optional) Set to true to use MongoDB Atlas for Government.
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

  # OR use PAK (legacy)
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

### Environment Variables Reference

| Provider Attribute | Environment Variable | Description |
|---|---|---|
| `client_id` | `MONGODB_ATLAS_CLIENT_ID` | Service Account Client ID |
| `client_secret` | `MONGODB_ATLAS_CLIENT_SECRET` | Service Account Client Secret |
| `access_token` | `MONGODB_ATLAS_ACCESS_TOKEN` | Service Account Access Token |
| `public_key` | `MONGODB_ATLAS_PUBLIC_API_KEY` | PAK Public Key |
| `private_key` | `MONGODB_ATLAS_PRIVATE_API_KEY` | PAK Private Key |
| `base_url` | `MONGODB_ATLAS_BASE_URL` | Custom Atlas API endpoint |
| `realm_base_url` | `MONGODB_REALM_BASE_URL` | Custom Realm API endpoint |
| `assume_role.role_arn` | `ASSUME_ROLE_ARN` | AWS IAM Role ARN |
| `secret_name` | `SECRET_NAME` | AWS Secrets Manager secret name |
| `region` | `AWS_REGION` | AWS region |
| `aws_access_key_id` | `AWS_ACCESS_KEY_ID` | AWS Access Key ID |
| `aws_secret_access_key` | `AWS_SECRET_ACCESS_KEY` | AWS Secret Access Key |
| `aws_session_token` | `AWS_SESSION_TOKEN` | AWS Session Token |
| `sts_endpoint` | `STS_ENDPOINT` | AWS STS endpoint |

## Credential Priority and Warnings

If multiple credentials are provided in the same source, the provider will:
1. Display a warning about multiple credentials
2. Use credentials in this priority order: Access Token → Service Account → PAK

Example warning messages:
- "Access Token will be used although Service Account is also set"
- "Service Account will be used although API Key is also set"

## Security Best Practices

- **Never hard-code credentials** in your Terraform configuration files
- Use environment variables or a secrets management system
- Regularly rotate your credentials
- Use the principle of least privilege when assigning roles
- Consider the risks of inadvertently committing secrets to version control
- Use Terraform's `sensitive` attribute for credential variables
- Consider using Terraform Cloud or Enterprise for secure variable storage

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
   - Ensure your Service Account or PAK has appropriate organization/project roles
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