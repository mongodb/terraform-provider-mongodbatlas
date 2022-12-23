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

In order to set up authentication with the MongoDB Atlas provider, you must generate a programmatic API key for MongoDB Atlas with the appropriate [role](https://docs.atlas.mongodb.com/reference/user-roles/) and IP access list entries.
The [MongoDB Atlas documentation](https://docs.atlas.mongodb.com/tutorial/manage-programmatic-access/index.html) contains the most up-to-date instructions for creating and managing your key(s), setting the appropriate role, and IP access.  

**Role**: If unsure of which role level to grant your key, we suggest creating an organization API Key with an Organization Owner role. This ensures that you have sufficient access for all actions.

**API Key Access List**: Some Atlas API resources such as Cloud Backup Restores, Cloud Backup Snapshots, and Cloud Backup Schedules **require** an Atlas API Key Access List to utilize these feature.  Hence, if using Terraform, or any other programmatic control, to manage these resources you must have the IP address or CIDR block that the connection is coming from added to the Atlas API Key Access List of the Atlas API key you are using.   See [Resources that require API Key List](https://www.mongodb.com/docs/atlas/configure-api-access/#use-api-resources-that-require-an-access-list)

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

In order to enable the Terraform MongoDB Atlas Provider to use AWS SM, first create Atlas API Keys and add them as a secret to AWS SM with a basic key with a raw value. See below example:  
``` 
     {
      "public_key": "iepubky",
      "private_key":"prvkey"
     }
```

Next, add assume_role block with `role_arn`, `secret_name`, and AWS `region` to match the AWS region where secret is stored with AWS SM. See below example:
```terraform
# Configure the MongoDB Atlas Provider to Authenticate with AWS Secrets Manager 
provider "mongodbatlas" {
  assume_role {
    role_arn = "arn:aws:iam::476xxx451:role/mdbsts"
  }
  secret_name           = "mongodbsecret"
  aws_access_key_id     = "ASIXXBNEK"
  aws_secret_access_key = "ZUZgVb8XYZWEXXEDURGFHFc5Au"
  aws_session_token     = "IQoXX3+Q="
  region                = "us-east-2"
  sts_endpoint          = "https://sts.us-east-2.amazonaws.com/"
}
```
Note: `aws_access_key_id`, `aws_secret_access_key`, `aws_session_token`, `region` can also be passed in using environment variables i.e. aws_access_key_id will accept AWS_ACCESS_KEY_ID and TF_VAR_AWS_ACCESS_KEY_ID as a default value in place of value in a terraform file variable.

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

## Helpful Links/Information

[Upgrade Guide for Terraform MongoDB Atlas 0.4.0](https://www.mongodb.com/blog/post/upgrade-guide-for-terraform-mongodb-atlas-040)

[MongoDB Atlas and Terraform Landing Page](https://www.mongodb.com/atlas/terraform)

[Report bugs](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)

[Request Features](https://feedback.mongodb.com/forums/924145-atlas?category_id=370723)

[Support](https://docs.atlas.mongodb.com/support/) covered by MongoDB Atlas support plans, Developer and above.

## Examples from MongoDB and the Community

We have [example configurations](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples)
in our GitHub repo that will help both beginner and more advanced users.

Have a good example you've created and want to share?
Let us know the details via an [issue](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)
or submit a PR of your work to add it to the `examples` directory in our [GitHub repo](https://github.com/mongodb/terraform-provider-mongodbatlas/).
