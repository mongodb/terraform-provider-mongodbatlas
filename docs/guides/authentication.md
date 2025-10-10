---
page_title: "Authentication Configuration"
---

# Authentication Configuration

This guide provides comprehensive information about the authentication methods available in the MongoDB Atlas Provider. The provider supports multiple authentication mechanisms to suit different security requirements and deployment scenarios.

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

## Service Account (Recommended)

Service Accounts simplify authentication by eliminating the need to create new Atlas-specific user identities and permission credentials. See [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/) and [MongoDB Atlas Service Account Limits](https://www.mongodb.com/docs/manual/reference/limits/#mongodb-atlas-service-account-limits) for more information.

### Configuration

**Using provider attributes:**
```terraform
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

**Using environment variables:**
```shell
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**Note:** Service Accounts can't be used with `mongodbatlas_event_trigger` resources as its API doesn't support it yet.

## Service Account Token

Instead of using Client ID and Client Secret, you can generate and use an SA Token directly. See [Generate Service Account Token](https://www.mongodb.com/docs/atlas/api/service-accounts/generate-oauth2-token/#std-label-generate-oauth2-token-atlas) for details. Note that tokens have an expiration time.

### Configuration

**Using provider attributes:**
```terraform
provider "mongodbatlas" {
  access_token = var.mongodbatlas_access_token
}
```

**Using environment variables:**
```shell
export MONGODB_ATLAS_ACCESS_TOKEN="<ATLAS_ACCESS_TOKEN>"
```

**Important:** The MongoDB Terraform provider currently does not support additional Token OAuth features like scopes.

## Programmatic Access Key (PAK)

PAK is the legacy authentication method. You need to generate a Programmatic Access Key with the appropriate [role](https://docs.atlas.mongodb.com/reference/user-roles/). See [MongoDB Atlas documentation](https://www.mongodb.com/docs/atlas/configure-api-access-org/) for instructions on creating and managing your keys.

**Role recommendation:** If unsure of which role level to grant your key, we suggest creating an organization API Key with an Organization Owner role to ensure sufficient access for all actions.

### Configuration

**Using provider attributes:**
```terraform
provider "mongodbatlas" {
  public_key  = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}
```

**Using environment variables:**
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

## MongoDB Atlas for Government

To use the provider with MongoDB Atlas for Government, add the `is_mongodbgov_cloud` parameter:

```terraform
provider "mongodbatlas" {
  client_id           = var.mongodbatlas_client_id
  client_secret       = var.mongodbatlas_client_secret
  is_mongodbgov_cloud = true
}
```

See [Atlas for Government Considerations](https://www.mongodb.com/docs/atlas/government/api/#atlas-for-government-considerations) for more information.

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

## Troubleshooting

For authentication issues:
1. Verify credentials are correctly set in your chosen source
2. Check that the credentials have appropriate permissions
3. Ensure there are no conflicting credentials across different sources
4. Review provider debug logs if issues persist

## Additional Resources

- [MongoDB Atlas API Documentation](https://www.mongodb.com/docs/atlas/api/)
- [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/)
- [Configure API Access](https://www.mongodb.com/docs/atlas/configure-api-access/)