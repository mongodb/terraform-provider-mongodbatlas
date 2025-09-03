---
page_title: "Migration Guide: Team Usernames Attribute to Cloud User Team Assignment"
---
  
# Migration Guide: Team Usernames Attribute to Cloud User Team Assignment
  
**Objective**: Migrate from the deprecated `usernames` attribute on the `mongodbatlas_team` resource to the new `mongodbatlas_cloud_user_team_assignment` resource.
  
---  
  
## Before you begin
  
- Create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state).
- Ensure you are using the MongoDB Atlas Terraform Provider `2.0.0` or later (version that includes `mongodbatlas_cloud_user_team_assignment` resource).

---  

## Why should I migrate?

- **Future Compatibility:** The `usernames` attribute on `mongodbatlas_team` is deprecated and may be removed in future provider versions. Migrating ensures your Terraform configuration remains functional.
- **Flexibility:** Manage teams and user assignments independently, without coupling membership changes to team creation or updates.  
- **Clarity:** Clear separation between the `mongodbatlas_team` resource (team definition) and `mongodbatlas_cloud_user_team_assignment` (membership management).  

---  
  
## What’s changing?
  
- `mongodbatlas_team` included a `usernames` argument that allowed assigning users to a team directly inside the resource. This argument is now deprecated.
- New attribute `users` in `mongodbatlas_team` data source can be used to retrieve information about all the users assigned to that team.
- `mongodbatlas_cloud_user_team_assignment` manages the user’s team membership (pending or active) and exposes both `username` and `user_id`. It supports import using either `ORG_ID/TEAM_ID/USERNAME` or `ORG_ID/TEAM_ID/USER_ID`.

---  
  
## From `mongodbatlas_team.usernames` to `mongodbatlas_cloud_user_team_assignment`
  
### Original configuration
  
```terraform
locals {
  usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
}

resource "mongodbatlas_team" "this" {  
  org_id    = var.org_id  
  name      = var.team_name
  usernames = local.usernames
} 
```  
  
---  
  
### Step 1: Use `mongodbatlas_team` data source to retrieve user IDs
  
We first need to retrieve each user's `user_id` via the new `users` attribute in `mongodbatlas_team` data source.
  
```terraform  
# Use data source to get team members (with user_id)  
locals {
    usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
    team_assignments = {      
    for user in data.mongodbatlas_team.this.users :      
        user.id => {      
            org_id   = var.org_id
            team_id  = mongodbatlas_team.this.team_id
            user_id  = user.id
        }
    }
}

resource "mongodbatlas_team" "this" {  
    org_id = var.org_id  
    name   = var.team_name
    usernames = local.usernames
} 

data "mongodbatlas_team" "this" {  
    org_id  = var.org_id  
    team_id = mongodbatlas_team.this.team_id  
} 
```

---  

### Step 2: Add `mongodbatlas_cloud_user_team_assignment` and use import blocks

```terraform   
locals {
    usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
    team_assignments = {
    for user in data.mongodbatlas_team.this.users :
        user.id => {
            org_id   = var.org_id
            team_id  = mongodbatlas_team.this.team_id
            user_id  = user.id
        }
    }
}

resource "mongodbatlas_team" "this" {
    org_id = var.org_id
    name   = var.team_name
    usernames = local.usernames
}

data "mongodbatlas_team" "this" {
    org_id  = var.org_id
    team_id = mongodbatlas_team.this.team_id
}
  
# New resource for each (user, team) assignment  
resource "mongodbatlas_cloud_user_team_assignment" "this" {           
    for_each = local.team_assignments

    org_id  = each.value.org_id   
    team_id = each.value.team_id     
    user_id = each.value.user_id  # Use user_id instead of username  
}  
  
# Import existing team-user relationships into the new resource  
import {  
    for_each = local.team_assignments

    to = mongodbatlas_cloud_user_team_assignment.this[each.key] 
    id = "${each.value.org_id}/${each.value.team_id}/${each.value.user_id}" 
} 
```
  
---  
## Step 3: Remove deprecated `usernames` from `mongodbatlas_team`  
  
Once the new resources are in place:  
  
```terraform  
resource "mongodbatlas_team" "this" {  
  org_id = var.org_id  
  name   = "this"  
  # usernames = local.usernames  # Remove this line
}  
```  

---
## Step 4: Run migration

Run `terraform plan` (you should see **import** operations), then `terraform apply`.
  
---  

  
## Step 5: Update any references to `mongodbatlas_team.usernames`  
  
Before:  
  
```terraform  
output "team_usernames" {  
  value = mongodbatlas_team.this.usernames  
}  
```  
  
After:  
  
```terraform  
output "team_usernames" {  
  value = [for u in data.mongodbatlas_team.this.users : u.username]  
}  
```

Run `terraform plan`. There should be **no changes**.

---  

## Data source migration

If you previously used the `usernames` attribute in the `data.mongodbatlas_team` data source:  
  
**Original:**  

