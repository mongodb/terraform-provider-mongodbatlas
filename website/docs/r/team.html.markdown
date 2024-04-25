---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: team"
sidebar_current: "docs-mongodbatlas-resource-team"
description: |-
    Provides a Team resource.
---

# Resource: mongodbatlas_team

`mongodbatlas_team` provides a Team resource. The resource lets you create, edit and delete Teams. Also, Teams can be assigned to multiple projects, and team members’ access to the project is determined by the team’s project role.

> **IMPORTANT:** MongoDB Atlas Team are limited to a maximum of 250 teams in an organization and 100 teams per project.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_team" "test" {
  org_id     = "<ORGANIZATION-ID>"
  name       = "myNewTeam"
  usernames  = ["user1@email.com", "user2@email.com", "user3@email.com"]
}
```

## Argument Reference

* `org_id` - (Required) The unique identifier for the organization you want to associate the team with.
* `name` - (Required) The name of the team you want to create.
* `usernames` - (Required) The Atlas usernames (email address). You can only add Atlas users who are part of the organization. Users who have not accepted an invitation to join the organization cannot be added as team members. There is a maximum of 250 Atlas users per team. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` -	The Terraform's unique identifier used internally for state management.
* `team_id` - The unique identifier for the team.

## Import

Teams can be imported using the organization ID and team id, in the format ORGID-TEAMID, e.g.

```
$ terraform import mongodbatlas_team.my_team 1112222b3bf99403840e8934-1112222b3bf99403840e8935
```

See detailed information for arguments and attributes: [MongoDB API Teams](https://docs.atlas.mongodb.com/reference/api/teams-create-one/)
