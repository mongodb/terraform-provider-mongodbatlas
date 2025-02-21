# Module Maintainer - `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`

The purpose of this example is to demonstrate how a Terraform module definition can migrate from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster` while minimizing impact to its clients.
The example contains three module versions which represent the three steps of the migration:

Step | Purpose | Resources
--- | --- | ---
[Step 1](./v1) | Baseline | `mongodbatlas_cluster`
[Step 2](./v2) | Migrate to advanced_cluster with no change in variables or plan | `mongodbatlas_advanced_cluster`
[Step 3](./v3) | Use the latest features of advanced_cluster | `mongodbatlas_advanced_cluster`

The rest of this document summarizes the different implementations:

- [Step 1: Module `v1` Implementation Summary](#step-1-module-v1-implementation-summary)
  - [`variables.tf`](#variablestf)
  - [`main.tf`](#maintf)
  - [`outputs.tf`](#outputstf)
- [Step 2: Module `v2` Implementation Changes and Highlights](#step-2-module-v2-implementation-changes-and-highlights)
  - [`variables.tf` unchanged from `v1`](#variablestf-unchanged-from-v1)
  - [`versions.tf`](#versionstf)
  - [`main.tf`](#maintf-1)
  - [`outputs.tf`](#outputstf-1)
- [Step 3: Module `v3` Implementation Changes and Highlights](#step-3-module-v3-implementation-changes-and-highlights)
  - [`variables.tf`](#variablestf-1)
  - [`main.tf`](#maintf-2)
  - [`outputs.tf`](#outputstf-2)


## Step 1: Module `v1` Implementation Summary
This module creates a `mongodbatlas_cluster`.

### [`variables.tf`](v1/variables.tf)
An abstraction for the `mongodbatlas_cluster` resource:
- Not all arguments are exposed, but the arguments follow the schema closely.
- `disk_size` and `auto_scaling_disk_gb_enabled` are mutually exclusive and validated in the `precondition` in `main.tf`.

### [`main.tf`](v1/main.tf)
It uses `dynamic` blocks to represent:
- `tags`
- `replication_specs`
- `regions_config` (nested inside replication_specs)

### [`outputs.tf`](v1/outputs.tf)
- Expose some attributes of `mongodbatlas_cluster` but also the full resource with `mongodbatlas_cluster` output variable:
```terraform
output "mongodbatlas_cluster" {
  value       = mongodbatlas_cluster.this
  description = "Full cluster configuration for mongodbatlas_cluster resource"
}
```

## Step 2: Module `v2` Implementation Changes and Highlights
This module uses HCL code to create a `mongodbatlas_advanced_cluster` resource that is compatible with the input variables of `v1`.
The module supports standalone usage when there is no existing `mongodbatlas_cluster` and also upgrading from `v1` using a `moved` block.

### [`variables.tf`](v2/variables.tf) unchanged from `v1`
### [`versions.tf`](v2/versions.tf)
- `required_version` of Terraform CLI bumped to `1.8` for `moved` block [support](https://developer.hashicorp.com/terraform/plugin/framework/resources/state-move) between resource types.
- `mongodbatlas.version` bumped to `1.27.0` for new `mongodbatlas_advanced_cluster` v2 schema support.

### [`main.tf`](v2/main.tf)
- `locals.replication_specs` an intermediate variable transforming the `variables` to a compatible [replication_specs](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster%2520%2528preview%2520provider%2520v2%2529#replication_specs-1) for `mongodbatlas_advanced_cluster`
  - We use the Terraform builtin [range](https://developer.hashicorp.com/terraform/language/functions/range) function (`range(old_spec.num_shards)`) to flatten `num_shards`.
  - We expand `read_only_specs` and `electable_specs` into nested attributes.
  - We use the `var.provider_name` in the `region_configs.*.instance_size`
- `moved` block:
```terraform
moved {
  from = mongodbatlas_cluster.this
  to   = mongodbatlas_advanced_cluster.this
}
```
- `resource "mongodbatlas_advanced_cluster" "this"`
  - We reference the `local.replication_specs` as input to `replication_specs` (`replication_specs = local.replication_specs`).
  - Tags can be passed directly instead of the dynamic block (`tags = var.tags`).
- Adds `data "mongodbatlas_cluster" "this"` to avoid breaking changes in `outputs.tf` (see below).

### [`outputs.tf`](v2/outputs.tf)
- All outputs can use `mongodbatlas_advanced_cluster` except for:
  - `replication_specs`, which uses `data.mongodbatlas_cluster.this.replication_specs` to keep the same format.
  - `mongodbatlas_cluster`, which uses the `data.mongodbatlas_cluster.this` to keep the same format.


## Step 3: Module `v3` Implementation Changes and Highlights
This module adds variables to support the latest `mongodbatlas_advanced_cluster` features while staying compatible with the old input variables.
The module supports standalone usage when there is no existing `mongodbatlas_cluster` and also upgrading from `v1` using a `moved` block.
The module also supports changing an existing `mongodbatlas_advanced_cluster` created in `v2`.

### [`variables.tf`](v3/variables.tf)
- Add `replication_specs_new`. This is almost fully equivalent to the [`replication_specs`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster%2520%2528preview%2520provider%2520v2%2529#replication_specs-1) of the latest `mongodbatlas_advanced_cluster` schema.
  - Use a `[]` for default to allow continued usage of the old `replication_specs`.
  - Usage of `optional` to simplify the caller.
- Add `default` to `instance_size` and `provider_name` as these are not required when `replication_specs_new` is used.
- Change `[]` default to `replication_specs` to allow usage of `replication_specs_new`.

### [`main.tf`](v3/main.tf)
- Add *defaults* to old variables in `locals`:
  - `old_disk_size`
  - `old_instance_size`
  - `old_provider_name`
- Add `_old` suffix to `locals.replication_specs` to make conditional code (see below) more readable.
- Add `precondition` for `replication_specs` to validate only `var.replication_specs_new` or `replication_specs` is used.
```terraform
    precondition {
      condition     = !((local.use_new_replication_specs && length(var.replication_specs) > 0) || (!local.use_new_replication_specs && length(var.replication_specs) == 0))
      error_message = "Must use either replication_specs_new or replication_specs, not both."
    }
```
- Use a conditional for`replication_specs` in `resource "mongodbatlas_advanced_cluster" "this"`:
```terraform
  # other attributes...
  replication_specs      = local.use_new_replication_specs ? var.replication_specs_new : local.replication_specs_old
  tags                   = var.tags
}
```
- Use `count` for data source to avoid error when Asymmetric Shards are used:
```terraform
data "mongodbatlas_cluster" "this" {
  count      = local.use_new_replication_specs ? 0 : 1 # Not safe when Asymmetric Shards are used
  name       = mongodbatlas_advanced_cluster.this.name
  project_id = mongodbatlas_advanced_cluster.this.project_id

  depends_on = [mongodbatlas_advanced_cluster.this]
}
```

### [`outputs.tf`](v3/outputs.tf)
- Update `replication_specs` and `mongodbatlas_cluster` to handle the case when the new schema is used:
```terraform
output "replication_specs" {
  value       = local.use_new_replication_specs ? [] : data.mongodbatlas_cluster.this[0].replication_specs # updated
  description = "Replication Specs for cluster, will be empty if var.replication_specs_new is set"
}

output "mongodbatlas_cluster" {
  value       = local.use_new_replication_specs ? null : data.mongodbatlas_cluster.this[0] # updated
  description = "Full cluster configuration for mongodbatlas_cluster resource, will be null if var.replication_specs_new is set"
}
```
