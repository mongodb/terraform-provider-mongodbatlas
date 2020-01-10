---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: database_user"
sidebar_current: "docs-mongodbatlas-datasource-database-user"
description: |-
    Describes a Database User.
---

# mongodbatlas_database_user

`mongodbatlas_database_user` describe a Database User. This represents a database user which will be applied to all clusters within the project.

Each user has a set of roles that provide access to the project’s databases. User's roles apply to all the clusters in the project: if two clusters have a `products` database and a user has a role granting `read` access on the products database, the user has that access on both clusters.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_database_user" "test" {
	username      = "test-acc-username"
	password      = "test-acc-password"
	project_id      = "<PROJECT-ID>"
	database_name = "admin"
	
	roles {
		role_name     = "readWrite"
		database_name = "admin"
	}

    roles {
		role_name     = "atlasAdmin"
		database_name = "admin"
	}

	labels {
		key   = "key 1"
		value = "value 1"
	}
	labels {
		key   = "key 2"
		value = "value 2"
	}
}

data "mongodbatlas_database_user" "test" {
	project_id = mongodbatlas_database_user.test.project_id
	username = mongodbatlas_database_user.test.username
}

```

## Argument Reference

* `username` - (Required) Username for authenticating to MongoDB.
* `project_id` - (Required) The unique ID for the project to create the database user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The database user's name.
* `roles` - List of user’s roles and the databases / collections on which the roles apply. A role allows the user to perform particular actions on the specified database. A role on the admin database can include privileges that apply to the other databases as well. See [Roles](#roles) below for more details.
* `database_name` - The user’s authentication database. A user must provide both a username and authentication database to log into MongoDB. In Atlas deployments of MongoDB, the authentication database is always the admin database.

### Roles

Block mapping a user's role to a database / collection. A role allows the user to perform particular actions on the specified database. A role on the admin database can include privileges that apply to the other databases as well.

-> **NOTE:** The available privilege actions for custom MongoDB roles support a subset of MongoDB commands. See Unsupported Commands in M10+ Clusters for more information.

~> **IMPORTANT:** If a user is assigned a custom MongoDB role, they cannot be assigned any other roles.

* `name` - Name of the role to grant.
* `database_name` -  Database on which the user has the specified role. A role on the `admin` database can include privileges that apply to the other databases.
* `collection_name` - Collection for which the role applies. You can specify a collection for the `read` and `readWrite` roles. If you do not specify a collection for `read` and `readWrite`, the role applies to all collections in the database (excluding some collections in the `system`. database).

### Labels
Containing key-value pairs that tag and categorize the database user. Each key and value has a maximum length of 255 characters.

* `key` - The key that you want to write.
* `value` - The value that you want to write.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/database-users-get-single-user/) Documentation for more information.