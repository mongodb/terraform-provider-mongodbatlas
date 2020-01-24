---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project"
sidebar_current: "docs-mongodbatlas-datasource-project"
description: |-
    Describes a Project.
---

# mongodbatlas_project

`mongodbatlas_project` describe a Project. This represents a project created.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_project" "test" {
  name   = "project-name"
  org_id = "<ORG_ID>"

  teams {
    team_id    = "5e0fa8c99ccf641c722fe645"
    role_names = ["GROUP_OWNER"]

  }
  teams {
    team_id    = "5e1dd7b4f2a30ba80a70cd4rw"
    role_names = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_WRITE"]
  }
}

data "mongodbatlas_project" "test" {
  project_id = "${mongodbatlas_project.test.id}"
}
```

## Argument Reference

* `project_id` - The unique ID for the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the project you want to create. (Cannot be changed via this Provider after creation.)
* `org_id` - The ID of the organization you want to create the project within.

* `teams.#.team_id` - The unique identifier of the team you want to associate with the project. The team and project must share the same parent organization.

* `teams.#.role_names` - Each string in the array represents a project role assigned to the team. Every user associated with the team inherits these roles.
The following are the valid roles and their associated mappings:

* `GROUP_OWNER`
* `GROUP_READ_ONLY`
* `GROUP_DATA_ACCESS_ADMIN`
* `GROUP_DATA_ACCESS_READ_WRITE`
* `GROUP_DATA_ACCESS_READ_ONLY`


See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/project-get-one/) Documentation for more information.