# MongoDB Atlas Provider

You can use the MongoDB Atlas provider to interact with the resources supported by [MongoDB Atlas](https://www.mongodb.com/cloud/atlas).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available provider resources and data sources.

See [CHANGELOG](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/CHANGELOG.md) for current version information.  

## Provider and terraform version constraints

We recommend that you pin your Atlas [provider version](https://developer.hashicorp.com/terraform/language/providers/requirements#version) to at least the [major version](#versioning-strategy) (e.g. `~> 2.0`) to avoid accidental upgrades to incompatible new versions. Starting on `2.0.0`, the [MongoDB Atlas Provider Versioning Policy](#mongodb-atlas-provider-versioning-policy) ensures that minor and patch versions do not include [Breaking Changes](#definition-of-breaking-changes). 

For Terraform version, we recommend that you use the latest [HashiCorp Terraform Core Version](https://github.com/hashicorp/terraform). For more details see [HashiCorp Terraform Version Compatibility Matrix](#hashicorp-terraform-version-compatibility-matrix).

## Authenticate the Provider

The MongoDB Atlas provider offers Service Account (SA) and Programmatic Access Key (PAK) authentication mechanisms.

Credentials can be stored in these sources: AWS Secrets Manager, provider attributes or environment variables. The first source in this order to have credentials will be used.
The following table contains the credential attribute names to use:

| Attribute | AWS Secrets Manager | Provider attribute | Environment variable |
|---|---|---|---|
| SA Client ID | `client_id` | `client_id` | `MONGODB_ATLAS_CLIENT_ID` |
| SA Client Secret | `client_secret` | `client_secret` | `MONGODB_ATLAS_CLIENT_SECRET` |
| Access Token | `access_token` | `access_token` | `MONGODB_ATLAS_ACCESS_TOKEN` |
| PAK Public Key | `public_key` | `public_key` | `MONGODB_ATLAS_PUBLIC_API_KEY`, `MONGODB_ATLAS_PUBLIC_KEY` |
| PAK Private Key | `private_key` | `private_key` | `MONGODB_ATLAS_PRIVATE_API_KEY`, `MONGODB_ATLAS_PRIVATE_KEY` |

~> *IMPORTANT* Hard-coding your MongoDB Atlas SA or PAK key pair into the MongoDB Atlas Provider configuration is not recommended.
Consider the risks, especially the inadvertent submission of a configuration file containing secrets to a public repository.

### Service Account

Service Accounts (SA) is the preferred authentication method for the MongoDB Atlas provider.
The [MongoDB Atlas documentation](https://www.mongodb.com/docs/atlas/configure-api-access/#grant-programmatic-access-to-an-organization) contains the most up-to-date instructions for creating your organization's SA and granting the required access.
These are the [MongoDB Atlas Service Account Limits](https://www.mongodb.com/docs/manual/reference/limits/#mongodb-atlas-service-account-limits).


See [Migration Guide: Service Accounts](guides/migrate-to-service-accounts) if you are currently using Programmatic Access Key and want to move to Service Accounts.

An example using provider attributes is:

```terraform
provider "mongodbatlas" {
  client_id = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

An example using environment variables is:

```shell
$ export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
$ export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

See below for AWS Secret Manager examples.

**IMPORTANT:**  Currently, the MongoDB Terraform provider does not support additional Token OAuth features like scopes.

**NOTE:** Service Accounts can't be used with `mongodbatlas_event_trigger` resources as its API doesn't support it yet.

### Service Account Token

Instead of using your Client ID and Client Secret, you can alternatively generate the SA Token yourself, see [Generate Service Account Token](https://www.mongodb.com/docs/atlas/api/service-accounts/generate-oauth2-token/#std-label-generate-oauth2-token-atlas) for more details on creating an SA token. Note that tokens has an expiration time.

An example using provider attributes is:

```terraform
provider "mongodbatlas" {
  access_token = var.mongodbatlas_access_token
}
```

An example using environment variables is:

```shell
$ export MONGODB_ATLAS_ACCESS_TOKEN="<ATLAS_ACCESS_TOKEN>"
```

**IMPORTANT:**  Currently, the MongoDB Terraform provider does not support additional Token OAuth features like scopes.

### Programmatic Access Key

You have to generate a Programmatic Access Key (PAK) with the appropriate [role](https://docs.atlas.mongodb.com/reference/user-roles/).
The [MongoDB Atlas documentation](https://www.mongodb.com/docs/atlas/configure-api-access-org/) contains the most up-to-date instructions for creating and managing your key(s), setting the appropriate role, and optionally configuring IP access.

**Role**: If unsure of which role level to grant your key, we suggest creating an organization API Key with an Organization Owner role. This ensures that you have sufficient access for all actions.

An example using provider attributes is:

```terraform
provider "mongodbatlas" {
  public_key = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}
```

An example using environment variables is:

```shell
$ export MONGODB_ATLAS_PUBLIC_API_KEY="<ATLAS_PUBLIC_API_KEY>"
$ export MONGODB_ATLAS_PRIVATE_API_KEY="<ATLAS_PRIVATE_API_KEY>"
```

**NOTE:** We recommend that you use `MONGODB_ATLAS_PUBLIC_API_KEY` and `MONGODB_ATLAS_PRIVATE_API_KEY` because they are compatible with other MongoDB tools, such as Atlas CLI. You can still use `MONGODB_ATLAS_PUBLIC_KEY` and `MONGODB_ATLAS_PRIVATE_KEY` as alternative keys in your local environment. However, these environment variables are not guaranteed to work across all tools in the MongoDB ecosystem.

See below for AWS Secret Manager examples.
 
### AWS Secrets Manager

AWS Secrets Manager (AWS SM) helps to manage, retrieve, and rotate SAs, PAKs, database credentials, and other secrets throughout their lifecycles. See [product page](https://aws.amazon.com/secrets-manager/) and [documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/what-is-systems-manager.html) for more details.

You can configure AWS credentials to access AWS Secrets Manager using environment variables or provider attributes:

| Attribute | Provider attribute | Environment variable |
|---|---|---|
| Assume Role ARN | `assume_role.role_arn` | `ASSUME_ROLE_ARN` |
| Secret Name | `secret_name` | `SECRET_NAME` |
| AWS Region | `region` | `AWS_REGION` |
| Access Key ID | `aws_access_key_id` | `AWS_ACCESS_KEY_ID` |
| Secret Access Key | `aws_secret_access_key` |  `AWS_SECRET_ACCESS_KEY` |
| Session Token | `aws_session_token` | `AWS_SESSION_TOKEN` |
| STS Endpoint | `sts_endpoint` | `STS_ENDPOINT` |

In order to enable the Terraform MongoDB Atlas Provider with AWS SM, please follow the below steps: 

1. Create a SA or PAK and add them as one secret to AWS SM with a raw value. Take note of which AWS Region secret is being stored in. Each attribute needs to be entered as their own key value pair. See below example for SA:

``` 
{
  "client_id": "secret1",
  "client_secret":"secret2"
}
```

And this is an example for PAK:

``` 
{
  "public_key": "secret3",
  "private_key":"secret4"
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
export AWS_ACCESS_KEY_ID='<AWS_ACCESS_KEY_ID>'
export AWS_SECRET_ACCESS_KEY='<AWS_SECRET_ACCESS_KEY>'
```

4. In terminal, use the AWS CLI command: `aws sts assume-role --role-arn ROLE_ARN_FROM_ABOVE --role-session-name newSession` 

Note: AWS STS secrets are short lived by default, use the ` --duration-seconds` flag to specify longer duration as needed.

5. Store each of the 3 new created secrets from AWS STS as environment variables (hardcoding secrets into config file with additional risk is also supported). For example:

```
export AWS_ACCESS_KEY_ID='<AWS_ACCESS_KEY_ID>'
export AWS_SECRET_ACCESS_KEY='<AWS_SECRET_ACCESS_KEY>'
export AWS_SESSION_TOKEN="<AWS_SESSION_TOKEN>"
```

6. Add assume_role block with `role_arn`, `secret_name`, and AWS `region` where secret is stored as part of AWS SM. Each of these 3 fields are REQUIRED. For example:

```terraform
# Configure the MongoDB Atlas Provider to Authenticate with AWS Secrets Manager 
provider "mongodbatlas" {
  assume_role {
    role_arn = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/mdbsts"
  }
  secret_name           = "mongodbsecret"
  // fully qualified secret_name ARN also supported as input "arn:aws:secretsmanager:af-south-1:<AWS_ACCOUNT_ID>:secret:test789-TO06Hy" 
  region                = "us-east-2"
  
  aws_access_key_id     = "<AWS_ACCESS_KEY_ID>"
  aws_secret_access_key = "<AWS_SECRET_ACCESS_KEY>"
  aws_session_token     = "<AWS_SESSION_TOKEN>"
  sts_endpoint          = "https://sts.us-east-2.amazonaws.com/"
}
```

Note: `aws_access_key_id`, `aws_secret_access_key`, and `aws_session_token` can also be passed in using environment variables i.e. aws_access_key_id will accept AWS_ACCESS_KEY_ID and TF_VAR_AWS_ACCESS_KEY_ID as a default value in place of value in a terraform file variable. 

Note: Fully qualified `secret_name` ARN as input is REQUIRED for cross-AWS account secrets. For more detatils see:
* https://aws.amazon.com/blogs/security/how-to-access-secrets-across-aws-accounts-by-attaching-resource-based-policies/ 
* https://aws.amazon.com/premiumsupport/knowledge-center/secrets-manager-share-between-accounts/

Note: `sts_endpoint` parameter is REQUIRED for cross-AWS region or cross-AWS account secrets. 

7. In terminal, `terraform init`

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

## MongoDB Atlas Provider Versioning Policy

In order to promote stability, predictability, and transparency, the MongoDB Atlas Terraform Provider will implement **semantic versioning** with a **scheduled release cadence**. Our goal is to deliver regular improvements to the provider without overburdening users with frequent breaking changes.

---

### Definition of Breaking Changes

Our definition of breaking changes aligns with the impact updates have on the customer:

Breaking changes are defined as any change that requires user intervention to address.
This may include:

- Modifying existing schema (e.g., removing or renaming fields, renaming resources)
- Changes to business logic (e.g., implicit default values or server-side behavior)
- Provider-level changes (e.g., changing retry behavior)

Final confirmation of a breaking change — possibly leading to an exemption — is subject to:

- MongoDB’s understanding of the adoption level of the feature
- Timing of the next planned major release
- The change's relation to a bug fix

---

### Versioning Strategy

We follow [semantic versioning](https://semver.org/) for all updates:

- **Major (X.0.0):** Introduces breaking changes (as defined by MongoDB)
- **Minor (X.Y.0):** Adds non-breaking changes and announces deprecations
- **Patch (X.Y.Z):** Includes bug fixes and documentation updates

We do not utilize pre-release versioning at this time.

---

### Release Cadence

To minimize unexpected changes, we follow a scheduled cadence:

- **Minor and patch** versions follow a **biweekly** release pattern
- **Major** versions are released **once per year**, with a maximum of **two per calendar year**
- The provider team may adjust the schedule based on need

**Off-cycle releases** may occur for critical security flaws or regressions.

---

### Deprecation Policy

We use a structured deprecation window to notify customers in advance:

- Breaking changes are **deprecated in a minor version** with:
  - Warnings in migration guides, changelogs, and resource usage
- Deprecated functionality is **removed in the next 1–2 major versions**, unless otherwise stated

---

### Customer Communication

We are committed to clear and proactive communication:

- **Each release** includes a [changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/releases) clearly labeling:
  - `breaking`, `deprecated`, `bug-fix`, `feature`, and `enhancement` changes
- **Major versions** include migration guides
- **Minor and patch versions** generally do not include migration guides, but may if warranted
- **GitHub tags** with `vX.Y.Z` format are provided for all releases

---

## [HashiCorp Terraform Version](https://www.terraform.io/downloads.html) Compatibility Matrix

<!-- DO NOT remove below placeholder comments as this table is auto-generated -->
<!-- MATRIX_PLACEHOLDER_START -->
| HashiCorp Terraform Release | HashiCorp Terraform Release Date  | HashiCorp Terraform Full Support End Date  | MongoDB Atlas Support End Date |
|:-------:|:------------:|:------------:|:------------:|
| 1.13.x | 2025-08-20 | 2027-08-31 | 2027-08-31 |
| 1.12.x | 2025-05-14 | 2027-05-31 | 2027-05-31 |
| 1.11.x | 2025-02-27 | 2027-02-28 | 2027-02-28 |
| 1.10.x | 2024-11-27 | 2026-11-30 | 2026-11-30 |
| 1.9.x | 2024-06-26 | 2026-06-30 | 2026-06-30 |
| 1.8.x | 2024-04-10 | 2026-04-30 | 2026-04-30 |
| 1.7.x | 2024-01-17 | 2026-01-31 | 2026-01-31 |
| 1.6.x | 2023-10-04 | 2025-10-31 | 2025-10-31 |
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
We have [example configurations](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.1/examples)
in our GitHub repo that will help both beginner and more advanced users.

Have a good example you've created and want to share?
Let us know the details via an [issue](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)
or submit a PR of your work to add it to the `examples` directory in our [GitHub repo](https://github.com/mongodb/terraform-provider-mongodbatlas/).

## Terraform MongoDB Atlas Modules
You can now leverage our [Terraform Modules](https://registry.terraform.io/namespaces/terraform-mongodbatlas-modules) to easily get started with MongoDB Atlas and critical features like [Push-based log export](https://registry.terraform.io/modules/terraform-mongodbatlas-modules/push-based-log-export/mongodbatlas/latest), [Private Endpoints](https://registry.terraform.io/modules/terraform-mongodbatlas-modules/private-endpoint/mongodbatlas/latest), etc.
