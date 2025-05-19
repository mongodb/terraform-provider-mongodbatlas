# MongoDB Atlas Provider
[![Code Health](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/code-health.yml/badge.svg)](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/code-health.yml)
[![Acceptance Tests](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/acceptance-tests.yml/badge.svg)](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/acceptance-tests.yml)


This is the repository for the Terraform MongoDB Atlas Provider, which allows one to use Terraform with MongoDB's Database as a Service offering, Atlas.
Learn more about Atlas at  [https://www.mongodb.com/cloud/atlas](https://www.mongodb.com/cloud/atlas)

For general information about Terraform, visit the [official website](https://www.terraform.io) and the [GitHub project page](https://github.com/hashicorp/terraform).

## Support, Bugs, Feature Requests

Support for the Terraform MongoDB Atlas Provider is provided under MongoDB Atlas support plans.   Please submit support questions within the Atlas UI.  Support questions submitted under the Issues section of this repo will be handled on a "best effort" basis.

Bugs should be filed under the Issues section of this repo.

Feature requests can be submitted at https://feedback.mongodb.com/forums/924145-atlas - just select the Terraform plugin as the category or vote for an already suggested feature.

## Requirements  
- [HashiCorp Terraform Version Compatibility Matrix](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs#hashicorp-terraform-versionhttpswwwterraformiodownloadshtml-compatibility-matrix) 

- [Go Version](https://golang.org/doc/install) 1.23 (to build the provider plugin)

## Using the Provider

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [`Terraform documentation on provider versioning`](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

Documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/).

## Logs
To help with issues, you can turn on Logs with `export TF_LOG=TRACE`. Note: this is very noisy. 

To export logs to file, you can use `export TF_LOG_PATH=terraform.log`


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


## Troubleshooting
See [Troubleshooting](./docs/troubleshooting.md).

## Developing the Provider
See our [contributing guides](./contributing/README.md).

## Issues

### Autoclose stale issues and PRs
- After 5 days of no activity (no comments or commits on an issue/PR) we automatically tag it as "stale" and add a message: ```This issue/PR has gone 5 days without any activity and meets the project's definition of "stale". This will be auto-closed if there is no new activity over the next 5 days. If the issue is still relevant and active, you can simply comment with a "bump" to keep it open, or add the label "not_stale". Thanks for keeping our repository healthy!```
- After 5 more days of no activity we automatically close the issue/PR.

### One-click reproducible issues principle
Our support will prioritise issues that contain all the required information that follows the following principles:

* We should be able to make no changes to your provided script and **successfully run a local execution reproducing the issue**.
  * This means that you should kindly **provide all the required instructions**. This includes but not limited to:
    * Terraform Atlas provider version used to reproduce the issue
    * Terraform version used to reproduce the issue
  * Configurations that **cannot be properly executed will be de-prioritised** in favour of the ones that succeed.
* Before opening an issue, you have to try to specifically isolate it to **Terraform MongoDB Atlas** provider by **removing as many dependencies** as possible. If the issue only happens with other dependencies, then:
  * If other terraform providers are required, please make sure you also include those. _Same "one-click reproducible issue" principle applies_.
  * If external components are required to replicate it, please make sure you also provides instructions on those parts.
* Please confirm if the platform being used is Terraform OSS, Terraform Cloud, or Terraform Enterprise deployment


## Thanks

We'd like to thank [Akshay Karle](https://github.com/akshaykarle) for writing the first version of a Terraform Provider for MongoDB Atlas and paving the way for the creation of this one.
