# MongoDB Atlas Provider

The MongoDB Atlas provider is used to interact with the resources supported by [MongoDB Atlas](https://www.mongodb.com/cloud/atlas). The provider needs to be configured with proper credentials before it can be used.

## Example Usage

This example shows how to set up the MongoDB Atlas provider and create a cluster:

```terraform
provider "mongodbatlas" {
  client_id     = var.mongodbatlas_client_id
  client_secret = var.mongodbatlas_client_secret
}

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

The MongoDB Atlas provider uses Service Account (SA) as the recommended authentication method.

For detailed authentication configuration, see:
- [Service Account (SA)](guides/provider-configuration#service-account-recommended)
- [Programmatic Access Key (PAK)](guides/provider-configuration#programmatic-access-key)
- [AWS Secrets Manager integration](guides/provider-configuration#aws-secrets-manager)

## MongoDB Atlas for Government

MongoDB Atlas for Government is a dedicated deployment option for government agencies and contractors requiring FedRAMP compliance. 
For more details on configuration, see the [Provider Configuration Guide](guides/provider-configuration#provider-arguments).

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

In order to promote stability, predictability, and transparency, the MongoDB Atlas Terraform Provider will implement **semantic versioning** with a **scheduled release cadence**. Our goal is to deliver regular improvements to the provider without overburdening users with frequent breaking changes.

---

### Definition of Breaking Changes

Our definition of breaking changes aligns with the impact updates have on the customer:

Breaking changes are defined as any change that requires user intervention to address.
This may include:

- Modifying existing schema (e.g., removing or renaming fields, renaming resources)
- Changes to business logic (e.g., implicit default values or server-side behavior)
- Provider-level changes (e.g., changing retry behavior)

Final confirmation of a breaking change — possibly leading to an exemption — is subject to:

- MongoDB’s understanding of the adoption level of the feature
- Timing of the next planned major release
- The change's relation to a bug fix

---

### Versioning Strategy

We follow [semantic versioning](https://semver.org/) for all updates:

- **Major (X.0.0):** Introduces breaking changes (as defined by MongoDB)
- **Minor (X.Y.0):** Adds non-breaking changes and announces deprecations
- **Patch (X.Y.Z):** Includes bug fixes and documentation updates

We do not utilize pre-release versioning at this time.

---

### Release Cadence

To minimize unexpected changes, we follow a scheduled cadence:

- **Minor and patch** versions follow a **biweekly** release pattern
- **Major** versions are released **once per year**, with a maximum of **two per calendar year**
- The provider team may adjust the schedule based on need

**Off-cycle releases** may occur for critical security flaws or regressions.

---

### Deprecation Policy

We use a structured deprecation window to notify customers in advance:

- Breaking changes are **deprecated in a minor version** with:
  - Warnings in migration guides, changelogs, and resource usage
- Deprecated functionality is **removed in the next 1–2 major versions**, unless otherwise stated

---

### Customer Communication

We are committed to clear and proactive communication:

- **Each release** includes a [changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/releases) clearly labeling:
  - `breaking`, `deprecated`, `bug-fix`, `feature`, and `enhancement` changes
- **Major versions** include migration guides
- **Minor and patch versions** generally do not include migration guides, but may if warranted
- **GitHub tags** with `vX.Y.Z` format are provided for all releases

---

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
<!-- MATRIX_PLACEHOLDER_END -->
For the safety of our users, we require only consuming versions of HashiCorp Terraform that are currently receiving Security / Maintenance Updates. For more details see [Support Period and End-of-Life (EOL) Policy](https://support.hashicorp.com/hc/en-us/articles/360021185113-Support-Period-and-End-of-Life-EOL-Policy).   

HashiCorp Terraform versions that are not listed on this table are no longer supported by MongoDB Atlas. For latest HashiCorp Terraform versions see [here](https://endoflife.date/terraform ).

## Supported OS and Architectures

The MongoDB Atlas Provider supports multiple operating systems and architectures. See the [Provider Configuration Guide](guides/provider-configuration#supported-os-and-architectures) for the complete list of supported platforms.

## Helpful Links/Information

[Upgrade Guide for Terraform MongoDB Atlas 0.4.0](https://www.mongodb.com/blog/post/upgrade-guide-for-terraform-mongodb-atlas-040)

[MongoDB Atlas and Terraform Landing Page](https://www.mongodb.com/atlas/terraform)

[Report bugs](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)

[Request Features](https://feedback.mongodb.com/forums/924145-atlas?category_id=370723)

[Support covered by MongoDB Atlas support plans, Developer and above](https://docs.atlas.mongodb.com/support/) 

## Examples from MongoDB and the Community

<!-- NOTE: the below examples link is updated during the release process, when doing changes in the following sentence verify scripts/update-examples-reference-in-docs.sh is not impacted-->
We have [example configurations](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.1.0/examples)
in our GitHub repo that will help both beginner and more advanced users.

Have a good example you've created and want to share?
Let us know the details via an [issue](https://github.com/mongodb/terraform-provider-mongodbatlas/issues)
or submit a PR of your work to add it to the `examples` directory in our [GitHub repo](https://github.com/mongodb/terraform-provider-mongodbatlas/).

## Terraform MongoDB Atlas Modules
You can now leverage our [Terraform Modules](https://registry.terraform.io/namespaces/terraform-mongodbatlas-modules) to easily get started with MongoDB Atlas and critical features like [Push-based log export](https://registry.terraform.io/modules/terraform-mongodbatlas-modules/push-based-log-export/mongodbatlas/latest), [Private Endpoints](https://registry.terraform.io/modules/terraform-mongodbatlas-modules/private-endpoint/mongodbatlas/latest), etc.
