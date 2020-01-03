---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: teams"
sidebar_current: "docs-mongodbatlas-resource-teams"
description: |-
    Provides a Team resource.
---

# mongodbatalas_teams

`mongodbatalas_teams` provides a Team resource. The resource lets you create, edit and delete Teams. Also, Teams can be assigned to multiple projects, and team members’ access to the project is determined by the team’s project role.

> **IMPORTANT:** MongoDB Atlas Team limits: max 250 teams in an organization and max 100 teams per project.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

MongoDB Atlas Team limits: max 250 teams in an organization and max 100 teams per project.

## Example Usage

```hcl
resource "mongodbatalas_teams" "test" {
    org_id    = "<Your Organization ID>"
    name      = "myNewTeam"
    usernames = ["user1", "user2", "user3"]
}
```

## Argument Reference

* `org_id` - (Required) The unique identifier for the organization you want to associate the team with.
* `name` - (Required) 	The name of the team you want to create.
* `usernames` - (Optional) Usernames to add to the new team.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `id` -	The Terraform's unique identifier used internally for state management.
* `team_id` - The unique identifier for the team.

## Import

Clusters can be imported using the organization ID and team id, in the format `ORGID-TEAMID`, e.g.

```
$ terraform import mongodbatalas_teams.my_team 1112222b3bf99403840e8934-1112222b3bf99403840e8935
```

See detailed information for arguments and attributes: [MongoDB API Teams](https://docs.atlas.mongodb.com/reference/api/teams-create-one/)
