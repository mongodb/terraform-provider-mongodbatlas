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

**IP access list**: Some API resources, such as backup resources, require an IP access list by default. We highly suggest that you add an IP access list as soon as possible.  See [Require IP Access List for Public API](https://docs.atlas.mongodb.com/tutorial/manage-organization-settings/#require-ip-access-list-for-public-api) for more info.

**API Key List**: Some API resources such as Organization API Access List Entries, Cloud Backup Restores, Cloud Backup Snapshots, Cloud Backup Schedules, Legacy Backups, require an API Key list to utilize this feature. See [Resources that require API Key List](https://www.mongodb.com/docs/atlas/configure-api-access/#use-api-resources-that-require-an-access-list)
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
