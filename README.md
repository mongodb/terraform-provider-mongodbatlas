# MongoDB Atlas Provider
[![Automated Tests](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/automated-test.yml/badge.svg)](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/automated-test.yml)
[![Automated Acceptances Tests](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/automated-test-acceptances.yml/badge.svg)](https://github.com/mongodb/terraform-provider-mongodbatlas/actions/workflows/automated-test-acceptances.yml)

This is the repository for the Terraform MongoDB Atlas Provider, which allows one to use Terraform with MongoDB's Database as a Service offering, Atlas.
Learn more about Atlas at  [https://www.mongodb.com/cloud/atlas](https://www.mongodb.com/cloud/atlas)

For general information about Terraform, visit the [official website](https://www.terraform.io) and the [GitHub project page](https://github.com/hashicorp/terraform).

## Support, Bugs, Feature Requests

Support for the Terraform MongoDB Atlas Provider is provided under MongoDB Atlas support plans.   Please submit support questions within the Atlas UI.  Support questions submitted under the Issues section of this repo will be handled on a "best effort" basis.

Bugs should be filed under the Issues section of this repo.

Feature requests can be submitted at https://feedback.mongodb.com/forums/924145-atlas - just select the Terraform plugin as the category or vote for an already suggested feature.

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

## Using the Provider

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [`Terraform documentation on provider versioning`](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

Documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/).

## Beta Features
In order to use and/or test beta resources and datasources in this provider you'll need to set the environment variable `MONGODB_ATLAS_ENABLE_BETA` to true.


```bash
export MONGODB_ATLAS_ENABLE_BETA = true
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](https://golang.org/doc/install) installed on your machine (please check the [requirements](#Requirements) before proceeding).

Note: This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](https://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e $HOME/development/terraform-providers/).

Clone repository to: `$HOME/development/terraform-providers/`

```bash
mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
git clone git@github.com:mongodb/terraform-provider-mongodbatlas
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

```bash
make tools
```

To compile the provider, run `make build`. This will build the provider and put its binary in the ./bin directory.

```bash
make build
...
# ./bin/terraform-provider-mongodbatlas
...
```

### Using development provider in Terraform 0.14+
From terraform version 0.14, we can override provider use for development purposes.

Just create a `.trfc` file to hold the configuration to override terraform local configuration

```terraform
provider_installation {

  dev_overrides {
    "mongodb/mongodbatlas" = "[PATH THAT CONTAINS CUSTOM PLUGIN]"
  }

  direct {}
}
```

and set the env var `TF_CLI_CONFIG_FILE` like this:
`export TF_CLI_CONFIG_FILE=PATH/TO/dev.trfc`

For more explained information about "plugin override" check [Development Overrides for Provider Developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers)

### Testing the Provider

In order to test the provider, you can run `make test`. You can use [meta-arguments](https://www.terraform.io/docs/configuration/providers.html) such as `alias` and `version`. The following arguments are supported in the MongoDB Atlas `provider` block:

* `public_key` - (Optional) This is the MongoDB Atlas API public_key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PUBLIC_KEY`
  environment variable.

* `private_key` - (Optional) This is the MongoDB Atlas private_key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PRIVATE_KEY`
  environment variable.

~> **Notice:**  If you do not have a `public_key` and `private_key` you must create a programmatic API key to configure the provider (see [Creating Programmatic API key](#Programmatic-API-key)). If you already have one, you can continue with [Configuring environment variables](#Configuring-environment-variables)

#### Testing Atlas Provider Versions that are NOT hosted on Terraform Registry (i.e. pre-release versions)
To test development / pre-release versions of the Terraform Atlas Provider that are not hosted on the Terraform Registry, you will need to create a [Terraform Provider Network Mirror](https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol). 

The provider network mirror protocol is an optional protocol which you can implement to provide an alternative installation source for Terraform providers, regardless of their origin registries. Terraform uses network mirrors only if you activate them explicitly in the CLI configuration's `provider_installation` block. When enabled, a network mirror can serve providers belonging to any registry hostname, which can allow an organization to serve all of the Terraform providers they intend to use from an internal server, rather than from each provider's origin registry.

To do this you can: 
1.	Create a versions.tf file in a new directory with an existing live version from Terraform Registry. For example this can include:  
```
terraform {
  required_providers {
    mongodbatlas = {
      source = "mongodb/mongodbatlas"
    }
  }
  required_version = ">= 0.13"
}
```

2.	Use `terraform init` to create required file structures 

3.	`mkdir` a “tf_cache” sub-directory and `cd` into that directory 

4.	Create a .terraformrc file and insert below (modify accordingly to your own local path directory): 
```
provider_installation {
  filesystem_mirror {
    path    = "C:\Users\ZuhairAhmed\Desktop\Tenant_Upgrade\tf_cache"
    include = ["registry.terraform.io/hashicorp/*"]
  }
  direct {
    exclude = ["registry.terraform.io/hashicorp/*"]
  }
}
plugin_cache_dir = "C:\Users\ZuhairAhmed\Desktop\Tenant_Upgrade\tf_cache"
disable_checkpoint=true
```
5.	`cd` back up to original directory and `mv` the “.terraform/providers/registry.terraform.io” directory to “tf_cache” 

6.	Create required environment variables (modify accordingly to your own local path directory):
```
export TF_PLUGIN_CACHE_DIR=/mnt/c/Users/ZuhairAhmed/Desktop/Tenant_Upgrade/tf_cache
export TF_CLI_CONFIG_FILE=/mnt/c/Users/ZuhairAhmed/Desktop/Tenant_Upgrade/tf_cache/terraform.rc
```
7.  Delete the .terraform and .terraform.lock.hcl directories altogether. At this point you should only have the “tf_cache” directory and the “versions.tf” config file remaining. 

9. Next in the “tf_cache” directory replace existing terraform provider core files (Terraform Atlas Provider version binary, CHANGELOG.md, LICENSE, and README.md) with the version you seek to test locally. Make sure to keep folder structure the same. 

8. Lastly, in the terminal run `terraform init` again and this time terraform will pull provider version from tf_cache network mirror. You can confirm this by noting the “Terraform has been successfully initialized! Using mongodb/mongodbatlas Vx.x.x from the shared cache directory” message.  

#### Logs
To help with dubbing issues, you can turn on Logs with `export TF_LOG=TRACE`. Note: this is very noisy. 

To export logs to file, you can use `export TF_LOG_PATH=terraform.log`

### Running the acceptance test

#### Programmatic API key

It's necessary to generate and configure an API key for your organization for the acceptance test to succeed. To grant programmatic access to an organization or project using only the [API](https://docs.atlas.mongodb.com/api/) you need to know:

  1. The programmatic API key has two parts: a Public Key and a Private Key. To see more details on how to create a programmatic API key visit https://docs.atlas.mongodb.com/configure-api-access/#programmatic-api-keys.

  1. The programmatic API key must be granted roles sufficient for the acceptance test to succeed. The Organization Owner and Project Owner roles should be sufficient. You can see the available roles at https://docs.atlas.mongodb.com/reference/user-roles.

  1. You must [configure Atlas API Access](https://docs.atlas.mongodb.com/configure-api-access/) for your programmatic API key. You should allow API access for the IP address from which the acceptance test runs.

#### Configuring environment variables

You must also configure the following environment variables before running the test:

##### MongoDB Atlas env variables
```sh
export MONGODB_ATLAS_PROJECT_ID=<YOUR_PROJECT_ID>
export MONGODB_ATLAS_ORG_ID=<YOUR_ORG_ID>
export MONGODB_ATLAS_PUBLIC_KEY=<YOUR_PUBLIC_KEY>
export MONGODB_ATLAS_PRIVATE_KEY=<YOUR_PRIVATE_KEY>

# This env variable is optional and allow you to run terraform with a custom server
export MONGODB_ATLAS_BASE_URL=<CUSTOM_SERVER_URL>
```

- For `Authentication database user` resource configuration:
```sh
$ export MONGODB_ATLAS_DB_USERNAME=<YOUR_DATABASE_NAME>
```

- For `Project(s)` resource configuration:
```sh
$ export MONGODB_ATLAS_TEAMS_IDS=<YOUR_TEAMS_IDS>
$ export MONGODB_ATLAS_API_KEYS_IDS=<API_KEYS_IDS>
```
~> **Notice:** It should be at least one team id up to 3 teams ids depending of acceptance testing using separator comma like this `teamId1,teamdId2,teamId3`.

- For skip acceptances testing that requires additional credentials such as AWS, AZURE and GCP:
```sh
export SKIP_TEST_EXTERNAL_CREDENTIALS=TRUE
```

- For `Federated Settings` resource configuration:
```sh
$ export MONGODB_ATLAS_FEDERATION_SETTINGS_ID=<YOUR_FEDERATION_SETTINGS_ID>
$ export ONGODB_ATLAS_FEDERATED_ORG_ID=<YOUR_FEDERATED_ORG_ID>
$ export MONGODB_ATLAS_FEDERATED_PROJECT_ID=<YOUR_VPC_FEDERATED_PROJECT_ID>
$ export MONGODB_ATLAS_FEDERATED_GROUP_ID=<YOUR_FEDERATED_GROUP_ID>
$ export MONGODB_ATLAS_FEDERATED_ROLE_MAPPING_ID=<YOUR_FEDERATED_ROLE_MAPPING_ID>
$ export MONGODB_ATLAS_FEDERATED_OKTA_IDP_ID=<YOUR_FEDERATED_OKTA_IDP_ID>
$ export MONGODB_ATLAS_FEDERATED_SSO_URL=<YOUR_FEDERATED_SSO_URL>
$ export MONGODB_ATLAS_FEDERATED_ISSUER_URI=<YOUR_FEDERATED_ISSUER_URI>
```
~> **Notice:** For more information about the Federation configuration resource, see: https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/

##### AWS env variables

- For `Network Peering` resource configuration:
```sh
$ export AWS_ACCOUNT_ID=<YOUR_ACCOUNT_ID>
$ export AWS_VPC_ID=<YOUR_VPC_ID>
$ export AWS_VPC_CIDR_BLOCK=<YOUR_VPC_CIDR_BLOCK>
$ export AWS_REGION=<YOUR_REGION>
$ export AWS_SUBNET_ID=<YOUR_SUBNET_ID>
$ export AWS_SECURITY_GROUP_ID=<YOUR_SECURITY_GROUP_ID>
```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/

- For `Encryption at Rest` resource configuration:
```sh
$ export AWS_ACCESS_KEY_ID=<YOUR_ACCESS_KEY_ID>
$ export AWS_SECRET_ACCESS_KEY=<YOUR_SECRET_ACCESS_KEY>
$ export AWS_CUSTOMER_MASTER_KEY_ID=<YOUR_CUSTOMER_MASTER_KEY_ID>
$ export AWS_REGION=<YOUR_REGION>

$ export AWS_ACCESS_KEY_ID_UPDATED=<YOUR_ACCESS_KEY_ID_UPDATED>
$ export AWS_SECRET_ACCESS_KEY_UPDATED=<YOUR_SECRET_ACCESS_KEY_UPDATED>
$ export AWS_CUSTOMER_MASTER_KEY_ID_UPDATED=<YOUR_CUSTOMER_MASTER_KEY_ID_UPDATED>
$ export AWS_REGION_UPDATED=<YOUR_REGION_UPDATED>
```
~> **Notice:** For more information about the Encryption at Rest resource, see: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/

- For `Private Endpoint Link` resource configuration:
```sh
$ export AWS_ACCESS_KEY_ID=<YOUR_ACCESS_KEY_ID>
$ export AWS_SECRET_ACCESS_KEY=<YOUR_SECRET_ACCESS_KEY>
$ export AWS_CUSTOMER_MASTER_KEY_ID=<YOUR_CUSTOMER_MASTER_KEY_ID>
$ export AWS_REGION=<YOUR_REGION>
$ export AWS_VPC_ID=<YOUR_VPC_ID>
$ export AWS_SUBNET_ID=<YOUR_SUBNET_ID>
$ export AWS_SECURITY_GROUP_ID=<YOUR_SECURITY_GROUP_ID>
```
~> **Notice:** For more information about the PrivateLink (for AWS only), see: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/https://docs.atlas.mongodb.com/reference/api/private-endpoint/

##### AZURE env variables

- For `Network Peering` resource configuration:
```sh
$ export AZURE_DIRECTORY_ID=<YOUR_DIRECTORY_ID>
$ export AZURE_SUBSCRIPTION_ID=<YOUR_SUBSCRIPTION_ID>
$ export AZURE_RESOURCE_GROUP_NAME=<YOUR_RESOURCE_GROUP_NAME>
$ export AZURE_VNET_NAME=<YOUR_VNET_NAME>
```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/


- For Encryption at Rest resource configuration:
```sh
export AZURE_CLIENT_ID=<YOUR_CLIENT_ID>
export AZURE_SUBSCRIPTION_ID=<YOUR_SUBSCRIPTION_ID>
export AZURE_RESOURCE_GROUP_NAME=<YOUR_RESOURCE_GROUP_NAME>
export AZURE_SECRET=<YOUR_SECRET>
export AZURE_KEY_VAULT_NAME=<YOUR_KEY_VAULT_NAME>
export AZURE_KEY_IDENTIFIER=<YOUR_KEY_IDENTIFIER>
export AZURE_TENANT_ID=<YOUR_TENANT_ID>
export AZURE_DIRECTORY_ID=<YOUR_DIRECTORY_ID>

export AZURE_CLIENT_ID_UPDATED=<YOUR_CLIENT_ID_UPDATED>
export AZURE_RESOURCE_GROUP_NAME_UPDATED=<YOUR_RESOURCE_GROUP_NAME_UPDATED>
export AZURE_SECRET_UPDATED=<YOUR_SECRET_UPDATED>
export AZURE_KEY_VAULT_NAME_UPDATED=<YOUR_KEY_VAULT_NAME_UPDATED>
export AZURE_KEY_IDENTIFIER_UPDATED=<YOUR_KEY_IDENTIFIER_UPDATED>
```
~> **Notice:** For more information about the Encryption at Rest resource, see: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/

##### GCP env variables
- For `Network Peering` resource configuration:
```sh
$export GCP_PROJECT_ID=<YOUR_PROJECT_ID>
```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/


- For Encryption at Rest resource configuration:
```sh
$ export GCP_SERVICE_ACCOUNT_KEY=<YOUR_GCP_SERVICE_ACCOUNT_KEY>
$ export GCP_KEY_VERSION_RESOURCE_ID=<YOUR_GCP_KEY_VERSION_RESOURCE_ID>
```
~> **Notice:** For more information about the Encryption at Rest resource, see: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/


In order to run the full suite of Acceptance tests, run ``make testacc``.

~> **Notice:** Acceptance tests create real resources, and often cost money to run. Please note in any PRs made if you are unable to pay to run acceptance tests for your contribution. We will accept "best effort" implementations of acceptance tests in this case and run them for you on our side. This may delay the contribution but we do not want your contribution blocked by funding.

```
$ make testacc
```
### Running the integration tests

The integration tests helps the validation for resources interacting with third party providers (AWS, Azure or Google Cloud) using terratest [environment setup details](integrationtesting/README.md)

```
  cd integrationtesting
  go test -tags=integration
```

### Autoclose stale issues and PRs

- After 30 days of no activity (no comments or commits on an issue/PR) we automatically tag it as "stale" and add a message: ```This issue/PR has gone 30 days without any activity and meets the project's definition of "stale". This will be auto-closed if there is no new activity over the next 30 days. If the issue is still relevant and active, you can simply comment with a "bump" to keep it open, or add the label "not_stale". Thanks for keeping our repository healthy!```
- After 30 more days of no activity we automatically close the issue/PR.

## Thanks

We'd like to thank [Akshay Karle](https://github.com/akshaykarle) for writing the first version of a Terraform Provider for MongoDB Atlas and paving the way for the creation of this one.
