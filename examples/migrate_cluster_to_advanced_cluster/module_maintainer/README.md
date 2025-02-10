# Module Maintainer - `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`

The purpose of this example is to demonstrate how a Terraform module can help the migration from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`.
The example contains three module versions:

<!-- Update Variable count -->
Version | Purpose | Variables | Resources
--- | --- | --- | ---
[v1](./v1) | Baseline | 5 | `mongodbatlas_cluster`
[v2](./v2) | Migrate to advanced_cluster with no change in variables or plan | 5 | `mongodbatlas_advanced_cluster`
[v3](./v3) | Use the latest features of advanced_cluster | 10 | `mongodbatlas_advanced_cluster`

<!-- TODO: Add highlights of module implementations -->
