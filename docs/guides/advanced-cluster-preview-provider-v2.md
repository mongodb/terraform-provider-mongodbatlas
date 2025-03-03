---
page_title: "Migration Guide: Advanced Cluster (Preview for MongoDB Atlas Provider 2.0.0)"
---

# Migration Guide: Advanced Cluster (Preview for MongoDB Atlas Provider 2.0.0)

**Objective**: This guide explains the changes introduced for the `mongodbatlas_advanced_cluster` resource in the Preview for MongoDB Atlas Provider 2.0.0 of the and how to migrate to it.

 `mongodbatlas_advanced_cluster` in the Preview for MongoDB Atlas Provider 2.0.0 is implemented using the recommended [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework). This improves the overall user experience and provides a more consistent and predictable behavior. It supports the latest Terraform features and best practices, including support for `moved` block between different resource types, for more information see the [Migration Guide: Cluster to Advanced Cluster](cluster-to-advanced-cluster-migration-guide#moved-block).

The [resource documentation page](../resources/advanced_cluster%2520%2528preview%2520provider%2520v2%2529) contains all the details about the `mongodbatlas_advanced_cluster` resource in the Preview for MongoDB Atlas Provider 2.0.0.

## Enable the Preview for MongoDB Atlas Provider 2.0.0 for `mongodbatlas_advanced_cluster`

In order to enable the Preview for MongoDB Atlas Provider 2.0.0 for `mongodbatlas_advanced_cluster`, set the environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`. This will allow you to use the new `mongodbatlas_advanced_cluster` resource. You can also define the environment variable in your local development environment so your tools can use the new format and help you with linting and auto-completion.

This environment variable only affects the `mongodbatlas_advanced_cluster` resource and corresponding data sources. It doesn't affect other resources. `mongodbatlas_advanced_cluster` definition will use the new format and new features like `moved block` from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster` will be available.


## Best Practices Before Migrating
Before doing any migration create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state).

## How to migrate `mongodbatlas_advanced_cluster` to Preview for MongoDB Atlas Provider 2.0.0 

The process to migrate from current `mongodbatlas_advanced_cluster` to the one in Preview for MongoDB Atlas Provider 2.0.0 is as follows:
- Before starting, run `terraform plan` to make sure that there are no planned changes.
- Set environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true` in order to use the Preview for MongoDB Atlas Provider 2.0.0.
- Run `terraform plan` and you'll see errors as definition file hasn't been updated yet.
- Apply definition changes explained on this page until there are no errors and no planned changes. **Important**: Don't apply until the plan is empty. If it shows other changes, you must update the `mongodbatlas_advanced_cluster` configuration until it matches the original configuration.
- Run `terraform apply` to apply the changes. Although there are no plan changes shown to the user, the `mongodbatlas_advanced_cluster` state will be updated to support the Preview for MongoDB Atlas Provider 2.0.0.
