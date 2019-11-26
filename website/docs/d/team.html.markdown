---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: team"
sidebar_current: "docs-mongodbatlas-datasource-team"
description: |-
    Describes a Team.
---

# mongodbatlas_team

`mongodbatlas_team` describes a Team. The resource requires your Organization ID and Team ID.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_team" "test" {
    org_id    = "<Your Organization ID>"
    name      = "myNewTeam"
    usernames = ["user1", "user2", "user3"]
}

data "mongodbatlas_team" "test" {
    org_id    = mongodbatlas_team.test.org_id
    team_id   = mongodbatlas_team.test.team_id
}
```

## Argument Reference

* `org_id` - (Required) The unique identifier for the organization you want to associate the team with.
* `team_id` - (Required) The unique identifier for the team.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `id` -	The Terraform's unique identifier used internally for state management.
* `name` - (Required) 	The name of the team you want to create.
* `users` - (Optional) Usernames to add to the new team. See [User](#user).


### User

Indicates a user assigned to a Team

* `email_address` - The email address associated to the user.
* `first_name` - The first name of the user.
* `id` -  The unique identifier for the team.
* `last_name` -  The last name of the user.
* `roles` -  Each object in the roles array represents the Atlas organization role the user has for the associated orgId or groupId. See [Role](#role).
* `teams_id` -  Array of string IDs for each team the user is a member of.
* `username` - Username associated to the user.

### Role

Represents the Atlas organization role the user has for the associated orgId or groupId.

* `org_id` - ID of the organization where the user has the assigned roles.roleName organization role.
* `group_id` - ID of the project where the user has the assigned roles.roleName project role.
* `role_name` - The organization role assigned to the user for the specified roles.orgId or roles.groupId.

See detailed information for arguments and attributes: [MongoDB API Teams](https://docs.atlas.mongodb.com/reference/api/teams-create-one/)
