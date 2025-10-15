# Development Setup

## Table of Contents
- [Prerequisite Tools](#prerequisite-tools)
- [Environment](#prerequisite-tools)
- [Open a Pull Request](#open-a-pull-request)
- [Testing the Provider](#testing-the-provider)
- [Running Acceptance Tests](#running-acceptance-tests)

### Prerequisite Tools

- [Git](https://git-scm.com/)
- [Go (at least Go 1.23)](https://golang.org/dl/)

### Environment

- Fork the repository.
- Clone your forked repository locally.
- We use Go Modules to manage dependencies, so you can develop outside your `$GOPATH`.
- We use [golangci-lint](https://github.com/golangci/golangci-lint) to lint our code, you can install it locally via `make tools`.
### Building
- Enter the provider directory
- Run `make tools` to install the needed tools for the provider
- Run `make link-git-hooks` to install githooks 
- Run `make build` to build the binary in the `./bin` directory: 
- Use the local provider binary in the `./bin` folder:
  - Create the following `dev.trfc` file inside your directory 
  ```terraform
  provider_installation {

    dev_overrides {
      "mongodb/mongodbatlas" = "/Users/<YourUser>/terraform-provider-mongodbatlas/bin" # path to the provider binary
    }

    direct {} 
  }
  ```
  - Define the env var `TF_CLI_CONFIG_FILE` in your console session
  ```bash
  export TF_CLI_CONFIG_FILE=PATH/TO/dev.trfc
  ```
- Run `terraform init` to initialize terraform
- Run `terraform apply` to use terraform with the local binary

For more explained information about plugin override check [Development Overrides for Provider Developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers)

### Open a Pull Request
- Sign the [contributor's agreement](http://www.mongodb.com/contributor). This will allow us to review and accept contributions.
- Implement your feature, improvement or bug fix, ensuring it adheres to the [Terraform Plugin Best Practices](https://www.terraform.io/docs/extend/best-practices/index.html).
- Make sure that the PR title follows [*Conventional Commits*](https://www.conventionalcommits.org/).
- Add comments around your new code that explain what's happening.
- Commit and push your changes to your branch then submit a pull request against the `master` branch.
  - We recommend squashing commits into one. This enables us to correctly cherry-pick your changes for merge, see [Merging a Pull Request](#merging-a-pull-request) below for details.
- A repo maintainer will review your pull request. **Note**: If you have an active [MongoDB Atlas Support](https://www.mongodb.com/services/support/atlas-support-plans) contract, we recommend also creating a support ticket for any questions related to this process.

### Merging a Pull Request
Due to security reasons, there are restrictions on how external contributions can be handled, especially concerning the use of repository secrets and running tests from forks.
Our GitHub setup does not allow sharing of variables and secrets across forks made by external contributors. This is a necessary security measure to prevent unauthorized access and potential vulnerabilities.

As external contributors cannot merge changes from their forks directly nor create branches from the master, the contribution process involves co-authoring. Here's what will happen:

- We will cherry-pick your commit into a new branch in our main repository.
- We will then create a new pull request from this branch, adding you as a co-author to ensure your contribution is not lost.
- Our team will carry out any necessary testing and review on the new pull request. You are welcome to participate in this phase by commenting and making further suggestions.
- Once approved, the pull request will be merged into the main branch.
- **Exception**: 
For minor changes that do not affect the core functionality of the project and do not require running tests, exceptions to the above process may be made at the discretion of the project maintainers.

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
- `a ! in the description of the PR`: The commit introduces a breaking change (correlating with MAJOR in Semantic Versioning). A breaking change can be part of commits of any type.
Examples:
  - `fix!: Description of the ticket`
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
  export MONGODB_ATLAS_PUBLIC_KEY=<ATLAS_PUBLIC_KEY>
  export MONGODB_ATLAS_PRIVATE_KEY=<ATLAS_PRIVATE_KEY>
  export MONGODB_ATLAS_ORG_ID=<ATLAS_ORG_ID>
  export MONGODB_ATLAS_PROJECT_ID=<ATLAS_PROJECT_ID>

  # This env variable is optional and allow you to run terraform with a custom server
  export MONGODB_ATLAS_BASE_URL=<CUSTOM_SERVER_URL>
  ```

- For `Authentication database user` resource configuration:
  ```bash
  export MONGODB_ATLAS_DB_USERNAME=<ATLAS_DATABASE_NAME>
  ```

- For `Project(s)` resource configuration:
  ```bash
  export MONGODB_ATLAS_TEAMS_IDS=<ATLAS_TEAMS_IDS>
  ```
~> **Notice:** It should be at least one team id up to 3 teams ids depending of acceptance testing using separator comma like this `teamId1,teamdId2,teamId3`.

- For `Federated Settings` resource configuration:
  ```bash
  export MONGODB_ATLAS_FEDERATION_SETTINGS_ID=<ATLAS_FEDERATION_SETTINGS_ID>
  export MONGODB_ATLAS_FEDERATED_ORG_ID=<ATLAS_FEDERATED_ORG_ID>
  export MONGODB_ATLAS_FEDERATED_GROUP_ID=<ATLAS_FEDERATED_GROUP_ID>
  export MONGODB_ATLAS_FEDERATED_ROLE_MAPPING_ID=<ATLAS_FEDERATED_ROLE_MAPPING_ID>
  export MONGODB_ATLAS_FEDERATED_OKTA_IDP_ID=<ATLAS_FEDERATED_OKTA_IDP_ID>
  export MONGODB_ATLAS_FEDERATED_SSO_URL=<ATLAS_FEDERATED_SSO_URL>
  export MONGODB_ATLAS_FEDERATED_ISSUER_URI=<ATLAS_FEDERATED_ISSUER_URI>
  ```
~> **Notice:** For more information about the Federation configuration resource, see: https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-federated-authentication

##### AWS env variables

- For `Network Peering` resource configuration:
  ```bash
  export AWS_ACCOUNT_ID=<AWS_ACCOUNT_ID>
  export AWS_VPC_ID=<AWS_VPC_ID>
  export AWS_VPC_CIDR_BLOCK=<AWS_VPC_CIDR_BLOCK>
  export AWS_REGION=<AWS_REGION>
  export AWS_SUBNET_ID=<AWS_SUBNET_ID>
  export AWS_SECURITY_GROUP_ID=<AWS_SECURITY_GROUP_ID>
  ```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/

- For `Encryption at Rest` resource configuration:
  ```bash
  export AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID>
  export AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY>
  export AWS_CUSTOMER_MASTER_KEY_ID=<AWS_CUSTOMER_MASTER_KEY_ID>
  export AWS_REGION=<AWS_REGION>

  export AWS_ACCESS_KEY_ID_UPDATED=<AWS_ACCESS_KEY_ID_UPDATED>
  export AWS_SECRET_ACCESS_KEY_UPDATED=<AWS_SECRET_ACCESS_KEY_UPDATED>
  export AWS_CUSTOMER_MASTER_KEY_ID_UPDATED=<AWS_CUSTOMER_MASTER_KEY_ID_UPDATED>
  export AWS_REGION_UPDATED=<AWS_REGION_UPDATED>
  ```
~> **Notice:** For more information about the Encryption at Rest resource, see: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/

- For `Private Endpoint Link` resource configuration:
  ```bash
  export AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID>
  export AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY>
  export AWS_CUSTOMER_MASTER_KEY_ID=<AWS_CUSTOMER_MASTER_KEY_ID>
  export AWS_REGION=<AWS_REGION>
  export AWS_VPC_ID=<AWS_VPC_ID>
  export AWS_SUBNET_ID=<AWS_SUBNET_ID>
  export AWS_SECURITY_GROUP_ID=<AWS_SECURITY_GROUP_ID>
  ```
~> **Notice:** For more information about the PrivateLink (for AWS only), see: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/https://docs.atlas.mongodb.com/reference/api/private-endpoint/

##### AZURE env variables

- For `Network Peering` resource configuration:
  ```bash
  export AZURE_DIRECTORY_ID=<AZURE_DIRECTORY_ID>
  export AZURE_SUBSCRIPTION_ID=<AZURE_SUBSCRIPTION_ID>
  export AZURE_RESOURCE_GROUP_NAME=<AZURE_RESOURCE_GROUP_NAME>
  export AZURE_VNET_NAME=<AZURE_VNET_NAME>
  ```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/


- For Encryption at Rest resource configuration:
  ```bash
  export AZURE_CLIENT_ID=<AZURE_CLIENT_ID>
  export AZURE_SUBSCRIPTION_ID=<AZURE_SUBSCRIPTION_ID>
  export AZURE_RESOURCE_GROUP_NAME=<AZURE_RESOURCE_GROUP_NAME>
  export AZURE_APP_SECRET=<AZURE_APP_SECRET>
  export AZURE_KEY_VAULT_NAME=<AZURE_KEY_VAULT_NAME>
  export AZURE_KEY_IDENTIFIER=<AZURE_KEY_IDENTIFIER>
  export AZURE_TENANT_ID=<AZURE_TENANT_ID>
  export AZURE_DIRECTORY_ID=<AZURE_DIRECTORY_ID>
  export AZURE_KEY_VAULT_NAME_UPDATED=<AZURE_KEY_VAULT_NAME_UPDATED>
  export AZURE_KEY_IDENTIFIER_UPDATED=<AZURE_KEY_IDENTIFIER_UPDATED>
  ```
~> **Notice:** For more information about the Encryption at Rest resource, see: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/

##### GCP env variables
- For `Network Peering` resource configuration:
  ```bash
  $export GCP_PROJECT_ID=<GCP_PROJECT_ID>
  ```
~> **Notice:** For more information about the Network Peering resource, see: https://docs.atlas.mongodb.com/reference/api/vpc/


- For Encryption at Rest resource configuration:
  ```bash
  export GCP_SERVICE_ACCOUNT_KEY=<GCP_SERVICE_ACCOUNT_KEY>
  export GCP_KEY_VERSION_RESOURCE_ID=<GCP_KEY_VERSION_RESOURCE_ID>
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
        path    = "C:\Users\<YourUser>\Desktop\Tenant_Upgrade\tf_cache"
        include = ["registry.terraform.io/hashicorp/*"]
      }
      direct {
        exclude = ["registry.terraform.io/hashicorp/*"]
      }
    }
    plugin_cache_dir = "C:\Users\<YourUser>\Desktop\Tenant_Upgrade\tf_cache"
    disable_checkpoint=true
    ```
-	`cd` back up to original directory and `mv` the `.terraform/providers/registry.terraform.io` directory to `tf_cache` 

-	Create required environment variables (modify accordingly to your own local path directory):
    ```bash
    export TF_PLUGIN_CACHE_DIR=/mnt/c/Users/<YourUser>/Desktop/Tenant_Upgrade/tf_cache
    export TF_CLI_CONFIG_FILE=/mnt/c/Users/<YourUser>/Desktop/Tenant_Upgrade/tf_cache/terraform.rc
    ```
-  Delete the `.terraform` and `.terraform.lock.hcl` directories altogether. At this point you should only have the `tf_cache` directory and the `versions.tf` config file remaining. 

- Next in the `tf_cache` directory replace existing terraform provider core files (Terraform Atlas Provider version binary, `CHANGELOG.md`, `LICENSE`, and `README.md`) with the version you seek to test locally. Make sure to keep folder structure the same. 

- Lastly, in the terminal run `terraform init` again and this time terraform will pull provider version from `tf_cache` network mirror. You can confirm this by noting the `Terraform has been successfully initialized! Using mongodb/mongodbatlas Vx.x.x from the shared cache directory` message.  
