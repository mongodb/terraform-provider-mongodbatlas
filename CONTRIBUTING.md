# Contributing

Thanks for your interest in contributing to MongoDB Atlas Terraform Provider, this document describes some guidelines necessary to participate in the community.

## Table of Contents

- [Development Setup](#development-setup)
  - [Prerequisite Tools](#prerequisite-tools)
  - [Environment](#prerequisite-tools)
  - [Open a Pull Request](#open-a-pull-request)
  - [Testing the Provider](#testing-the-provider)
  - [Running Acceptance Tests](#running-acceptance-tests)
- [Code and Test Best Practices](#code-and-test-best-practices)
  - [Creating New Resource and Data Sources](#creating-new-resources-and-data-sources)
    - [Scaffolding Initial Code and File Structure](#scaffolding-initial-code-and-file-structure)
    - [Scaffolding Schema and Model Definitions](#scaffolding-schema-and-model-definitions)
- [Documentation Best Practices](#documentation-best-practices)
  - [Creating Resource and Data source Documentation](#creating-resource-and-data-source-documentation)
- [Discovering New API features](#discovering-new-api-features)


## Development Setup
### Prerequisite Tools

- [Git](https://git-scm.com/)
- [Go (at least Go 1.22)](https://golang.org/dl/)

### Environment

- Fork the repository.
- Clone your forked repository locally.
- We use Go Modules to manage dependencies, so you can develop outside your `$GOPATH`.
- We use [golangci-lint](https://github.com/golangci/golangci-lint) to lint our code, you can install it locally via `make setup`.
### Building
- Enter the provider directory
- Run `make tools` to install the needed tools for the provider
- Run `link-git-hooks` to install githooks 
- Run `make build` to build the binary in the `./bin` directory: 
- Use the local provider binary in the `./bin` folder:
  - Create the following `dev.trfc` file inside your directory 
  ```terraform
    provider_installation {

    dev_overrides {
      "mongodb/mongodbatlas" = "/Users/yourUser/terraform-provider-mongodbatlas/bin" # path to the provider binary
    }

    direct {} 
  }
  ```
  - Define the env var `TF_CLI_CONFIG_FILE` in your console session
  ```bash
  export TF_CLI_CONFIG_FILE=PATH/TO/dev.trfc
  ```
- Run `terraform init` to inizialize terraform
- Run `terraform apply` to use terraform with the local binary

For more explained information about plugin override check [Development Overrides for Provider Developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers)

### Open a Pull Request
- Sign the [contributor's agreement](http://www.mongodb.com/contributor). This will allow us to review and accept contributions.
- Implement your feature, improvement or bug fix, ensuring it adheres to the [Terraform Plugin Best Practices](https://www.terraform.io/docs/extend/best-practices/index.html).
- Make sure that the PR title follows [*Conventional Commits*](https://www.conventionalcommits.org/).
- Add comments around your new code that explain what's happening.
- Commit and push your changes to your branch then submit a pull request against the `master` branch.
- A repo maintainer will review the your pull request, and may either request additional changes or merge the pull request.

#### PR Title Format

Use [*Conventional Commits*](https://www.conventionalcommits.org/) to name pull requests, starting with the type of change followed by a description of the change. Use a third person point of view, [active voice](https://www.mongodb.com/docs/meta/style-guide/writing/use-active-voice/#std-label-use-active-voice), and start each description with an uppercase character:

- `fix: Description of the PR`: a commit of the type fix patches a bug in your codebase (this correlates with PATCH in Semantic Versioning).
- `chore: Description of the PR`: the commit includes a technical or preventative maintenance task that is necessary for managing the product or the repository, but it is not tied to any specific feature or user story (this correlates with PATCH in Semantic Versioning).
- `doc: Description of the PR`: The commit adds, updates, or revises documentation that is stored in the repository (this correlates with PATCH in Semantic Versioning).
- `test: Description of the PR`: The commit enhances, adds to, revised, or otherwise changes the suite of automated tests for the product (this correlates with PATCH in Semantic Versioning).
- `security: Description of the PR`: The commit improves the security of the product or resolves a security issue that has been reported (this correlates with PATCH in Semantic Versioning).
- `refactor: Description of the PR`: The commit refactors existing code in the product, but does not alter or change existing behavior in the product (this correlates with Minor in Semantic Versioning).
- `perf: Description of the PR`: The commit improves the performance of algorithms or general execution time of the product, but does not fundamentally change an existing feature (this correlates with Minor in Semantic Versioning).
- `ci: Description of the PR`: The commit makes changes to continuous integration or continuous delivery scripts or configuration files (this correlates with Minor in Semantic Versioning).
- `revert: Description of the PR`: The commit reverts one or more commits that were previously included in the product, but were accidentally merged or serious issues were discovered that required their removal from the main branch (this correlates with Minor in Semantic Versioning).
- `style: Description of the PR`: The commit updates or reformats the style of the source code, but does not otherwise change the product implementation (this correlates with Minor in Semantic Versioning).
- `feat: Description of the PR`: a commit of the type feat introduces a new feature to the codebase (this correlates with MINOR in Semantic Versioning).
- `deprecate: Description of the PR`: The commit deprecates existing functionality, but does not remove it from the product (this correlates with MINOR in Semantic Versioning).
- `BREAKING CHANGE`: a commit that has a footer BREAKING CHANGE:, or appends a ! after the type/scope, introduces a breaking API change (correlating with MAJOR in Semantic Versioning). A BREAKING CHANGE can be part of commits of any type.
Examples:
  - `fix!: Description of the ticket`
  - If the PR has `BREAKING CHANGE`: in its description is a breaking change
- `remove!: Description of the PR`: The commit removes a feature from the product. Typically features are deprecated first for a period of time before being removed. Removing a feature is a breaking change (correlating with MAJOR in Semantic Versioning).

Example PR title:
  ```bash
  chore: Upgrades `privatelink_endpoint_service_data_federation_online_archive` resource to auto-generated SDK
  ```

- The example PR title starts with a task type, "chore:"
- The description begins with an uppercase character, "U" in "Upgrades"
- The description uses active voice with the verb "Upgrades," where the subject is 
  implicitly the PR itself.

### Testing the Provider

In order to test the provider, you can run `make test`. You can use [meta-arguments](https://www.terraform.io/docs/configuration/providers.html) such as `alias` and `version`. The following arguments are supported in the MongoDB Atlas `provider` block:

* `public_key` - (Optional) This is the MongoDB Atlas API public_key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PUBLIC_KEY`
  environment variable.

* `private_key` - (Optional) This is the MongoDB Atlas private_key. It must be
  provided, but it can also be sourced from the `MONGODB_ATLAS_PRIVATE_KEY`
  environment variable.

~> **Notice:**  If you do not have a `public_key` and `private_key` you must create a programmatic API key to configure the provider (see [Creating Programmatic API key](#Programmatic-API-key)). If you already have one, you can continue with [Configuring environment variables](#Configuring-environment-variables)

### Running Acceptance Tests

#### Programmatic API key

It's necessary to generate and configure an API key for your organization for the acceptance test to succeed. To grant programmatic access to an organization or project using only the [API](https://docs.atlas.mongodb.com/api/) you need to know:

  - The programmatic API key has two parts: a Public Key and a Private Key. To see more details on how to create a programmatic API key visit https://docs.atlas.mongodb.com/configure-api-access/#programmatic-api-keys.

  - The programmatic API key must be granted roles sufficient for the acceptance test to succeed. The Organization Owner and Project Owner roles should be sufficient. You can see the available roles at https://docs.atlas.mongodb.com/reference/user-roles.

  - You must [configure Atlas API Access](https://docs.atlas.mongodb.com/configure-api-access/) for your programmatic API key. You should allow API access for the IP address from which the acceptance test runs.

#### Configuring environment variables

You must also configure the following environment variables before running the test:

##### MongoDB Atlas env variables
- Required env variables:
  ```bash
  export MONGODB_ATLAS_PROJECT_ID=<YOUR_PROJECT_ID>
  export MONGODB_ATLAS_ORG_ID=<YOUR_ORG_ID>
  export MONGODB_ATLAS_PUBLIC_KEY=<YOUR_PUBLIC_KEY>
  export MONGODB_ATLAS_PRIVATE_KEY=<YOUR_PRIVATE_KEY>

  # This env variable is optional and allow you to run terraform with a custom server
  export MONGODB_ATLAS_BASE_URL=<CUSTOM_SERVER_URL>
  ```

- For `Authentication database user` resource configuration:
  ```bash
  $ export MONGODB_ATLAS_DB_USERNAME=<YOUR_DATABASE_NAME>
  ```

- For `Project(s)` resource configuration:
  ```bash
  $ export MONGODB_ATLAS_TEAMS_IDS=<YOUR_TEAMS_IDS>
  ```
~> **Notice:** It should be at least one team id up to 3 teams ids depending of acceptance testing using separator comma like this `teamId1,teamdId2,teamId3`.

- For `Federated Settings` resource configuration:
  ```bash
  $ export MONGODB_ATLAS_FEDERATION_SETTINGS_ID=<YOUR_FEDERATION_SETTINGS_ID>
  $ export ONGODB_ATLAS_FEDERATED_ORG_ID=<YOUR_FEDERATED_ORG_ID>
  $ export MONGODB_ATLAS_FEDERATED_GROUP_ID=<YOUR_FEDERATED_GROUP_ID>
  $ export MONGODB_ATLAS_FEDERATED_ROLE_MAPPING_ID=<YOUR_FEDERATED_ROLE_MAPPING_ID>
  $ export MONGODB_ATLAS_FEDERATED_OKTA_IDP_ID=<YOUR_FEDERATED_OKTA_IDP_ID>
  $ export MONGODB_ATLAS_FEDERATED_SSO_URL=<YOUR_FEDERATED_SSO_URL>
  $ export MONGODB_ATLAS_FEDERATED_ISSUER_URI=<YOUR_FEDERATED_ISSUER_URI>
  ```
~> **Notice:** For more information about the Federation configuration resource, see: https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/

##### AWS env variables

- For `Network Peering` resource configuration:
  ```bash
  $ export AWS_ACCOUNT_ID=<YOUR_ACCOUNT_ID>
  $ export AWS_VPC_ID=<YOUR_VPC_ID>
  $ export AWS_VPC_CIDR_BLOCK=<YOUR_VPC_CIDR_BLOCK>
  $ export AWS_REGION=<YOUR_REGION>
  $ export AWS_SUBNET_ID=<YOUR_SUBNET_ID>
  $ export AWS_SECURITY_GROUP_ID=<YOUR_SECURITY_GROUP_ID>
  ```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/

- For `Encryption at Rest` resource configuration:
  ```bash
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
  ```bash
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
  ```bash
  $ export AZURE_DIRECTORY_ID=<YOUR_DIRECTORY_ID>
  $ export AZURE_SUBSCRIPTION_ID=<YOUR_SUBSCRIPTION_ID>
  $ export AZURE_RESOURCE_GROUP_NAME=<YOUR_RESOURCE_GROUP_NAME>
  $ export AZURE_VNET_NAME=<YOUR_VNET_NAME>
  ```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/


- For Encryption at Rest resource configuration:
  ```bash
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
  ```bash
  $export GCP_PROJECT_ID=<YOUR_PROJECT_ID>
  ```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/


- For Encryption at Rest resource configuration:
  ```bash
  $ export GCP_SERVICE_ACCOUNT_KEY=<YOUR_GCP_SERVICE_ACCOUNT_KEY>
  $ export GCP_KEY_VERSION_RESOURCE_ID=<YOUR_GCP_KEY_VERSION_RESOURCE_ID>
  ```
~> **Notice:** For more information about the Encryption at Rest resource, see: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/


#### Run Acceptance tests
~> **Notice:** Acceptance tests create real resources, and often cost money to run. Please note in any PRs made if you are unable to pay to run acceptance tests for your contribution. We will accept "best effort" implementations of acceptance tests in this case and run them for you on our side. This may delay the contribution but we do not want your contribution blocked by funding.
- Run `make testacc`



### Testing Atlas Provider Versions that are NOT hosted on Terraform Registry (i.e. pre-release versions)
To test development / pre-release versions of the Terraform Atlas Provider that are not hosted on the Terraform Registry, you will need to create a [Terraform Provider Network Mirror](https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol). 

The provider network mirror protocol is an optional protocol which you can implement to provide an alternative installation source for Terraform providers, regardless of their origin registries. Terraform uses network mirrors only if you activate them explicitly in the CLI configuration's `provider_installation` block. When enabled, a network mirror can serve providers belonging to any registry hostname, which can allow an organization to serve all of the Terraform providers they intend to use from an internal server, rather than from each provider's origin registry.

To do this you can: 
-	Create a versions.tf file in a new directory with an existing live version from Terraform Registry. For example this can include:  
    ```terraform
    terraform {
      required_providers {
        mongodbatlas = {
          source = "mongodb/mongodbatlas"
        }
      }
      required_version = ">= 0.13"
    }
    ```

-	Use `terraform init` to create required file structures 

-	`mkdir` a `tf_cache` sub-directory and `cd` into that directory 

-	Create a .terraformrc file and insert below (modify accordingly to your own local path directory): 
    ```terraform
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
-	`cd` back up to original directory and `mv` the `.terraform/providers/registry.terraform.io` directory to `tf_cache` 

-	Create required environment variables (modify accordingly to your own local path directory):
    ```bash
    export TF_PLUGIN_CACHE_DIR=/mnt/c/Users/ZuhairAhmed/Desktop/Tenant_Upgrade/tf_cache
    export TF_CLI_CONFIG_FILE=/mnt/c/Users/ZuhairAhmed/Desktop/Tenant_Upgrade/tf_cache/terraform.rc
    ```
-  Delete the `.terraform` and `.terraform.lock.hcl` directories altogether. At this point you should only have the `tf_cache` directory and the `versions.tf` config file remaining. 

- Next in the `tf_cache` directory replace existing terraform provider core files (Terraform Atlas Provider version binary, `CHANGELOG.md`, `LICENSE`, and `README.md`) with the version you seek to test locally. Make sure to keep folder structure the same. 

- Lastly, in the terminal run `terraform init` again and this time terraform will pull provider version from `tf_cache` network mirror. You can confirm this by noting the `Terraform has been successfully initialized! Using mongodb/mongodbatlas Vx.x.x from the shared cache directory` message.  

## Code and Test Best Practices

- Each resource (and associated data sources) is in a package in `internal/service`.
- There can be multiple helper files and they can also be used from other resources, e.g. `common_advanced_cluster.go` defines functions that are also used from other resources using `advancedcluster.FunctionName`.
- Unit and Acceptances tests are in the same `_test.go` file. They are not in the same package as the code tests, e.g. `advancedcluster` tests are in `advancedcluster_test` package so coupling is minimized.
- Migration tests are in `_migration_test.go` files.
- Helper methods must have their own tests, e.g. `common_advanced_cluster_test.go` has tests for `common_advanced_cluster.go`.
- `internal/testutils/acc` contains helper test methods for Acceptance and Migration tests.
- Tests that need the provider binary like End-to-End tests don’t belong to the source code packages and go in `test/e2e`.
- [Testify Mock](https://pkg.go.dev/github.com/stretchr/testify/mock) and [Mockery](https://github.com/vektra/mockery) are used for test doubles in unit tests. Mocked interfaces are generated in folder `internal/testutil/mocksvc`.


### Creating New Resource and Data Sources

A set of commands have been defined with the intention of speeding up development process, while also preserving common conventions throughout our codebase.

#### Scaffolding Initial Code and File Structure

This command can be used the following way:
```bash
make scaffold resource_name=streamInstance type=resource
```
- **resource_name**: The name of the resource, which must be defined in camel case.
- **type**: Describes the type of resource being created. There are 3 different types: `resource`, `data-source`, `plural-data-source`.

This will generate resource/data source files and accompanying test files needed for starting the development, and will contain multiple comments with `TODO:` statements which give guidance for the development.

As a follow up step, use [Scaffolding Schema and Model Definitions](#scaffolding-schema-and-model-definitions) to autogenerate the schema via the Open API specification. This will require making adjustments to the generated `./internal/service/<resource_name>/tfplugingen/generator_config.yml` file.

#### Scaffolding Schema and Model Definitions

Complementary to the `scaffold` command, there is a command which generates the initial Terraform schema definition and associated Go types for a resource or data source. This processes leverages [Code Generation Tools](https://developer.hashicorp.com/terraform/plugin/code-generation) developed by HashiCorp, which in turn make use of the [Atlas Admin API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/) OpenAPI Specification.

##### Running the command

Both `tfplugingen-openapi` and `tfplugingen-framework` must be installed. This can be done by running `make tools`.

The command takes a single argument which specifies the resource or data source where the code generation is run, defined in camel case, e.g.:
```bash
make scaffold-schemas resource_name=streamInstance
```

As a pre-requiste, the relevant resource/data source directory must define a configuration file in the path `./internal/service/<resource_name>/tfplugingen/generator_config.yml`. The content of this file will define which resource and/or data source schemas will be generated by providing the API endpoints they are mapped to. See the [Generator Config](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator#generator-config) documentation for more information on configuration options. An example defined in our repository can be found in [searchdeployment/tfplugingen/generator_config.yml](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/internal/service/searchdeployment/tfplugingen/generator_config.yml).

As a result of the execution, the schema definitions and associated model types will be defined in separate files depending on the resources and data sources that were configured in the generator_config.yml file:
- `data_source_<resource_name>_schema.go`
- `resource_<resource_name>_schema.go`

Note: if the resulting file paths already exist the content will be stored in files with a `_gen.go` postfix, and in this case any content will be overwritten. This can be useful for comparing the latest autogenerated schema against the existing implementation.

Note: you can override the Open API description of a field with a custom description via the [overrides](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator#overriding-attribute-descriptions) param. See this [example](internal/service/searchdeployment/tfplugingen/generator_config.yml).

##### Considerations over generated schema and types

- Generated Go type should include a TF prefix to follow the convention in our codebase, this will not be present in generated code.
- Some attribute names may need to be adjusted if there is a difference in how they are named in Terraform vs the API. An examples of this is `group_id` → `project_id`.
- Inferred characteristics of an attribute (computed, optional, required) may not always be an accurate representation and should be revised. Details of inference logic can be found in [OAS Types to Provider Attributes](https://github.com/hashicorp/terraform-plugin-codegen-openapi/blob/main/DESIGN.md#oas-types-to-provider-attributes).
- Missing [sensitive](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes/string#sensitive) field in attributes.
- Missing plan modifiers such as `RequiresReplace()` in attributes.
- Terraform specific attributes such as [timeouts](https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts#specifying-timeouts-in-configuration) need to be included manually.
- If nested attributes are defined a set of helper functions are generated for using the model. The usage of the generated functions can be considered optional as the current documentation is not very clear on the usage (More details in [terraform-plugin-codegen-framework/issues/80](https://github.com/hashicorp/terraform-plugin-codegen-framework/issues/80)).


## Documentation Best Practices

- In our documentation, when a resource field allows a maximum of only one item, we do not format that field as an array. Instead, we create a subsection specifically for this field. Within this new subsection, we enumerate all the attributes of the field. Let's illustrate this with an example: [cloud_backup_schedule.html.markdown](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/website/docs/r/cloud_backup_schedule.html.markdown?plain=1#L207)
- You can check how the documentation is rendered on the Terraform Registry via [doc-preview](https://registry.terraform.io/tools/doc-preview).

### Creating Resource and Data source Documentation
We autogenerate the documentation of our provider resources and data sources via [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs).

#### How to generate the documentation for a resource
- Make sure that the resource and data source schemas have defined the fields `MarkdownDescription` and `Description`.
  - We recommend to use [Scaffolding Schema and Model Definitions](#scaffolding-schema-and-model-definitions) to autogenerate the schema via the Open API specification.
- Add the resource/data source templates to the [templates](templates) folder. See [README.md](templates/README.md) for more info.
- Run the Makefile command `generate-doc`
```bash
export resource_name=search_deployment && make generate-doc
```

## Discovering New API features

Most of the new features of the provider are using [atlas-sdk](https://github.com/mongodb/atlas-sdk-go)
SDK is updated automatically, tracking all new Atlas features.


### Updating Atlas SDK 

To update Atlas SDK run:

```bash
make update-atlas-sdk
```

> NOTE: The update mechanism is only needed for major releases. Any other releases will be supported by dependabot.

> NOTE: Command can make import changes to +500 files. Please make sure that you perform update on main branch without any uncommited changes.

### SDK Major Release Update Procedure

1. If the SDK update doesn’t cause any compilation issues create a new SDK update PR
   1. Review [API Changelog](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/changelog) for any deprecated fields and breaking changes.
2. For SDK updates introducing compilation issues without graceful workaround
   1. Use the previous major version of the SDK (including the old client) for the affected resource
   1. Create an issue to identify the root cause and mitigation paths based on changelog information  
   2. If applicable: Make required notice/update to the end users based on the plan.
