# Module Maintainer - `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`

The purpose of this example is to demonstrate how a Terraform module can help in the migration from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`.
The example contains three module versions:

<!-- Update Variable count -->
Version | Purpose | Variables | Resources
--- | --- | --- | ---
[v1](./v1) | Baseline | 5 | `mongodbatlas_cluster`
[v2](./v2) | Migrate to advanced_cluster with no change in variables or plan | 5 | `mongodbatlas_advanced_cluster`
[v3](./v3) | Use the latest features of advanced_cluster | 10 | `mongodbatlas_advanced_cluster`

## `v2` Implementation Changes and Highlights

### `variables.tf` unchanged from `v1`
### `versions.tf`
- `required_version` of Terraform CLI bumped to `# todo: minimum moved block supported version` for `moved` block support
- `mongodbatlas.version` bumped to `# todo: PLACEHOLDER_TPF_RELEASE` for new `mongodbatlas_advanced_cluster` v2 schema support

### `main.tf`
<!-- TODO: Update link to (schema v2) docs page -->
- `locals.replication_specs` an intermediate variable transforming the `variables` to a compatible [replication_specs](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#replication_specs-1) for `mongodbatlas_advanced_cluster`
  - We use the Terraform builtin [range](https://developer.hashicorp.com/terraform/language/functions/range) function (`range(old_spec.num_shards)`) to flatten `num_shards`
  - We expand `read_only_specs` and `electable_specs` into nested attributes
  - We use the `var.provider_name` in the `region_configs.*.instance_size`
- `moved` block:
```terraform
moved {
  from = mongodbatlas_cluster.this
  to   = mongodbatlas_advanced_cluster.this
}
```
- `resource "mongodbatlas_advanced_cluster" "this"`
  - We reference the `local.replication_specs` as input to `replication_specs` (`replication_specs = local.replication_specs`)
  - Tags can be passed directly instead of the dynamic block (`tags = var.tags`)
- Adds `data "mongodbatlas_cluster" "this"` to avoid breaking changes in `outputs.tf` (see below)

### `outputs.tf`
- All outputs can use `mongodbatlas_advanced_cluster` except for
  - `replication_specs`, we use the `data.mongodbatlas_cluster.this.replication_specs` to keep the same format
  - `mongodbatlas_cluster`, we use the `data.mongodbatlas_cluster.this` to keep the same format


## `v3` Implementation Changes and Highlights
