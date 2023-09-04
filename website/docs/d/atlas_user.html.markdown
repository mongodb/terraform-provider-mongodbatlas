---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: atlas_user"
sidebar_current: "docs-mongodbatlas-datasource-atlas-user"
description: |-
    Provides a Atlas User Datasource.
---

# Data Source: mongodbatlas_atlas_user

`mongodbatlas_atlas_user` Provides a MongoDB Atlas User.

-> **NOTE:** If you are the owner of a MongoDB Atlas organization or project, you can also retrieve the user profile for any user with membership in that organization or project.

## Example Usage

### Using user_id attribute to query
```terraform
data "mongodbatlas_atlas_user" "test" {
  user_id = "<USER_ID>"
}
```

### Using username attribute to query
```terraform
data "mongodbatlas_atlas_user" "test" {
  username = "<USERNAME>"
}
```

## Argument Reference

* `user_id` - (Optional) Unique 24-hexadecimal digit string that identifies this user.
* `username` - (Optional) Email address that belongs to the MongoDB Atlas user account. You can't modify this address after creating the user.

~> **IMPORTANT:** Either `user_id` or `username` must be configurated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `country` - Two alphabet characters that identifies MongoDB Atlas user's geographic location. This parameter uses the ISO 3166-1a2 code format.
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

  
For additional documentation, see [MongoDB Atlas API - Get User By ID](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/MongoDB-Cloud-Users/operation/getUser) and [MongoDB Atlas API - Get User By Username](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/MongoDB-Cloud-Users/operation/getUserByUsername) respectively.
