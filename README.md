---
layout: "mongodbatlas"
page_title: "Provider: MongoDB Atlas"
sidebar_current: "docs-mongodbatlas"
description: |-
  The MongoDB Atlas provider is used to interact with the resources supported by MongoDB Altas Services. The provider needs to be configured with the proper credentials before it can be used.
---

# MongoDB Atlas Provider

- Website: [https://www.mongodb.com/cloud/atlas](https://www.mongodb.com/cloud/atlas)


# Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.11+
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

# Developing the Provider
If you wish to work on the provider, you'll first need [Go](https://golang.org/doc/install) installed on your machine (please check the [requirements](#Requirements) before proceeding).

Note: This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](https://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e $HOME/development/terraform-providers/).

Clone repository to: `$HOME/development/terraform-providers/`

```
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:mongodb/terraform-mongodbatlas
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

$ make tools
To compile the provider, run `make build`. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```
$ make build
...
$ $GOPATH/bin/terraform-provider-mongodbatlas
...
```

# Using the Provider

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [`Terraform documentation on provider versioning`](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory, run terraform init to initialize it.

For either installation method, documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/).

# Testing the Provider

In order to test the provider, you can run `make test`. You need to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g.
`alias` and `version`), the following arguments are supported in the MongoDB
Atlas `provider` block:

* `public_key` - (Optional) This is the MongoDB Atlas API publick_key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PUBLIC_KEY`
  environment variable.

* `private_key` - (Optional) This is the MongoDB Atlas private_key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PRIVATE_KEY`
  environment variable.

```
$ make test
```

In order to run the full suite of Acceptance tests, run ``make testacc``.

Note: Acceptance tests create real resources, and often cost money to run. Please read Running an Acceptance Test in the contribution guidelines for more information on usage.

```
$ make testacc
```
For more information about how to get this programmatic API Keys see the following [link](https://docs.atlas.mongodb.com/configure-api-access/#manage-programmatic-access-to-an-organization).

Contributing
---------------------------

Terraform is the work of thousands of contributors. We appreciate your help!

To contribute, please read the contribution guidelines: [Contributing to Terraform - MongoDB Atlas Provider](/CONTRIBUTING.md)

Issues on GitHub are intended to be related to bugs or feature requests with provider codebase. See https://www.terraform.io/docs/extend/community/index.html for a list of community resources to ask questions about Terraform.

