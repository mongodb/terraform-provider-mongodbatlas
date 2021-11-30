---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project"
sidebar_current: "docs-mongodbatlas-resource-project"
description: |-
    Provides a Project resource.
---

# mongodbatlas_project

`mongodbatlas_project` provides a Project resource. This allows project to be created.

~> **IMPORTANT WARNING:**  Changing the name of an existing Project in your Terraform configuration will result the destruction of that Project and related resources (including Clusters) and the re-creation of those resources.  Terraform will inform you of the destroyed/created resources before applying so be sure to verify any change to your environment before applying.

## Example Usage

```terraform
resource "mongodbatlas_project" "test" {
  name   = "project-name"
  org_id = "<ORG_ID>"
  project_owner_id = "<OWNER_ACCOUNT_ID>"

  teams {
    team_id    = "5e0fa8c99ccf641c722fe645"
    role_names = ["GROUP_OWNER"]

  }
  teams {
    team_id    = "5e1dd7b4f2a30ba80a70cd4rw"
    role_names = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_WRITE"]
  }
}
```

## Argument Reference

* `name` - (Required) The name of the project you want to create. (Cannot be changed via this Provider after creation.)
* `org_id` - (Required) The ID of the organization you want to create the project within.
* `project_owner_id` - (Optional) Unique 24-hexadecimal digit string that identifies the Atlas user account to be granted the [Project Owner](https://docs.atlas.mongodb.com/reference/user-roles/#mongodb-authrole-Project-Owner) role on the specified project. If you set this parameter, it overrides the default value of the oldest [Organization Owner](https://docs.atlas.mongodb.com/reference/user-roles/#mongodb-authrole-Organization-Owner).
* `with_default_alerts_settings` - (Optional) It allows users to disable the creation of the default alert settings. By default, this flag is set to true.

### Teams
Teams attribute is optional

~> **NOTE:** Atlas limits the number of users to a maximum of 100 teams per project and a maximum of 250 teams per organization.

* `team_id` - (Required) The unique identifier of the team you want to associate with the project. The team and project must share the same parent organization.

* `role_names` - (Required) Each string in the array represents a project role you want to assign to the team. Every user associated with the team inherits these roles. You must specify an array even if you are only associating a single role with the team.
 The following are valid roles:
  * `GROUP_OWNER`
  * `GROUP_READ_ONLY`
  * `GROUP_DATA_ACCESS_ADMIN`
  * `GROUP_DATA_ACCESS_READ_WRITE`
  * `GROUP_DATA_ACCESS_READ_ONLY`
  * `GROUP_CLUSTER_MANAGER`


~> **NOTE:** Project created by API Keys must belong to an existing organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The project id.
* `created` - The ISO-8601-formatted timestamp of when Atlas created the project..
* `cluster_count` - The number of Atlas clusters deployed in the project..

## Import

Project must be imported using project ID, e.g.

```
$ terraform import mongodbatlas_project.my_project 5d09d6a59ccf6445652a444a
```
For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/projects/) - [and MongoDB Atlas API - Teams](https://docs.atlas.mongodb.com/reference/api/teams/) Documentation for more information.
