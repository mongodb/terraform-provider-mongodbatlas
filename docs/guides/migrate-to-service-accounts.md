---
page_title: "Migration Guide: Service Accounts"
---

# Migration Guide: Service Accounts

This guide helps you migrate from Programmatic Access Key (PAK) authentication to Service
Accounts (SA) authentication without impacting your deployment. SA simplifies
authentication by eliminating the need to create new Atlas-specific user identities and
permission credentials.

This document contains the following sections:

- [Before You Begin](#before-you-begin)
- [Procedure](#procedure)
  - [Source Hierarchy](#source-hierarchy)
  - [AWS Secrets Manager](#aws-secrets-manager)
  - [Provider Attributes](#provider-attributes)
  - [Environment Variables](#environment-variables)
- [Attribute Names](#attribute-names)

**Note:** For more information on SA, see [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/) in the MongoDB documentation.

## Before You Begin

See [Get Started with the Atlas Administration API](https://www.mongodb.com/docs/atlas/configure-api-access/#grant-programmatic-access-to-an-organization) in the MongoDB Atlas documentation for detailed
instructions on how to create your organization's SA and granting the required access.

## Procedure

The following example shows the variables for PAK authentication:

```terraform
provider "mongodbatlas" {
  public_key  = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}
```

To migrate from PAK authentication to Service Accounts (SA) authentication, change the credentials provided.

Your credentials can come from the following sources:

- [AWS Secrets Manager](#aws-secrets-manager)
- [Provider attributes](#provider-attributes)
- [Environment variables](#environment-variables)

### Source Hierarchy

The source selection criteria applied is the following:

| Criteria (in order) | Source |
|---|---|
| 1. AWSAssumeRoleARN is set in provider | AWS Secrets Manager (getting AWS credentials from provider) |
| 2. AWSAssumeRoleARN is set in environment variable | AWS Secrets Manager (getting AWS credentials from environment variables) |
| 3. Any of AccessToken, ClientID, ClientSecret, PublicKey, PrivateKey in provider are set | Provider |
| 4. Any of AccessToken, ClientID, ClientSecret, PublicKey, PrivateKey in environment variables are set | Environment variables |
| 5. Else | No authentication |

The name of the attributes containing your credentials change depending on their source. See [Attribute Names](#attribute-names) for detailed information on attribute naming conventions.

### AWS Secrets Manager

The following example shows how to set up SA using the AWS Secrets Manager:

```terraform
provider "mongodbatlas" {
  assume_role {
    role_arn = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts"
  }
  secret_name           = "mongodbsecret" // fully qualified secret_name ARN also supported as input 
  region                = "us-east-2"
  aws_access_key_id     = "<AWS_ACCESS_KEY_ID>"
  aws_secret_access_key = "<AWS_SECRET_ACCESS_KEY>"
  aws_session_token     = "<AWS_SESSION_TOKEN>"
  sts_endpoint          = "https://sts.us-east-2.amazonaws.com/"
}
```

### Provider Attributes

You can set up SA through provider attributes by either providing your client ID and secret or
providing a valid access token.

#### Provide your Client ID and Secret (Provider)

The following example shows how to set up SA using the `client_id` and `client_secret` provider attributes:

```terraform
provider "mongodbatlas" {
  client_id = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

**NOTE:** Using `client_id` and `client_secret` is limited to 10 tokens per minute.

#### Provide a Valid Access Token (Provider)

The following example shows how to set up SA with a ``access_token`` attribute:

```terraform
provider "mongodbatlas" { 
  access_token = var.mongodbatlas_access_token
}
```

Consider that the access token is **valid for one hour only**.

See [Generate Service Account Token](https://www.mongodb.com/docs/atlas/api/service-accounts/generate-oauth2-token/#std-label-generate-oauth2-token-atlas) for more details on creating an SA token.

**IMPORTANT:**  Currently, the MongoDB Terraform provider does not support additional Token OAuth features like scopes.

**NOTE:** Service Accounts is not currently supported as the authentication method for the ``mongodbatlas_event_trigger`` resource.

### Environment Variables

You can set up SA through environment variables by either providing your client ID and secret or
providing a valid access token.

#### Provide your Client ID and Secret (Environment Variables)

The following example shows how to set up SA using the ``MONGODB_ATLAS_CLIENT_ID`` and
``MONGODB_ATLAS_CLIENT_SECRET`` variables:

```shell
$ export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"

$ export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"

provider "mongodbatlas" {}
```

#### Provide a Valid Access Token (Environment Variables)

The following example shows how to set up SA with a ``MONGODB_ATLAS_ACCESS_TOKEN`` variable:

```shell
$ export MONGODB_ATLAS_ACCESS_TOKEN="<ATLAS_ACCESS_TOKEN>" 

provider "mongodbatlas" {}
```

## Attribute Names

The following table details the attributes' naming conventions according to each source:

| Attribute Name | Provider | AWS Secrets Manager | Environment variable (checked in order) |
|---|---|---|---|
| AccessToken | `access_token` | `access_token` | `MONGODB_ATLAS_ACCESS_TOKEN`, `MCLI_ACCESS_TOKEN` |
| ClientID | `client_id` | `client_id` | `MONGODB_ATLAS_CLIENT_ID`, `MCLI_CLIENT_ID` |
| ClientSecret | `client_secret` | `client_secret` | `MONGODB_ATLAS_CLIENT_SECRET`, `MCLI_CLIENT_SECRET` |
| PublicKey | `public_key` | `public_key` | `MONGODB_ATLAS_PUBLIC_API_KEY`, `MONGODB_ATLAS_PUBLIC_KEY`, `MCLI_PUBLIC_API_KEY` |
| PrivateKey | `private_key` | `private_key` | `MONGODB_ATLAS_PRIVATE_API_KEY`, `MONGODB_ATLAS_PRIVATE_KEY`, `MCLI_PRIVATE_API_KEY` |
| BaseURL | `base_url` | `base_url` | `MONGODB_ATLAS_BASE_URL`, `MCLI_OPS_MANAGER_URL` |
| RealmBaseURL | `realm_base_url` | `realm_base_url` | `MONGODB_REALM_BASE_URL` |
| AWSAssumeRoleARN | `assume_role.role_arn` |  | `ASSUME_ROLE_ARN`, `TF_VAR_ASSUME_ROLE_ARN` |
| AWSSecretName | `secret_name` |  | `SECRET_NAME`, `TF_VAR_SECRET_NAME` |
| AWSRegion | `region` |  | `AWS_REGION`, `TF_VAR_AWS_REGION` |
| AWSAccessKeyID | `aws_access_key_id` |  | `AWS_ACCESS_KEY_ID`, `TF_VAR_AWS_ACCESS_KEY_ID` |
| AWSSecretAccessKey | `aws_secret_access_key` |  | `AWS_SECRET_ACCESS_KEY`, `TF_VAR_AWS_SECRET_ACCESS_KEY` |
| AWSSessionToken | `aws_session_token` |  | `AWS_SESSION_TOKEN`, `TF_VAR_AWS_SESSION_TOKEN` |
| AWSEndpoint | `sts_endpoint` |  | `STS_ENDPOINT`, `TF_VAR_STS_ENDPOINT` |

