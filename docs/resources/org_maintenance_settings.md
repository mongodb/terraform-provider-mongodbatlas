---
subcategory: "Organizations"
---

# Resource: mongodbatlas_org_maintenance_settings

`mongodbatlas_org_maintenance_settings` provides a resource to manage organization-level maintenance wave settings for a MongoDB Atlas organization. Use this resource to control how [Atlas assigns projects to maintenance waves](https://www.mongodb.com/docs/atlas/tutorial/cluster-maintenance-window/), either explicitly by project (`MANUAL` mode) or automatically based on each project's environment tag (`ENV_TAG_MAPPING` mode).

-> **NOTE:** Only one `mongodbatlas_org_maintenance_settings` resource can be defined per organization.

## Example Usage

### Manual wave assignment

Set the organization to `MANUAL` mode so you can assign each project to a specific maintenance wave using the `wave_assignment` attribute on [`mongodbatlas_maintenance_window`](maintenance_window.md).

```terraform
resource "mongodbatlas_org_maintenance_settings" "example" {
  org_id               = var.org_id
  wave_assignment_mode = "MANUAL"
}
```

### Automatic wave assignment based on environment tags

Set the organization to `ENV_TAG_MAPPING` mode to have Atlas derive the maintenance wave from each project's environment tag.

```terraform
resource "mongodbatlas_org_maintenance_settings" "example" {
  org_id               = var.org_id
  wave_assignment_mode = "ENV_TAG_MAPPING"
}
```

### Further Examples

<!-- TODO(CLOUDP-414003): Replace with versioned link once Mar's examples land -->
- [Configure Organization Maintenance Settings](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/[version]/examples/mongodbatlas_org_maintenance_settings)

## Argument Reference

`mongodbatlas_org_maintenance_settings` supports the following arguments:

* `org_id` - (Required) Unique 24-hexadecimal digit string that identifies the Atlas organization. This attribute cannot be changed after the resource is created.
* `wave_assignment_mode` - (Optional) Controls how Atlas assigns projects to maintenance waves. Accepted values are `MANUAL` and `ENV_TAG_MAPPING`. Remove this attribute from your configuration and run `terraform apply` to reset the mode to `MANUAL`.

## Import

Organization maintenance settings can be imported using the organization ID, in the format `ORG_ID`, e.g.

```
$ terraform import mongodbatlas_org_maintenance_settings.example 5d09d6a59ccf6445652a444a
```

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/maintenance-windows/)
