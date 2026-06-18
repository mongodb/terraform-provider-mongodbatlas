---
subcategory: "Organizations"
---

# Data Source: mongodbatlas_org_maintenance_settings

`mongodbatlas_org_maintenance_settings` provides a data source to read the organization-level [maintenance wave settings](https://www.mongodb.com/docs/atlas/tutorial/cluster-maintenance-window/) for a MongoDB Atlas organization.

## Example Usage

```terraform
data "mongodbatlas_org_maintenance_settings" "example" {
  org_id = var.org_id
}

output "maintenance_settings" {
  value = {
    wave_assignment_mode           = data.mongodbatlas_org_maintenance_settings.example.wave_assignment_mode
    effective_wave_assignment_mode = data.mongodbatlas_org_maintenance_settings.example.effective_wave_assignment_mode
  }
}
```

## Argument Reference

* `org_id` - (Required) Unique 24-hexadecimal digit string that identifies the Atlas organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `wave_assignment_mode` - Configured wave assignment mode for the organization. Accepted values are `MANUAL` and `ENV_TAG_MAPPING`. Defaults to `MANUAL` when unset.
* `effective_wave_assignment_mode` - Wave assignment mode Atlas actually uses for scheduling. Its value can differ from `wave_assignment_mode` in some cases. For more details, see the [`mongodbatlas_maintenance_window` data source](../data-sources/maintenance_window.md).

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-maintenance-windows)