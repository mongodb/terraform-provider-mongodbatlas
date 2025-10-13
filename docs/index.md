# MongoDB Atlas Provider

The MongoDB Atlas provider is used to interact with the resources supported by [MongoDB Atlas](https://www.mongodb.com/cloud/atlas). The provider needs to be configured with proper credentials before it can be used.

## Example Usage

This example shows how to set up the MongoDB Atlas provider and create a cluster:

```terraform
# Configure the MongoDB Atlas Provider
# Service Account credentials are read from environment variables:
# - MONGODB_ATLAS_CLIENT_ID
# - MONGODB_ATLAS_CLIENT_SECRET
provider "mongodbatlas" {}

# Create a project
resource "mongodbatlas_project" "this" {
  name   = "my-project"
  org_id = var.org_id
}

# Create a cluster
resource "mongodbatlas_advanced_cluster" "this" {
  project_id   = mongodbatlas_project.this.id
  name         = "my-cluster"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          region_name   = "US_EAST_1"
          priority      = 7
          provider_name = "AWS"
          electable_specs = {
            instance_size = "M10"
            node_count    = 3
          }
        }
      ]
    }
  ]
}
```

## Authentication

The MongoDB Atlas provider uses Service Accounts (SA) as the recommended authentication method. Create an SA in your [MongoDB Atlas organization](https://www.mongodb.com/docs/atlas/configure-api-access/#grant-programmatic-access-to-an-organization) and set the credentials as environment variables:

```shell
export MONGODB_ATLAS_CLIENT_ID="your-client-id"
export MONGODB_ATLAS_CLIENT_SECRET="your-client-secret"
```

For detailed authentication configuration including Programmatic Access Key (PAK), AWS Secrets Manager integration, and MongoDB Atlas for Government, see the [Provider Configuration Guide](guides/provider-configuration).

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

We follow [semantic versioning](https://semver.org/):
- **Major (X.0.0):** Introduces breaking changes
- **Minor (X.Y.0):** Adds non-breaking changes and deprecations
- **Patch (X.Y.Z):** Bug fixes and documentation updates

Release Cadence:
- **Minor and patch** versions: Biweekly releases
- **Major** versions: Once per year (maximum two per calendar year)
- **Off-cycle releases** may occur for critical security flaws or regressions

Breaking Changes and Deprecation:

- Breaking changes are announced as deprecations in minor versions
- Deprecated functionality is removed in the next 1-2 major versions
- Each release includes a [changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/releases)
- Major versions include migration guides

For complete details including breaking changes policy and deprecation guidelines, see the [full versioning policy](#mongodb-atlas-provider-versioning-policy-full-details) below.

## [HashiCorp Terraform Version](https://www.terraform.io/downloads.html) Compatibility Matrix

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
