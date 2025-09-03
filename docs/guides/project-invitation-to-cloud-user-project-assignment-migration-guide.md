---
page_title: "Migration Guide: Project Invitation to Cloud User Project Assignment"
---

# Migration Guide: Project Invitation to Cloud User Project Assignment

**Objective**: Migrate from the deprecated `mongodbatlas_project_invitation` resource and data source to the `mongodbatlas_cloud_user_project_assignment` resource.

## Before you begin

- Back up your Terraform state file.
- Ensure you are using MongoDB Atlas Terraform Provider version 2.0.0 or later that includes `mongodbatlas_cloud_user_project_assignment`.

## What’s changing?

- `mongodbatlas_project_invitation` only managed invitations and is deprecated. If the user accepted the invitation and is now a project member, the provider removed the invitation from Terraform state and you should remove it from your configuration as well. See the resource [documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project_invitation) for more details.
- `mongodbatlas_cloud_user_project_assignment` manages the user’s project membership (active members).
- Pending project invitations are not discoverable with the new APIs. The only migration path for existing PENDING invites is to re-create them using `mongodbatlas_cloud_user_project_assignment` with the same `username` and `roles`.
 - For details on the new resource, see the `mongodbatlas_cloud_user_project_assignment` resource documentation: https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_user_project_assignment

---

## Migrating PENDING invitation (resource exists in config)

Original configuration:

```terraform
locals {
  username = "user1@email.com"
  roles    = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

resource "mongodbatlas_project_invitation" "this" {
  project_id = var.project_id
  username   = local.username
  roles      = local.roles
}
```

### Step 1: Add the new resource alongside existing configuration

Add the new resource to re-create the pending invite via the new API:

```terraform
resource "mongodbatlas_cloud_user_project_assignment" "this" {
  project_id = var.project_id
  username   = local.username
  roles      = local.roles
}
```

Use the same `roles` as the original invitation to avoid drift.

### Step 2: Remove the deprecated resource from the configuration and state

#### Option A) [Recommended] Removed block

Remove the resource block and replace it with a `removed` block to cleanly remove the old resource from state:

```terraform
removed {
  from = mongodbatlas_project_invitation.this

  lifecycle {
    destroy = false
  }
}
```

#### Option B) Manual state removal

Remove the `mongodbatlas_project_invitation` resource from configuration and then remove it from the Terraform state using the command line (this does not affect the actual invitation in Atlas):

```bash
terraform state rm mongodbatlas_project_invitation.this
```

### Step 3: Apply the changes

Run `terraform apply` to create the assignment with the new resource. Afterwards, run `terraform plan` and ensure no further changes are pending.

---

## Examples

For complete, working configurations that demonstrate the migration process, see the examples in the provider repository: [migrate_project_invitation_to_cloud_user_project_assignment](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_project_invitation_to_cloud_user_project_assignment).

The examples include:
- **v1**: Original configuration using deprecated `mongodbatlas_project_invitation`
- **v2**: Migration phase with re-creation using new resource and clean state removal
- **v3**: Final clean configuration using only `mongodbatlas_cloud_user_project_assignment`

These examples provide practical validation of the migration steps and demonstrate the re-creation approach for pending invitations.

---

## Notes and tips

- After successful migration, ensure no references to `mongodbatlas_project_invitation` remain in configuration or state.
- Pending invitations are not discoverable by the new APIs and resources; there is no data source replacement for reading pending invites. Re-create them using the new resource as shown above.
- For additional details on how accepted invitations are handled, see the `mongodbatlas_project_invitation` resource [documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project_invitation).
