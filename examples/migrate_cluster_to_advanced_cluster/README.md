# How to Transition from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`

This directory contains examples demonstrating how to transition from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster` using the `moved block`. For more details, please refer to the [Migration Guide: Cluster to Advanced Cluster](TODO-LINK-HERE).

The examples are organized as follows:
- **For users directly utilizing the `mongodbatlas_cluster` resource**: please check the [basic/](./basic/README.md) folder.
- **For users employing `modules` to manage `mongodbatlas_cluster`**: please see the [module_maintainer/](./module_maintainer/README.md) and [module_user/](./module_user/README.md) folders. These folders illustrate the migration process from both the maintainer's and the user's perspectives, highlighting how the migration can be executed in phases to manage breaking changes that may affect module users.