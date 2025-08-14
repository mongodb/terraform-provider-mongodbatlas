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
- `mongodbatlas_cloud_user_team_assignment` manages the user’s organization membership (pending or active) and exposes both `username` and `user_id`. It supports import using either `ORG_ID/TEAM_ID/USERNAME` or `ORG_ID/TEAM_ID/USER_ID`.

---  
  
## From `mongodbatlas_team.usernames` to `mongodbatlas_cloud_user_team_assignment`
  
### Original configuration  
  
```terraform
locals {
  usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
}

resource "mongodbatlas_team" "this" {  
  org_id    = var.org_id  
  name      = "this"
  usernames = local.usernames
} 
```  
  
---  
  
### Step 1: User `mongodb_atlas_team` data source to retrieve user IDs
  
We first need to retrieve each user's `user_id` from the Atlas API via a data source.  
This is **required** if you already have a deployed team and want to migrate without recreating resources.  
  
```terraform  
# Use data source to get team members (with user_id)  
locals {
    usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
    team_assignments = {      
    for user in data.mongodbatlas_team.this.users :      
        user.id => {      
            org_id   = var.org_id
            team_id  = mongodbatlas_team.this.team_id
            user_id  = user.id  # Look up user_id here
        }
    }
}

resource "mongodbatlas_team" "this" {  
    org_id = var.org_id  
    name   = var.team_name
} 

data "mongodbatlas_team" "this" {  
    org_id  = var.org_id  
    team_id = mongodbatlas_team.this.team_id  
} 
```

---

### Step 2: Add `mongodbatlas_cloud_user_team_assignment`  and use import blocks

Handling migration in modules

- Terraform import blocks cannot live inside modules; they must be defined in the root module. See ([Terraform issue](https://github.com/hashicorp/terraform/issues/33474)).
- Module maintainers cannot ship import steps. Each module user must add root-level import blocks for every instance to import.

```terraform  
# Use data source to get team members (with user_id)  
locals {
    usernames = ["user1@email.com", "user2@email.com", "user3@email.com"]
    team_assignments = {
    for user in data.mongodbatlas_team.this.users :
        user.id => {
            org_id   = var.org_id
            team_id  = mongodbatlas_team.this.team_id
            user_id  = user.id  # Look up user_id here
        }
    }
}

resource "mongodbatlas_team" "this" {
    org_id = var.org_id
    name   = var.team_name
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
  
## Step 3: Run migration

Run `terraform plan` (you should see **import** operations), then `terraform apply`.
  
---  
  
## Step 4: Remove deprecated `usernames` from `mongodbatlas_team`  
  
Once the new resources are in place:  
  
```terraform  
resource "mongodbatlas_team" "this" {  
  org_id = var.org_id  
  name   = "this"  
  # usernames = local.usernames  # Remove this line
}  
```  
  
Run `terraform plan`. There should be **no changes**.  
  
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

---

## Data source migration  
  
If you previously used the `usernames` attribute in the `data.mongodbatlas_team` data source:  
  
**Original:**  

```terraform  
output "team_usernames" {  
  description = "Usernames in the MongoDB Atlas team"  
  value       = data.mongodbatlas_team.test.usernames  
}  
```  
  
**Replace with:**  

```terraform  
output "team_usernames" { 
  description = "Usernames in the MongoDB Atlas team"  
  value = [for u in data.mongodbatlas_team.team_1.users : u.username]  
}  
```  
  
---  
  
## Notes and tips  

- **Import format** for `mongodbatlas_cloud_user_team_assignment`:

```
  ORG_ID/TEAM_ID/USERNAME
  ORG_ID/TEAM_ID/USER_ID
```

- If using modules, remember to put import blocks in the root module.
- After successful migration, ensure **no references to** `mongodbatlas_team.usernames` remain.
