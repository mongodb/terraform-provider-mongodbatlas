---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: atlas_users"
sidebar_current: "docs-mongodbatlas-datasource-mongodbatlas-atlas_users"
description: |-
    Provides a Atlas Users Datasource.
---

# Data Source: atlas_users

`atlas_users` provides Atlas Users associated with a specified Organization, Project, or Team.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage
### Using org_id attribute to query Organization Atlas Users

```terraform
data "mongodbatlas_atlas_users" "test" {
  org_id = "<ORG_ID>"
}
```

### Using project_id attribute to query Project Atlas Users

```terraform
data "mongodbatlas_atlas_users" "test" {
  project_id = "<PROJECT_ID>"
}
```

### Using team_id and org_id attribute to query Team Atlas Users

```terraform
data "mongodbatlas_atlas_users" "test" {
  team_id = "<TEAM_ID>"
  org_id = "<ORG_ID>"
}
```

## Argument Reference

* `org_id` - (Optional) Unique 24-hexadecimal digit string that identifies the organization whose users you want to return. Also needed when `team_id` attributes is defined.
* `project_id` - (Optional) Unique 24-hexadecimal digit string that identifies the project whose users you want to return. 
* `team_id` - (Optional) Unique 24-hexadecimal digit string that identifies the team whose users you want to return.

* `page_num` - (Optional) Number of the page that displays the current set of the total objects that the response returns. Defaults to `1`.
* `items_per_page` - (Optional) Number of items that the response returns per page, up to a maximum of `500`. Defaults to `100`.

~> **IMPORTANT:** Either `org_id`, `project_id`, or `team_id` with `org_id` must be configurated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `results` - A list where each element contains a Atlas User.
* `total_count` - Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.

### Atlas User

* `user_id` - Unique 24-hexadecimal digit string that identifies this user.
* `username` - Email address that belongs to the MongoDB Atlas user account. You cannot modify this address after creating the user.
* `country` - Two alphabet characters that identifies MongoDB Cloud user's geographic location. This parameter uses the ISO 3166-1a2 code format.
* `created_at` - Date and time when the current account is created. This value is in the ISO 8601 timestamp format in UTC.
* `email_address` - Email address that belongs to the MongoDB Atlas user.
* `first_name` - First or given name that belongs to the MongoDB Atlas user.
* `last_auth` - Date and time when the current account last authenticated. This value is in the ISO 8601 timestamp format in UTC.
* `last_name` - Last name, family name, or surname that belongs to the MongoDB Atlas user.
* `mobile_number` - Mobile phone number that belongs to the MongoDB Atlas user.
* `team_ids` - List of unique 24-hexadecimal digit strings that identifies the teams to which this MongoDB Atlas user belongs.
* `links.#.href` - Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with https://cloud.mongodb.com/api/atlas.
* `links.#.rel` - Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with https://cloud.mongodb.com/api/atlas.
* `roles.#.group_id` - Unique 24-hexadecimal digit string that identifies the project to which this role belongs. You can set a value for this parameter or orgId but not both in the same request.
* `roles.#.org_id` - Unique 24-hexadecimal digit string that identifies the organization to which this role belongs. You can set a value for this parameter or groupId but not both in the same request.
* `roles.#.role_name` - Human-readable label that identifies the collection of privileges that MongoDB Atlas grants a specific API key, user, or team. These roles include organization- and project-level roles. The [MongoDB Documentation](https://www.mongodb.com/docs/atlas/reference/user-roles/#service-user-roles) describes the valid roles that can be assigned.

  
For additional documentation, see:
- For obtaining users of an Organization: [MongoDB Atlas API - List Organization Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/listOrganizationUsers) 
- For obtaining users of a Project: [MongoDB Atlas API - List Project Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Projects/operation/listProjectUsers)
- For obtaining users of a Team: [MongoDB Atlas API - List Team Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Teams/operation/listTeamUsers)
