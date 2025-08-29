---
page_title: "Migration Guide: Project Teams Attribute to Team Project Assignment Resource"
---

  
# Migration Guide: Project Teams Attribute to Team Project Assignment Resource
  
**Objective:** Migrate from the deprecated `teams` attribute on the `mongodbatlas_project` resource to the new `mongodbatlas_team_project_assignment` resource.  
  
---  
  
## Before you begin  
  
- **Backup** your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state).  
- Ensure you are using the MongoDB Atlas Terraform Provider `2.0.0` or later (version that includes `mongodbatlas_team_project_assignment` resource).
  
---  
  
## Why should I migrate?  
  
- **Future compatibility:** The `teams` attribute inside `mongodbatlas_project` is deprecated and will be removed in a future provider release.  
- **Separation of concerns:** Manage projects and team-to-project role assignments independently.  
- **Clearer diffs:** Role or team modifications won't require re‑applying the entire project resource.  
  
---  
  
## What's changing?  
  
- Historically, `mongodbatlas_project` accepted an inline `teams` block to assign one or more teams to a project with specific roles.  
- Now, each project-team role mapping must be managed with `mongodbatlas_team_project_assignment`.

---

## From `mongodbatlas_project.teams` to `mongodbatlas_team_project_assignment`

### Original configuration
  
```hcl  
locals {  
  project_id = <PROJECT_ID>  
  team_map = { # team_id => set(role_names)
    <TEAM_ID_1>  = ["GROUP_OWNER"]
    <TEAM_ID_2>  = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_WRITE"]
  }
}

resource "mongodbatlas_project" "this" {
  name             = "<PROJECT_NAME>"
  org_id           = "<ORG_ID>"
  project_owner_id = "<OWNER_ID>"

  dynamic "teams" {
    for_each = local.team_map
    content {  
      team_id    = teams.key  
      role_names = teams.value  
    }  
  }  
}  
```  

---  
  
### Step 1: Ignore `teams` and remove from configuration

-> **Note:** The `teams` attribute is a `SetNestedBlock` and cannot be marked `Optional`/`Computed` for a smooth migration. For now, `ignore_changes` is required during Step 1. Support for removing `teams` entirely will come in a future Atlas Provider release.

Replace the `mongodbatlas_project.teams` block with:  
  
```hcl  
resource "mongodbatlas_project" "this" {  
  name             = "<PROJECT_NAME>"
  org_id           = "<ORG_ID>"
  project_owner_id = "<OWNER_ID>"  
  
  lifecycle {  
    ignore_changes = ["teams"]  
  }  
}  
```  
  
Then run:  
  
```shell  
terraform plan  
terraform apply  
```  
  
This removes the `teams` block from the config but keeps the assignments in Atlas unchanged until we explicitly manage them in new resources.  
  
---  
  
### Step 2: Add the new `mongodbatlas_team_project_assignment` resources  
  
```hcl  
resource "mongodbatlas_project" "this" {  
  name             = "<PROJECT_NAME>"
  org_id           = "<ORG_ID>"
  project_owner_id = "<OWNER_ID>"  
  
  lifecycle {  
    ignore_changes = ["teams"]  
  }  
}

resource "mongodbatlas_team_project_assignment" "this" {  
  for_each = local.team_map  
  
  project_id = mongodbatlas_project.this.id  
  team_id    = each.key  
  role_names = each.value  
}  
 
import {  
  for_each = local.team_map

  to       = mongodbatlas_team_project_assignment.this[each.key]
  id       = "${mongodbatlas_project.this.id}/${each.key}"
}  
```
  
Run `terraform plan` (you should see **import** operations), then `terraform apply`. 
  
---  
  
### Step 4: Verify and clean up  
  
- After successful import and apply, `terraform plan` should show **no changes**.  
- Keep the `ignore_changes = ["teams"]` lifecycle rule until the provider releases a version without the `teams` argument in `mongodbatlas_project`.  
  
---

## Examples

For complete, working configurations that demonstrate the migration process, see the examples in the provider repository: [migrate_team_project_assignment](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.0/examples/migrate_team_project_assignment).

The examples include:
- **v1**: Original configuration using deprecated `teams` attribute in `mongodbatlas_project` resource.
- **v2**: Final configuration using `mongodbatlas_team_project_assignment` resource for team-to-project assignments.
  
## Notes and tips  
  
- **Import format** for `mongodbatlas_team_project_assignment`:  
```  
PROJECT_ID/TEAM_ID  
```  
- **Modules:** Terraform import blocks cannot live inside modules ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)). 
- If you manage team assignments in modules, import each at the root level using the correct resource address (e.g. `module.<name>.mongodbatlas_team_project_assignment.<name>`).  
- You can use `terraform plan` to confirm imports before applying.  
  
---  
  
## FAQ  

**Q: Do I need to delete the old `teams` from state?**
A: No — using `ignore_changes` ensures they remain in Atlas until the provider removes the field. Then you can drop the lifecycle rule.  
  
---  
  
## Further resources  
- [`mongodbatlas_team_project_assignment` docs](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/team_project_assignment)
