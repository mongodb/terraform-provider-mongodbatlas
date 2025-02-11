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

File | Highlights | Code
--- | --- | ---
variables.tf | • Added `default` | <pre> variable "instance_size" { <br>  type    = string <br>  default = "" # optional in v3 <br>} <br>variable "provider_name" { <br>  type    = string <br>  default = "" # optional in v3 <br>} </pre>
variables.tf | • Using a default of empty `[]` for `replication_specs` | <pre>variable "replication_specs" {<br>  description = "List of replication specifications in legacy mongodbatlas_cluster format"<br>  default     = []<br>  # everything else the same<br>} </pre>
variables.tf | • New `replication_specs_new` variable<br> • Usage of `optional` to simplify caller<br> • Using `[]` for default | <pre>variable "replication_specs_new" {<br>  description = "List of replication specifications using new mongodbatlas_advanced_cluster format"<br>  default     = []<br>  type = list(object({<br>    zone_name = optional(string, "Zone 1")<br><br>    region_configs = list(object({<br>      region_name   = string<br>      provider_name = string<br>      priority      = optional(number, 7)<br><br>      auto_scaling = optional(object({<br>        disk_gb_enabled = optional(bool, false)<br>      }), null)<br><br>      read_only_specs = optional(object({<br>        node_count      = number<br>        instance_size   = string<br>        disk_size_gb    = optional(number, null)<br>        ebs_volume_type = optional(string, null)<br>        disk_iops       = optional(number, null)<br>      }), null)<br>      analytics_specs = optional(object({<br>        node_count      = number<br>        instance_size   = string<br>        disk_size_gb    = optional(number, null)<br>        ebs_volume_type = optional(string, null)<br>        disk_iops       = optional(number, null)<br>      }), null)<br>      electable_specs = object({<br>        node_count      = number<br>        instance_size   = string<br>        disk_size_gb    = optional(number, null)<br>        ebs_volume_type = optional(string, null)<br>        disk_iops       = optional(number, null)<br>      })<br>    }))<br>  }))<br>} </pre>
main.tf | • Add *defaults* to old variables in `locals`<br> • Add `_old` suffix to `locals.replication_specs` | <pre>  old_disk_size     = var.auto_scaling_disk_gb_enabled ? null : var.disk_size<br>  old_instance_size = coalesce(var.instance_size, "M10")<br>  old_provider_name = coalesce(var.provider_name, "AWS")<br>  replication_specs_old = flatten([<br> </pre>
main.tf | • Add `precondition` for `replication_specs`<br> • Use a conditional for `replication_specs` | <pre>resource "mongodbatlas_advanced_cluster" "this" {<br>  lifecycle {<br><br>    precondition {<br>      condition     = local.use_new_replication_specs &#124;&#124; !(var.auto_scaling_disk_gb_enabled && var.disk_size &gt; 0)<br>      error_message = "Must use either auto_scaling_disk_gb_enabled or disk_size, not both."<br>    }<br>    precondition {<br>      condition     = !((local.use_new_replication_specs && length(var.replication_specs) &gt; 0) &#124;&#124; (!local.use_new_replication_specs && length(var.replication_specs) == 0))<br>      error_message = "Must use either replication_specs_new or replication_specs, not both."<br>    }<br>  }<br><br>  project_id             = var.project_id<br>  name                   = var.cluster_name<br>  cluster_type           = var.cluster_type<br>  mongo_db_major_version = var.mongo_db_major_version<br>  replication_specs      = local.use_new_replication_specs ? var.replication_specs_new : local.replication_specs_old<br>  tags                   = var.tags<br>}<br> </pre>
