---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: database_user"
sidebar_current: "docs-mongodbatlas-resource-database_user"
description: |-
    Provides a Database User resource.
---

# mongodb_atlas_database_user

`mongodb_atlas_database_user` provides a Database User resource. This represents a database user which will be applied to all clusters within the project.

Each user has a set of roles that provide access to the project’s databases. User's roles apply to all the clusters in the project: if two clusters have a `products` database and a user has a role granting `read` access on the products database, the user has that access on both clusters.

~> **NOTE:** All arguments including the password will be stored in the raw state as plain-text. [Read more about sensitive data in state](https://www.terraform.io/docs/state/sensitive-data.html)

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_database_user" "test" {
	username      = "test-acc-username"
	password      = "test-acc-password"
	project_id    = "<PROJECT-ID>"
	database_name = "admin"
	
	roles {
		role_name     = "readWrite"
		database_name = "admin"
	}

    roles {
		role_name     = "%s"
		database_name = "admin"
	}
}
```

## Argument Reference

* `database_name` - (Required) The user’s authentication database. A user must provide both a username and authentication database to log into MongoDB. In Atlas deployments of MongoDB, the authentication database is always the admin database.
* `project_id` - (Required) The unique ID for the project to create the database user.
* `password` - (Optional) User's initial password. This is required to create the user but may be removed after.

~> **NOTE:** Password may show up in logs, and it will be stored in the state file as plain-text. Password can be changed in the web interface to increase security.



* `roles` - (Required) 	List of user’s roles and the databases / collections on which the roles apply. A role allows the user to perform particular actions on the specified database. A role on the admin database can include privileges that apply to the other databases as well. See [Roles](#roles) below for more details.
* `username` - (Required) Username for authenticating to MongoDB.

### Roles

Block mapping a user's role to a database / collection. A role allows the user to perform particular actions on the specified database. A role on the admin database can include privileges that apply to the other databases as well.

**NOTE:** The available privilege actions for custom MongoDB roles support a subset of MongoDB commands. See Unsupported Commands in M10+ Clusters for more information.

**IMPORTANT** If a user is assigned a custom MongoDB role, they cannot be assigned any other roles.

* `name` - (Required) Name of the role to grant. See [Create a Database User](https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/) `roles.roleName` for valid values and restrictions.
* database_name - (Required) Database on which the user has the specified role. A role on the `admin` database can include privileges that apply to the other databases.
* `collection_name` - (Optional) Collection for which the role applies. You can specify a collection for the `read` and `readWrite` roles. If you do not specify a collection for `read` and `readWrite`, the role applies to all collections in the database (excluding some collections in the `system`. database).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The database user's name.

## Import

Database users can be imported using project ID and username, in the format `PROJECTID-USERNAME`, e.g.

```
$ terraform import mongodbatlas_database_user.my_user 1112222b3bf99403840e8934-my_user
```

~> **NOTE:** Terraform will want to change the password after importing the user if a `password` argument is specified.