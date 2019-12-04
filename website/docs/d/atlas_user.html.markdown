---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: users"
sidebar_current: "docs-mongodbatlas-datasource-user"
description: |-
    Describes a Atlas User.
---

# mongodbatlas_users

`mongodbatlas_users` Describe an Atlas User resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_user" "test" {
  username     = "john.doe@example.com"
  password     = "myPassword1@"
  emailAddress = "john.doe@example.com"
  mobileNumber = "2125550198"
  firstName    = "John"
  lastName     = "Doe"
  
  roles {
    orgId    = "8dbbe4570bd55b23f25444db"
    roleName = "ORG_MEMBER"

  }

  roles {
    groupId  = "2ddoa1233ef88z75f64578ff"
    roleName = "GROUP_READ_ONLY"

  }

  country = "US"
}

data "mongodbatlas_user" "test" {
	user_id = mongodbatlas_user.test.user_id
}


```

## Argument Reference

* `user_id` - Unique identifier for the Atlas user.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `id` -	The Terraform's unique identifier used internally for state management.
* `username` - (Required) The Atlas username. Must use e-mail address formatting. You cannot modify the username once set.
* `password` - (Required) Password. This field is NOT included in the entity returned from the server. it can only be sent in the entity body when creating a new user. You cannot update the password via API once set. The user must log into Atlas and update their password from the Atlas UI.
* `country` - (Required) The ISO 3166-1 alpha 2 country code of the Atlas user’s country of residence.
* `email_address` - (Required) The Atlas user’s email address.
* `mobile_number` - (Required) The Atlas user’s mobile or cell phone number.
* `first_name` - (Required) The Atlas user’s first name.
* `last_name` - (Required) The Atlas user’s last name.
* `roles` (Required)  - Each object in the array represents either an Atlas organization or project the Atlas user is assigned to and the Atlas role has for the associated organization or project. You can specify either roles.orgId or roles.groupId per object. See [Role](#role).
* `team_ids` - Array of string IDs for each team the user is a member of.



### Role

Represents the Atlas organization role the user has for the associated orgId or groupId.

* `org_id` - (Optional) The unique identifier of the organization in which the user has the specified roles.roleName.
* `group_id` - (Optional) The unique identifier of the project in which the user has the specified roles.roleName.
* `role_name` - (Required) The name of the role.

        When associated to roles.orgId, the valid roles and their mappings are:

        ORG_OWNER - Organization Owner
        ORG_GROUP_CREATOR - Organization Project Creator
        ORG_BILLING_ADMIN - Organization Billing Admin
        ORG_READ_ONLY - Organization Read Only
        ORG_MEMBER - Organization Member

        When associated to roles.groupId, the valid roles and their mappings are:

        GROUP_OWNER - Project Owner
        GROUP_CLUSTER_MANAGER - Project Cluster Manager
        GROUP_READ_ONLY - Project Read Only
        GROUP_DATA_ACCESS_ADMIN - Project Data Access Admin
        GROUP_DATA_ACCESS_READ_WRITE - Project Data Access Read/Write
        GROUP_DATA_ACCESS_READ_ONLY - Project Data Access Read Only


See detailed information for arguments and attributes: [MongoDB API Atlas Users](https://docs.atlas.mongodb.com/reference/api/user/)