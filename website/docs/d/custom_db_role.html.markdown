---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: custom_db_role"
sidebar_current: "docs-mongodbatlas-datasource-custom-db-role"
description: |-
    Describes a Custom DB Role.
---

# Data Source: mongodbatlas_custom_db_role

`mongodbatlas_custom_db_role` describe a Custom DB Role. This represents a custom db role.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_custom_db_role" "test_role" {
  project_id = "<PROJECT-ID>"
  role_name  = "myCustomRole"

  actions {
    action = "UPDATE"
    resources {
      collection_name = ""
      database_name   = "anyDatabase"
    }
  }
  actions {
    action = "INSERT"
    resources {
      collection_name = ""
      database_name   = "anyDatabase"
    }
  }
}

data "mongodbatlas_custom_db_role" "test" {
  project_id = mongodbatlas_custom_db_role.test_role.project_id
  role_name  = mongodbatlas_custom_db_role.test_role.role_name
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `role_name` - (Required) Name of the custom role. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier used for terraform for internal manages and can be used to import.

### Actions
Each object in the actions array represents an individual privilege action granted by the role. It is an required field.

* `action` - (Required) Name of the privilege action. For a complete list of actions available in the Atlas API, see Custom Role Actions.

* `resources` - (Required) Contains information on where the action is granted. Each object in the array either indicates a database and collection on which the action is granted, or indicates that the action is granted on the cluster resource.

* `resources.#.collection_name` - (Optional) Collection on which the action is granted. If this value is an empty string, the action is granted on all collections within the database specified in the actions.resources.db field.

* `resources.#.database_name`	Database on which the action is granted.

* `resources.#.cluster`	(Optional) Set to true to indicate that the action is granted on the cluster resource.

### Inherited Roles
Each object in the inheritedRoles array represents a key-value pair indicating the inherited role and the database on which the role is granted. It is an optional field.

* `database_name` (Required) Database on which the inherited role is granted.
* `role_name`	(Required) Name of the inherited role. This can either be another custom role or a built-in role.


See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/custom-roles-get-single-role/) Documentation for more information.