---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: teams"
sidebar_current: "docs-mongodbatlas-datasource-teams"
description: |-
    Describes a Team.
---

# Data Source: mongodbatlas_teams

`mongodbatlas_teams` describes a Team. The resource requires your Organization ID, Project ID and Team ID.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_teams" "test" {
  org_id     = "<ORGANIZATION-ID>"
  name       = "myNewTeam"
  usernames  = ["user1", "user2", "user3"]
}

data "mongodbatlas_teams" "test" {
	org_id     = mongodbatlas_teams.test.org_id
	team_id    = mongodbatlas_teams.test.team_id
}

```

```terraform
resource "mongodbatlas_teams" "test" {
  org_id     = "<ORGANIZATION-ID>"
  name       = "myNewTeam"
  usernames  = ["user1", "user2", "user3"]
}

data "mongodbatlas_teams" "test2" {
	org_id     = mongodbatlas_teams.test.org_id
	name       = mongodbatlas_teams.test.name
}
```


## Argument Reference

* `org_id` - (Required) The unique identifier for the organization you want to associate the team with.
* `team_id` - (Optional) The unique identifier for the team.
* `name` - (Optional) The team name.

~> **IMPORTANT:** Either `team_id` or `name` must be configured.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` -	The Terraform's unique identifier used internally for state management.
* `team_id` -  The unique identifier for the team.
* `name` -  The name of the team you want to create.
* `usernames` - The users who are part of the organization.

See detailed information for arguments and attributes: [MongoDB API Teams](https://docs.atlas.mongodb.com/reference/api/teams-create-one/)
