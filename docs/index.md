# MongoDB Atlas Provider

Use the MongoDB Atlas provider to interact with the resources supported by [MongoDB Atlas](https://www.mongodb.com/cloud/atlas). The provider needs to be configured with proper credentials before it can be used.

Use the navigation to the left to read about the available provider resources and data sources.

## Quick Start

This example shows how to set up the MongoDB Atlas provider and create your first cluster:

```terraform
# Configure the MongoDB Atlas Provider with Service Account credentials
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}

# Create a project
resource "mongodbatlas_project" "test" {
  name   = "test-project"
  org_id = var.org_id
}

# Create a cluster
resource "mongodbatlas_cluster" "test" {
  project_id = mongodbatlas_project.test.id
  name       = "test-cluster"

  # Minimum tier for a dedicated cluster
  cluster_type = "REPLICASET"

  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }

  cloud_backup = true
  auto_scaling_disk_gb_enabled = true

  # Provider settings
  provider_name               = "AWS"
  provider_instance_size_name = "M10"
}
```

## Authentication

The MongoDB Atlas provider uses Service Accounts as the recommended authentication method. You need to create a Service Account in your MongoDB Atlas organization and grant it appropriate permissions.

### Setting up authentication:

1. Create a Service Account following the [MongoDB Atlas documentation](https://www.mongodb.com/docs/atlas/configure-api-access/#grant-programmatic-access-to-an-organization)
2. Configure the provider with your credentials:

**Using environment variables (recommended):**
```shell
export MONGODB_ATLAS_CLIENT_ID="your-client-id"
export MONGODB_ATLAS_CLIENT_SECRET="your-client-secret"
```

Then in your Terraform configuration:
```terraform
provider "mongodbatlas" {
  # Credentials are read from environment variables
}
```

**Using provider attributes:**
```terraform
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}
```

If you need to use MongoDB Atlas for Government, see the [Provider Configuration Guide](guides/provider-configuration#mongodb-atlas-for-government) for setup instructions.

For detailed information about all authentication methods, including Programmatic Access Keys, AWS Secrets Manager integration, and migration guides, see the [Provider Configuration Guide](guides/provider-configuration).

## Provider Configuration

### Arguments

* `client_id` - (Optional) Service Account Client ID. Can also be set with the `MONGODB_ATLAS_CLIENT_ID` environment variable.
* `client_secret` - (Optional) Service Account Client Secret. Can also be set with the `MONGODB_ATLAS_CLIENT_SECRET` environment variable.
* `access_token` - (Optional) Service Account Access Token. Can also be set with the `MONGODB_ATLAS_ACCESS_TOKEN` environment variable. Note: tokens have expiration times.
* `public_key` - (Optional) MongoDB Atlas Programmatic Access Key Public Key. Can also be set with the `MONGODB_ATLAS_PUBLIC_API_KEY` environment variable.
* `private_key` - (Optional) MongoDB Atlas Programmatic Access Key Private Key. Can also be set with the `MONGODB_ATLAS_PRIVATE_API_KEY` environment variable.
* `base_url` - (Optional) MongoDB Atlas Base URL. Can also be set with the `MONGODB_ATLAS_BASE_URL` environment variable.
* `realm_base_url` - (Optional) MongoDB Realm Base URL. Can also be set with the `MONGODB_REALM_BASE_URL` environment variable.
* `is_mongodbgov_cloud` - (Optional) Set to true to use MongoDB Atlas for Government.
* `assume_role` - (Optional) AWS IAM role configuration for accessing secrets in AWS Secrets Manager. See [Provider Configuration Guide](guides/provider-configuration#aws-secrets-manager) for details.
* `secret_name` - (Optional) Name of the secret in AWS Secrets Manager.
* `region` - (Optional) AWS region where the secret is stored.
* `aws_access_key_id` - (Optional) AWS Access Key ID. Can also be set with the `AWS_ACCESS_KEY_ID` environment variable.
* `aws_secret_access_key` - (Optional) AWS Secret Access Key. Can also be set with the `AWS_SECRET_ACCESS_KEY` environment variable.
* `aws_session_token` - (Optional) AWS Session Token. Can also be set with the `AWS_SESSION_TOKEN` environment variable.
* `sts_endpoint` - (Optional) AWS STS endpoint. Can also be set with the `STS_ENDPOINT` environment variable.

## Version Requirements

### Provider Version

We recommend pinning your provider version to at least the major version (e.g., `~> 2.0`) to avoid accidental upgrades to incompatible new versions:

```terraform
terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 2.0"
    }
  }
}
```

Starting with version 2.0.0, the [MongoDB Atlas Provider Versioning Policy](#mongodb-atlas-provider-versioning-policy) ensures that minor and patch versions do not include breaking changes.

### Terraform Version

We recommend using the latest [HashiCorp Terraform Core Version](https://github.com/hashicorp/terraform). See the [HashiCorp Terraform Version Compatibility Matrix](#hashicorp-terraform-version-compatibility-matrix) below for supported versions.

## MongoDB Atlas Provider Versioning Policy

### Versioning Strategy

We follow [semantic versioning](https://semver.org/):
- **Major (X.0.0):** Introduces breaking changes
- **Minor (X.Y.0):** Adds non-breaking changes and deprecations
- **Patch (X.Y.Z):** Bug fixes and documentation updates

### Release Cadence

- **Minor and patch** versions: Biweekly releases
- **Major** versions: Once per year (maximum two per calendar year)
- **Off-cycle releases** may occur for critical security flaws or regressions

### Breaking Changes and Deprecation

- Breaking changes are announced as deprecations in minor versions
- Deprecated functionality is removed in the next 1-2 major versions
- Each release includes a [changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/releases)
- Major versions include migration guides

For the full versioning policy details, see [MongoDB Atlas Provider Versioning Policy](#mongodb-atlas-provider-versioning-policy) section below.

## HashiCorp Terraform Version Compatibility Matrix

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

For safety, we require only consuming versions of HashiCorp Terraform that are currently receiving Security/Maintenance updates. See [Support Period and End-of-Life Policy](https://support.hashicorp.com/hc/en-us/articles/360021185113-Support-Period-and-End-of-Life-EOL-Policy).

## Supported OS and Architectures

The MongoDB Atlas Provider supports multiple operating systems and architectures. See the [Provider Configuration Guide](guides/provider-configuration#supported-os-and-architectures) for the complete list of supported platforms.

## Resources

### Documentation and Examples
- [Example configurations](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.1/examples) - Beginner and advanced examples
- [Terraform Modules](https://registry.terraform.io/namespaces/terraform-mongodbatlas-modules) - Pre-built modules for common patterns
- [MongoDB Atlas and Terraform Landing Page](https://www.mongodb.com/atlas/terraform)
- [CHANGELOG](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/CHANGELOG.md) - Current version information

### Support
- [Report bugs](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)
- [Request features](https://feedback.mongodb.com/forums/924145-atlas?category_id=370723)
- [MongoDB Atlas support plans](https://docs.atlas.mongodb.com/support/) - Developer tier and above

### Community
Have a good example you've created? Submit a PR to add it to the `examples` directory in our [GitHub repo](https://github.com/mongodb/terraform-provider-mongodbatlas/).

---

## MongoDB Atlas Provider Versioning Policy (Full Details)

### Definition of Breaking Changes

Breaking changes are defined as any change that requires user intervention to address. This may include:
- Modifying existing schema (e.g., removing or renaming fields, renaming resources)
- Changes to business logic (e.g., implicit default values or server-side behavior)
- Provider-level changes (e.g., changing retry behavior)

Final confirmation of a breaking change — possibly leading to an exemption — is subject to:
- MongoDB's understanding of the adoption level of the feature
- Timing of the next planned major release
- The change's relation to a bug fix

### Customer Communication

We are committed to clear and proactive communication:
- **Each release** includes a [changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/releases) clearly labeling: `breaking`, `deprecated`, `bug-fix`, `feature`, and `enhancement` changes
- **Major versions** include migration guides
- **Minor and patch versions** generally do not include migration guides, but may if warranted
- **GitHub tags** with `vX.Y.Z` format are provided for all releases