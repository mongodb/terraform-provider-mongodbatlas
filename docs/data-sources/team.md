# Data Source: mongodbatlas_team

`mongodbatlas_team` describes a Team. The resource requires your Organization ID, Project ID and Team ID.

-> **NOTE:** Groups and projects are synonymous terms. You may find `group_id` in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_team" "test" {
  org_id     = "<ORGANIZATION-ID>"
  name       = "myNewTeam"
  usernames  = ["user1", "user2", "user3"]
}

data "mongodbatlas_team" "test" {
	org_id     = mongodbatlas_team.test.org_id
	team_id    = mongodbatlas_team.test.team_id
}

```

```terraform
resource "mongodbatlas_team" "test" {
  org_id     = "<ORGANIZATION-ID>"
  name       = "myNewTeam"
  usernames  = ["user1", "user2", "user3"]
}

data "mongodbatlas_team" "test2" {
	org_id     = mongodbatlas_team.test.org_id
	name       = mongodbatlas_team.test.name
}
```


## Argument Reference

* `org_id` - (Required) The unique identifier for the organization you want to associate the team with.
* `team_id` - (Optional) The unique identifier for the team.
* `name` - (Optional) The team name.

~> **IMPORTANT:** Either `team_id` or `name` must be configured.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Terraform's unique identifier used internally for state management.
* `team_id` -  The unique identifier for the team.
* `name` -  The name of the team you want to create.
* `usernames` - The users who are part of the organization.
* `users`- Returns a list of all pending and active MongoDB Cloud users associated with the specified organization.

### Users
* `id` - Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.
* `org_membership_status` - String enum that indicates whether the MongoDB Cloud user has a pending invitation to join the organization or are already active in the organization.
* `roles` - Organization and project-level roles assigned to one MongoDB Cloud user within one organization.
* `team_ids` - List of unique 24-hexadecimal digit strings that identifies the teams to which this MongoDB Cloud user belongs.
* `username` - Email address that represents the username of the MongoDB Cloud user.
* `country` - Two-character alphabetical string that identifies the MongoDB Cloud user's geographic location. This parameter uses the ISO 3166-1a2 code format.
* `invitation_created_at` - Date and time when MongoDB Cloud sent the invitation. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
* `invitation_expires_at` - Date and time when the invitation from MongoDB Cloud expires. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
* `inviter_username` - Username of the MongoDB Cloud user who sent the invitation to join the organization.
* `created_at` - Date and time when MongoDB Cloud created the current account. This value is in the ISO 8601 timestamp format in UTC.
* `first_name` - First or given name that belongs to the MongoDB Cloud user.
* `last_auth` - Date and time when the current account last authenticated. This value is in the ISO 8601 timestamp format in UTC.
* `last_name` - Last name, family name, or surname that belongs to the MongoDB Cloud user.
* `mobile_number` - Mobile phone number that belongs to the MongoDB Cloud user.


~> **NOTE:** - Users with pending invitations created using [`mongodbatlas_project_invitation`](../resources/project_invitation.md) resource or via the deprecated [Invite One MongoDB Cloud User to Join One Project](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createprojectinvitation) endpoint are excluded (or cannot be managed) with this resource. See  [MongoDB Atlas API]<link-to-resource-API> for details. 
To manage these users with this resource/data source, refer to our [migration guide]<link-to-migration-guide>.

See detailed information for arguments and attributes: [MongoDB API Teams](https://docs.atlas.mongodb.com/reference/api/teams-create-one/)
