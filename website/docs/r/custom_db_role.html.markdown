# Resource: mongodbatlas_custom_db_role

`mongodbatlas_custom_db_role` provides a Custom DB Role resource. The customDBRoles resource lets you retrieve, create and modify the custom MongoDB roles in your cluster. Use custom MongoDB roles to specify custom sets of actions which cannot be described by the built-in Atlas database user privileges.

-> **IMPORTANT**  You define custom roles at the project level for all clusters in the project. The `mongodbatlas_custom_db_role` resource supports a subset of MongoDB privilege actions. For a complete list of [privilege actions](https://docs.mongodb.com/manual/reference/privilege-actions/) available for this resource, see [Custom Role actions](https://docs.atlas.mongodb.com/reference/api/custom-role-actions/). Custom roles must include actions that all project's clusters support, and that are compatible with each MongoDB version used by your project's clusters. For example, if your project has MongoDB 4.2 clusters, you can't create custom roles that use actions introduced in MongoDB 4.4.


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
  actions {
    action = "REMOVE"
    resources {
      collection_name = ""
      database_name   = "anyDatabase"
    }
  }
}
```

## Example Usage with inherited roles

```terraform
resource "mongodbatlas_custom_db_role" "inherited_role_one" {
  project_id = "<PROJECT-ID>"
  role_name  = "insertRole"

  actions {
    action = "INSERT"
    resources {
      collection_name = ""
      database_name   = "anyDatabase"
    }
  }
}

resource "mongodbatlas_custom_db_role" "inherited_role_two" {
  project_id = mongodbatlas_custom_db_role.inherited_role_one.project_id
  role_name  = "statusServerRole"

  actions {
    action = "SERVER_STATUS"
    resources {
      cluster = true
    }
  }
}

resource "mongodbatlas_custom_db_role" "test_role" {
  project_id = mongodbatlas_custom_db_role.inherited_role_one.project_id
  role_name  = "myCustomRole"

  actions {
    action = "UPDATE"
    resources {
      collection_name = ""
      database_name   = "anyDatabase"
    }
  }
  actions {
    action = "REMOVE"
    resources {
      collection_name = ""
      database_name   = "anyDatabase"
    }
  }

  inherited_roles {
    role_name     = mongodbatlas_custom_db_role.inherited_role_one.role_name
    database_name = "admin"
  }

  inherited_roles {
    role_name     = mongodbatlas_custom_db_role.inherited_role_two.role_name
    database_name = "admin"
  }
}

```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `role_name` - (Required) Name of the custom role.

	-> **IMPORTANT** The specified role name can only contain letters, digits, underscores, and dashes. Additionally, you cannot specify a role name which meets any of the following criteria:

	* Is a name already used by an existing custom role in the project
	* Is a name of any of the built-in roles
	* Is `atlasAdmin`
	* Starts with `xgen-`


### Actions
Each object in the actions array represents an individual privilege action granted by the role. It is an required field.

* `action` - (Required) Name of the privilege action. For a complete list of actions available in the Atlas API, see [Custom Role Actions](https://docs.atlas.mongodb.com/reference/api/custom-role-actions)
-> **Note**: The privilege actions available to the Custom Roles API resource represent a subset of the privilege actions available in the Atlas Custom Roles UI.

* `resources` - (Required) Contains information on where the action is granted. Each object in the array either indicates a database and collection on which the action is granted, or indicates that the action is granted on the cluster resource.

* `resources.#.collection_name` - (Optional) Collection on which the action is granted. If this value is an empty string, the action is granted on all collections within the database specified in the actions.resources.db field.

	-> **NOTE** This field is mutually exclusive with the `actions.resources.cluster` field.

* `resources.#.database_name`	Database on which the action is granted.

	-> **NOTE** This field is mutually exclusive with the `actions.resources.cluster` field.

* `resources.#.cluster`	(Optional) Set to true to indicate that the action is granted on the cluster resource.

	-> **NOTE** This field is mutually exclusive with the `actions.resources.collection` and `actions.resources.db fields`.

### Inherited Roles
Each object in the inheritedRoles array represents a key-value pair indicating the inherited role and the database on which the role is granted. It is an optional field.

* `database_name` (Required) Database on which the inherited role is granted.

	-> **NOTE** This value should be admin for all roles except read and readWrite.

* `role_name`	(Required) Name of the inherited role. This can either be another custom role or a built-in role.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier used for terraform for internal manages and can be used to import.

## Import

Database users can be imported using project ID and username, in the format `PROJECTID-ROLENAME`, e.g.

```
$ terraform import mongodbatlas_custom_db_role.my_role 1112222b3bf99403840e8934-MyCustomRole
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/custom-roles/)