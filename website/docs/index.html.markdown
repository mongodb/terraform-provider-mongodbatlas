---
layout: "mongodbatlas"
page_title: "Provider: MongoDB Atlas"
sidebar_current: "docs-mongodbatlas-index"
description: |-
  The MongoDB Atlas provider is used to interact with the resources supported by MongoDB Atlas. The provider needs to be configured with the proper credentials before it can be used.
---

# MongoDB Atlas Provider

The MongoDB Atlas provider is used to interact with the resources supported by [MongoDB Atlas](https://www.mongodb.com/cloud/atlas). The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available provider resources and data sources.

You may want to consider pinning the [provider version](https://www.terraform.io/docs/configuration/providers.html#provider-versions) to ensure you have a chance to review and prepare for changes.   Speaking of changes, see [CHANGELOG](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/CHANGELOG.md) for current version information.  

## Example Usage

```hcl
# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  public_key = "${var.mongodbatlas_public_key}"
  private_key  = "${var.mongodbatlas_private_key}"
}

#Create the resources
...
```

## Configure Atlas Programmatic Access

In order to setup authentication with the MongoDB Atlas provider a programmatic API key must be generated for MongoDB Atlas with the appropriate permissions and IP whitelist entries.   The [MongoDB Atlas documentation](https://docs.atlas.mongodb.com/tutorial/manage-programmatic-access/index.html) contains the most up-to-date instructions for creating and managing your key(s) and IP access.   Be aware, not all API resources require an IP access list by default, but one can set Atlas to require IP access entries for all API resources, see the [organization settings documentation](https://docs.atlas.mongodb.com/tutorial/manage-organization-settings/#require-ip-whitelist-for-public-api) for more info.

## Authenticate the Provider

The MongoDB Atlas provider offers a flexible means of providing credentials for authentication. The following methods are supported and explained below:

### Static credentials

Static credentials can be provided by adding the following attributes in-line in the MongoDB Atlas provider block, either directly or via input variable/local value:

Usage:

```hcl
provider "mongodbatlas" {
  public_key = "atlas_public_api_key" #required
  private_key  = "atlas_private_api_key" #required
}
```

~> *IMPORTANT* Hard-coding your MongoDB Atlas programmatic API key pair into a Terraform configuration is not recommended.  Consider the risks, especially the inadvertent submission of a configuration file containing secrets to a public repository.

### Environment variables

You can also provide your credentials via the environment variables, MONGODB_ATLAS_PUBLIC_KEY and MONGODB_ATLAS_PRIVATE_KEY, for your public and private MongoDB Atlas programmatic API key pair respectively:

```hcl
provider "mongodbatlas" {}
```

Usage (prefix the export commands with a space to avoid the keys being recorded in OS history):

```shell
$  export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
$  export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
$ terraform plan
```

## Argument Reference

In addition to [generic `provider`
arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g.
`alias` and `version`), the following arguments are supported in the MongoDB
Atlas `provider` block:

* `public_key` - (Optional) This is the public key of your MongoDB Atlas API key pair. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PUBLIC_KEY`
  environment variable.

* `private_key` - (Optional) This is the private key of your MongoDB Atlas key pair. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PRIVATE_KEY`
  environment variable.

For more information on configuring and managing programmatic API Keys see the [MongoDB Atlas Documentation](https://docs.atlas.mongodb.com/tutorial/manage-programmatic-access/index.html).

## Helpful Links/Information

[Upgrade Guide for Terraform MongoDB Atlas 0.4.0](https://www.mongodb.com/blog/post/upgrade-guide-for-terraform-mongodb-atlas-040)

[MongoDB Atlas and Terraform Landing Page](https://www.mongodb.com/atlas/terraform)

[Report bugs](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)

[Request Features](https://feedback.mongodb.com/forums/924145-atlas?category_id=370723)

[Support](https://docs.atlas.mongodb.com/support/) covered by MongoDB Atlas support plans, Developer and above.
