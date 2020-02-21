---
layout: "mongodbatlas"
page_title: "Provider: MongoDB Atlas"
sidebar_current: "docs-mongodbatlas-index"
description: |-
  The MongoDB Atlas provider is used to interact with the resources supported by MongoDB Atlas Services. The provider needs to be configured with the proper credentials before it can be used.
---

# MongoDB Atlas Provider

The MongoDB Atlas provider is used to interact with the resources supported by MongoDB Atlas Services. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage
```hcl
# Configure the MongoDB Atlas Provider
provider "mongodbatlas" {
  public_key = "${var.mongodbatlas_public_key}"
  private_key  = "${var.mongodbatlas_private_key}"
}

#Create the resources
```

## Authentication

The MongoDB Atlas provider offers a flexible means of providing credentials for authentication. The following methods are supported, in this order, and explained below:

### Static credentials

Static credentials can be provided by adding the following attributes in-line in the MongoDB Atlas provider block:

Usage:

```hcl
provider "mongodbatlas" {
  public_key = "" #required
  private_key  = "" #required
}
```

### Environment variables

You can provide your credentials via environment variables, representing your MongoDB Atlas
authentication.

```hcl
provider "mongodbatlas" {}
```

Usage:

```shell
$ export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
$ export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
$ terraform plan
```

## Argument Reference

In addition to [generic `provider`
arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g.
`alias` and `version`), the following arguments are supported in the MongoDB
Atlas `provider` block:

* `public_key` - (Optional) This is the MongoDB Atlas API public_key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PUBLIC_KEY`
  environment variable.

* `private_key` - (Optional) This is the MongoDB Atlas private_key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PRIVATE_KEY`
  environment variable.

For more information about how to get these programmatic API Keys see the following [link](https://docs.atlas.mongodb.com/configure-api-access/#manage-programmatic-access-to-an-organization).

## Helpful Links/Information

[Upgrade Guide for Terraform MongoDB Atlas 0.4.0](https://www.mongodb.com/blog/post/upgrade-guide-for-terraform-mongodb-atlas-040)

[MongoDB Atlas and Terraform Landing Page](https://www.mongodb.com/atlas/terraform)

[Report bugs](https://github.com/terraform-providers/terraform-provider-mongodbatlas/issues)

[Request Features](https://feedback.mongodb.com/forums/924145-atlas?category_id=370723)

[Support](https://docs.atlas.mongodb.com/support/) covered by MongoDB Atlas support plans, Developer and above.