```terraform  
output "team_usernames" {  
  description = "Usernames in the MongoDB Atlas team"  
  value       = data.mongodbatlas_team.this.usernames  
}  
```
  
**Replace with:**  

```terraform  
output "team_usernames" { 
  description = "Usernames in the MongoDB Atlas team"  
  value = [for u in data.mongodbatlas_team.this.users : u.username]  
}  
```

Run `terraform plan`. There should be **no changes**.
  
--- 

## Migration using Modules

If you are using modules to manage teams and user assignments to teams, migrating from `mongodbatlas_team` to the new pattern requires special attention. Because the old `mongodbatlas_team.usernames` attribute corresponds to `mongodbatlas_cloud_user_team_assignment`, you cannot simply move the resource block inside your module and expect Terraform to handle the migration automatically. This section demonstrates how to migrate from a module using the `mongodbatlas_team` resource to a module using both `mongodbatlas_team` and the new `mongodbatlas_cloud_user_team_assignment` resources.

**Key points for module users:**
- You must use `terraform import` to bring existing user-team assignments into the new resources, even when they are managed inside a module.
- The import command must match the resource address as used in your module (e.g., `module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>`).
- If you were using a list of usernames in your previous configuration, you also need to include the `mongodbatlas_team` data source and use the new `users` attribute to retrieve the corresponding user IDs, along with team ID, for the import to work correctly.

**Example import blocks for modules**
```terraform
import {
   to = module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>
   id = "<ORG_ID>/<TEAM_ID>/<USER_ID>"
}
import {
   to = module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>
   id = "<ORG_ID>/<TEAM_ID>/<USERNAME>"
}
```

**Example import commands for modules:**
```shell
terraform import 'module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>' <ORG_ID>/<TEAM_ID>/<USER_ID>
terraform import 'module.<module_name>.mongodbatlas_cloud_user_team_assignment.<name>' <ORG_ID>/<TEAM_ID>/<USERNAME>
```

### 1. Old Module Usage (Legacy)

```hcl
module "user_team_assignment" {  
  source     = "./old_module"  
  org_id     = var.org_id  
  team_name  = var.team_name  
  usernames  = var.usernames  
}
```

### 2. New Module Usage (Recommended)

```hcl
data "mongodbatlas_team" "this" {  
  org_id = var.org_id  
  name   = var.team_name
}

locals {  
  team_assigments = {
    for user in data.mongodbatlas_team.this.users :
    user.id => {
      org_id  = var.org_id
      team_id = data.mongodbatlas_team.this.team_id
      user_id = user.id
    }
  }  
}

module "user_team_assignment" {
  source     = "./new_module"
  org_id     = var.org_id
  team_name  = var.team_name
  team_assigments = local.team_assigments
}
```

### 3. Migration Steps

1. **Add the new module to your configuration:**
   - Add the new module block as shown above, using the same input variables as appropriate.
   - Also add the `data.mongodbatlas_team` data source and declare the `team_assignments` local variable to retrieve user IDs and team ID.

2. **Import the existing user-team assignments into the new resources:**

-  An `import block` (available in Terraform 1.5 and later) can be used to import the resource and iterate through a list of users, e.g.:
   ```terraform
  import {
    for_each = local.team_assigments

    to       = module.user_team_assignment.mongodbatlas_cloud_user_team_assignment.this[each.key]
    id       = "${var.org_id}/${data.mongodbatlas_team.this.team_id}/${each.value.user_id}"
  }
```

- Alternatively, use the correct resource addresses for your module and each of the user-team assignments:
```shell
  terraform import 'module.user_team_assignment.mongodbatlas_cloud_user_team_assignment.this' <ORG_ID>/<TEAM_ID>/<USER_ID>
```
   

3. **Remove the old module block from your configuration.**
4. **Run `terraform plan` to review the changes.**
   - Ensure that Terraform imports the user-team assignments and does not plan to create these.
   - Ensure that Terraform does not plan to destroy and recreate the `mongodbatlas_team` resource.
5. **Run `terraform apply` to apply the migration.**

For complete working examples, see:
- [Old module example](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/migrate_user_team_assignment/module/old_module/)
- [New module example](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/migrate_user_team_assignment/module/new_module/)


## Notes and tips

- **Import format** for `mongodbatlas_cloud_user_team_assignment`:

```
  ORG_ID/TEAM_ID/USERNAME
  ORG_ID/TEAM_ID/USER_ID
```

- **Importing inside modules:** Terraform import blocks cannot live inside modules. See ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)). Each module user must add root-level import blocks for every instance to import.

- After successful migration, ensure **no references to** `mongodbatlas_team.usernames` remain.

## FAQ
**Q: Can I assign the same user to multiple teams?**
A: Yes, simply create multiple `mongodbatlas_cloud_user_team_assignment` resources for each team.

**Q: Where can I find a working example?**
A: See [examples/mongodbatlas_cloud_user_team_assignment/main.tf](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_cloud_user_team_assignment/main.tf).

## Further Resources
- [Cloud User Team Assignment Resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_user_team_assignment)
