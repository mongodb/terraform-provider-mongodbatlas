---
layout: "mongodbatlas"
page_title: "Provider: MongoDB Atlas"
sidebar_current: "docs-mongodbatlas-index"
description: |-
  The MongoDB Atlas provider is used to interact with the resources supported by MongoDB Atlas. The provider needs to be configured with the proper credentials before it can be used.
---

# MongoDB Atlas Provider

You can use the MongoDB Atlas provider to interact with the resources supported by [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available provider resources and data sources.

You may want to consider pinning the [provider version](https://www.terraform.io/docs/configuration/providers.html#provider-versions) to ensure you have a chance to review and prepare for changes.
Speaking of changes, see [CHANGELOG](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/CHANGELOG.md) for current version information.  

## Example Usage

```terraform
# Configure the MongoDB Atlas Provider 
provider "mongodbatlas" {
  public_key = var.mongodbatlas_public_key
  private_key  = var.mongodbatlas_private_key
}
# Create the resources
```
## Configure Atlas Programmatic Access

In order to set up authentication with the MongoDB Atlas provider, you must generate a programmatic API key for MongoDB Atlas with the appropriate [role](https://docs.atlas.mongodb.com/reference/user-roles/).
The [MongoDB Atlas documentation](https://docs.atlas.mongodb.com/tutorial/manage-programmatic-access/index.html) contains the most up-to-date instructions for creating and managing your key(s), setting the appropriate role, and optionally configuring IP access.

**Role**: If unsure of which role level to grant your key, we suggest creating an organization API Key with an Organization Owner role. This ensures that you have sufficient access for all actions.

## Configure MongoDB Atlas for Government

In order to enable the Terraform MongoDB Atlas Provider for use with MongoDB Atlas for Government add is_mongodbgov_cloud = true to your provider configuration:
```terraform
# Configure the MongoDB Atlas Provider for MongoDB Atlas for Government
provider "mongodbatlas" {
  public_key = var.mongodbatlas_public_key
  private_key  = var.mongodbatlas_private_key
  is_mongodbgov_cloud = true
}
# Create the resources
```
Also see [`Atlas for Government Considerations`](https://www.mongodb.com/docs/atlas/government/api/#atlas-for-government-considerations).  

## Authenticate the Provider

The MongoDB Atlas provider offers a flexible means of providing credentials for authentication.
You can use any the following methods:

### Environment Variables

You can also provide your credentials via the environment variables, 
`MONGODB_ATLAS_PUBLIC_KEY` and `MONGODB_ATLAS_PRIVATE_KEY`,
for your public and private MongoDB Atlas programmatic API key pair respectively:

```terraform
provider "mongodbatlas" {}
```

Usage (prefix the export commands with a space to avoid the keys being recorded in OS history):

```shell
$  export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
$  export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
$ terraform plan
```

As an alternative to `MONGODB_ATLAS_PUBLIC_KEY` and `MONGODB_ATLAS_PRIVATE_KEY`
if you are using [MongoDB CLI](https://docs.mongodb.com/mongocli/stable/) 
then `MCLI_PUBLIC_API_KEY` and `MCLI_PRIVATE_API_KEY` are also supported.

### AWS Secrets Manager
AWS Secrets Manager (AWS SM) helps to manage, retrieve, and rotate database credentials, API keys, and other secrets throughout their lifecycles. See [product page](https://aws.amazon.com/secrets-manager/) and [documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/what-is-systems-manager.html) for more details.

In order to enable the Terraform MongoDB Atlas Provider with AWS SM, please follow the below steps: 

1. Create Atlas API Keys and add them as one secret to AWS SM with a raw value. Take note of which AWS Region secret is being stored in. Public Key and Private Key each need to be entered as their own key value pair. See below example:  
``` 
     {
      "public_key": "secret1",
      "private_key":"secret2"
     }
```
2. Create an AWS IAM Role to attach to the AWS STS (Security Token Service) generated short lived API keys. This is required since STS generated API Keys by default have restricted permissions and need to have their permissions elevated in order to authenticate with Terraform. Take note of Role ARN and ensure IAM Role has permission for “sts:AssumeRole”. For example: 
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Statement1",
            "Effect": "Allow",
            "Principal": {
                "AWS": "*"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```
In addition, you are required to also attach the AWS Managed policy of `SecretsManagerReadWrite` to this IAM role.

Note: this policy may be overly broad for many use cases, feel free to adjust accordingly to your organization's needs.

3. In terminal, store as environmental variables AWS API Keys (while you can also hardcode in config files these will then be stored as plain text in .tfstate file and should be avoided if possible). For example:
``` 
export AWS_ACCESS_KEY_ID="secret"
export AWS_SECRET_ACCESS_KEY="secret”
```
4. In terminal, use the AWS CLI command: `aws sts assume-role --role-arn ROLE_ARN_FROM_ABOVE --role-session-name newSession` 

Note: AWS STS secrets are short lived by default, use the ` --duration-seconds` flag to specify longer duration as needed 

5. Store each of the 3 new created secrets from AWS STS as environment variables (hardcoding secrets into config file with additional risk is also supported). For example: 
```
export AWS_ACCESS_KEY_ID="ASIAYBYSK3S5FZEKLETV"
export AWS_SECRET_ACCESS_KEY="lgT6kL9lr1fxM6mCEwJ33MeoJ1M6lIzgsiW23FGH"
export AWS_SESSION_TOKEN="IQoXX3+Q"
```

6. Add assume_role block with `role_arn`, `secret_name`, and AWS `region` where secret is stored as part of AWS SM. Each of these 3 fields are REQUIRED. For example:
```terraform
# Configure the MongoDB Atlas Provider to Authenticate with AWS Secrets Manager 
provider "mongodbatlas" {
  assume_role {
    role_arn = "arn:aws:iam::476xxx451:role/mdbsts"
  }
  secret_name           = "mongodbsecret"
  // fully qualified secret_name ARN also supported as input "arn:aws:secretsmanager:af-south-1:553552370874:secret:test789-TO06Hy" 
  region                = "us-east-2"
  
  aws_access_key_id     = "ASIXXBNEK"
  aws_secret_access_key = "ZUZgVb8XYZWEXXEDURGFHFc5Au"
  aws_session_token     = "IQoXX3+Q="
  sts_endpoint          = "https://sts.us-east-2.amazonaws.com/"
}
```
Note: `aws_access_key_id`, `aws_secret_access_key`, and `aws_session_token` can also be passed in using environment variables i.e. aws_access_key_id will accept AWS_ACCESS_KEY_ID and TF_VAR_AWS_ACCESS_KEY_ID as a default value in place of value in a terraform file variable. 

Note: Fully qualified `secret_name` ARN as input is REQUIRED for cross-AWS account secrets. For more detatils see:
* https://aws.amazon.com/blogs/security/how-to-access-secrets-across-aws-accounts-by-attaching-resource-based-policies/ 
* https://aws.amazon.com/premiumsupport/knowledge-center/secrets-manager-share-between-accounts/

Note: `sts_endpoint` parameter is REQUIRED for cross-AWS region or cross-AWS account secrets. 

7. In terminal, `terraform init` 

### Static Credentials

Static credentials can be provided by adding the following attributes in-line in the MongoDB Atlas provider block, 
either directly or via input variable/local value:

```terraform
provider "mongodbatlas" {
  public_key = "atlas_public_api_key" #required
  private_key  = "atlas_private_api_key" #required
}
```

~> *IMPORTANT* Hard-coding your MongoDB Atlas programmatic API key pair into a Terraform configuration is not recommended.
Consider the risks, especially the inadvertent submission of a configuration file containing secrets to a public repository.

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the MongoDB Atlas `provider` supports the following arguments:

* `public_key` - (Optional) This is the public key of your MongoDB Atlas API key pair. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PUBLIC_KEY` or `MCLI_PUBLIC_API_KEY`
  environment variable.

* `private_key` - (Optional) This is the private key of your MongoDB Atlas key pair. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PRIVATE_KEY` or `MCLI_PRIVATE_API_KEY`
  environment variable.

For more information on configuring and managing programmatic API Keys see the [MongoDB Atlas Documentation](https://docs.atlas.mongodb.com/tutorial/manage-programmatic-access/index.html).

## [HashiCorp Terraform Version](https://www.terraform.io/downloads.html) Compatibility Matrix

<!-- DO NOT remove below placeholder comments as this table is auto-generated -->
<!-- MATRIX_PLACEHOLDER_START -->
| HashiCorp Terraform Release | HashiCorp Terraform Release Date  | HashiCorp Terraform Full Support End Date  | MongoDB Atlas Support End Date |
|:-------:|:------------:|:------------:|:------------:|
| 1.9.x | 2024-06-26 | 2026-06-30 | 2026-06-30 |
| 1.8.x | 2024-04-10 | 2026-04-30 | 2026-04-30 |
| 1.7.x | 2024-01-17 | 2026-01-31 | 2026-01-31 |
| 1.6.x | 2023-10-04 | 2025-10-31 | 2025-10-31 |
| 1.5.x | 2023-06-12 | 2025-06-30 | 2025-06-30 |
| 1.4.x | 2023-03-08 | 2025-03-31 | 2025-03-31 |
| 1.3.x | 2022-09-21 | 2024-09-30 | 2024-09-30 |
<!-- MATRIX_PLACEHOLDER_END -->
For the safety of our users, we require only consuming versions of HashiCorp Terraform that are currently receiving Security / Maintenance Updates. For more details see [Support Period and End-of-Life (EOL) Policy](https://support.hashicorp.com/hc/en-us/articles/360021185113-Support-Period-and-End-of-Life-EOL-Policy).   

HashiCorp Terraform versions that are not listed on this table are no longer supported by MongoDB Atlas. For latest HashiCorp Terraform versions see [here](https://endoflife.date/terraform ).

## Supported OS and Architectures
As per [HashiCorp's recommendations](https://developer.hashicorp.com/terraform/registry/providers/os-arch), we fully support the following operating system / architecture combinations:
- Darwin / AMD64
- Darwin / ARMv8
- Linux / AMD64
- Linux / ARMv8 (sometimes referred to as AArch64 or ARM64)
- Linux / ARMv6
- Windows / AMD64

We ship binaries but do not prioritize fixes for the following operating system / architecture combinations:
- Linux / 386
- Windows / 386
- FreeBSD / 386
- FreeBSD / AMD64

## Helpful Links/Information

[Upgrade Guide for Terraform MongoDB Atlas 0.4.0](https://www.mongodb.com/blog/post/upgrade-guide-for-terraform-mongodb-atlas-040)

[MongoDB Atlas and Terraform Landing Page](https://www.mongodb.com/atlas/terraform)

[Report bugs](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)

[Request Features](https://feedback.mongodb.com/forums/924145-atlas?category_id=370723)

[Support covered by MongoDB Atlas support plans, Developer and above](https://docs.atlas.mongodb.com/support/) 

## Examples from MongoDB and the Community

<!-- NOTE: the below examples link is updated during the release process, when doing changes in the following sentence verify scripts/update-examples-reference-in-docs.sh is not impacted-->
We have [example configurations](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.17.3/examples)
in our GitHub repo that will help both beginner and more advanced users.

Have a good example you've created and want to share?
Let us know the details via an [issue](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)
or submit a PR of your work to add it to the `examples` directory in our [GitHub repo](https://github.com/mongodb/terraform-provider-mongodbatlas/).
