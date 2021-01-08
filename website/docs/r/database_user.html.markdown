---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: database_user"
sidebar_current: "docs-mongodbatlas-resource-database-user"
description: |-
    Provides a Database User resource.
---

# mongodbatlas_database_user

`mongodbatlas_database_user` provides a Database User resource. This represents a database user which will be applied to all clusters within the project.

Each user has a set of roles that provide access to the project’s databases. User's roles apply to all the clusters in the project: if two clusters have a `products` database and a user has a role granting `read` access on the products database, the user has that access on both clusters.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** All arguments including the password will be stored in the raw state as plain-text. [Read more about sensitive data in state.](https://www.terraform.io/docs/state/sensitive-data.html)

## Example Usages

```hcl
resource "mongodbatlas_database_user" "test" {
  username           = "test-acc-username"
  password           = "test-acc-password"
  project_id         = "<PROJECT-ID>"
  auth_database_name = "admin"

  roles {
    role_name     = "readWrite"
    database_name = "dbforApp"
  }

  roles {
    role_name     = "readAnyDatabase"
    database_name = "admin"
  }

  labels {
    key   = "My Key"
    value = "My Value"
  }

  scopes {
    name   = "My cluster name"
    type = "CLUSTER"
  }

  scopes {
    name   = "My second cluster name"
    type = "CLUSTER"
  }
}
```


```hcl
resource "mongodbatlas_database_user" "test" {
  username           = "test-acc-username"
  x509_type          = "MANAGED"
  project_id         = "<PROJECT-ID>"
  auth_database_name = "$external"

  roles {
    role_name     = "readAnyDatabase"
    database_name = "admin"
  }

  labels {
    key   = "%s"
    value = "%s"
  }

  scopes {
    name   = "My cluster name"
    type = "CLUSTER"
  }
}
```

## Argument Reference

* `auth_database_name` - (Required) Database against which Atlas authenticates the user. A user must provide both a username and authentication database to log into MongoDB.
Accepted values include:
  * `admin` if `x509_type` and `aws_iam_type` are omitted or NONE.
  * `$external` if `x509_type` is MANAGED or CUSTOMER or `aws_iam_type` is USER or ROLE.
* `project_id` - (Required) The unique ID for the project to create the database user.
* `roles` - (Required) 	List of user’s roles and the databases / collections on which the roles apply. A role allows the user to perform particular actions on the specified database. A role on the admin database can include privileges that apply to the other databases as well. See [Roles](#roles) below for more details.
* `username` - (Required) Username for authenticating to MongoDB.
* `password` - (Required) User's initial password. A value is required to create the database user, however the argument but may be removed from your Terraform configuration after user creation without impacting the user, password or Terraform management. IMPORTANT --- Passwords may show up in Terraform related logs and it will be stored in the Terraform state file as plain-text. Password can be changed after creation using your preferred method, e.g. via the MongoDB Atlas UI, to ensure security.  If you do change management of the password to outside of Terraform be sure to remove the argument from the Terraform configuration so it is not inadvertently updated to the original password.

* `x509_type` - (Optional) X.509 method by which the provided username is authenticated. If no value is given, Atlas uses the default value of NONE. The accepted types are:
  * `NONE` -	The user does not use X.509 authentication.
  * `MANAGED` - The user is being created for use with Atlas-managed X.509.Externally authenticated users can only be created on the `$external` database.
  * `CUSTOMER` -  The user is being created for use with Self-Managed X.509. Users created with this x509Type require a Common Name (CN) in the username field. Externally authenticated users can only be created on the `$external` database.

* `aws_iam_type` - (Optional) If this value is set, the new database user authenticates with AWS IAM credentials. If no value is given, Atlas uses the default value of NONE. The accepted types are:
  * `NONE` -	The user does not use AWS IAM credentials.
  * `USER` - New database user has AWS IAM user credentials.
  * `ROLE` -  New database user has credentials associated with an AWS IAM role.

* `ldap_auth_type` - (Optional) Method by which the provided username is authenticated. If no value is given, Atlas uses the default value of NONE.
  * `NONE` -	Atlas authenticates this user through SCRAM-SHA, not LDAP.
  * `USER` - LDAP server authenticates this user through the user's LDAP user. `username` must also be a fully qualified distinguished name, as defined in RFC-2253.
  * `GROUP` -  LDAP server authenticates this user using their LDAP user and authorizes this user using their LDAP group.

### Roles

Block mapping a user's role to a database / collection. A role allows the user to perform particular actions on the specified database. A role on the admin database can include privileges that apply to the other databases as well.

-> **NOTE:** The available privilege actions for custom MongoDB roles support a subset of MongoDB commands. See Unsupported Commands in M10+ Clusters for more information.

~> **IMPORTANT:** If a user is assigned a custom MongoDB role, they cannot be assigned any other roles.

* `role_name` - (Required) Name of the role to grant. See [Create a Database User](https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/) `roles.roleName` for valid values and restrictions.
* `database_name` - (Required) Database on which the user has the specified role. A role on the `admin` database can include privileges that apply to the other databases.
* `collection_name` - (Optional) Collection for which the role applies. You can specify a collection for the `read` and `readWrite` roles. If you do not specify a collection for `read` and `readWrite`, the role applies to all collections in the database (excluding some collections in the `system`. database).

### Labels
Containing key-value pairs that tag and categorize the database user. Each key and value has a maximum length of 255 characters.

* `key` - The key that you want to write.
* `value` - The value that you want to write.

### Scopes
Array of clusters and Atlas Data Lakes that this user has access to. If omitted, Atlas grants the user access to all the clusters and Atlas Data Lakes in the project by default.

* `name` - (Required) Name of the cluster or Atlas Data Lake that the user has access to.
* `type` - (Required) Type of resource that the user has access to. Valid values are: `CLUSTER` and `DATA_LAKE`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The database user's name.

## Import

Database users can be imported using project ID and username, in the format `project_id`-`username`-`auth_database_name`, e.g.

```
$ terraform import mongodbatlas_database_user.my_user 1112222b3bf99403840e8934-my_user-admin
```

~> **NOTE:** Terraform will want to change the password after importing the user if a `password` argument is specified.
