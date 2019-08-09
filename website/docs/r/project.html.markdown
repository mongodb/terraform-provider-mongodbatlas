---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project"
sidebar_current: "docs-mongodbatlas-resource-project"
description: |-
    Provides a Project resource.
---

# mongodb_atlas_project

`mongodbatlas_project` provides a Project resource. This allows project to be created.

## Example Usage

```hcl
resource "mongodbatlas_project" "my_project" {
	name   = "testacc-project"
	org_id = "5b93ff2f96e82120w0aaec19"
}
```

## Argument Reference

* `name` - (Required) The name of the project you want to create.
* `org_id` - (Required) The ID of the organization you want to create the project within.

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
For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/projects/)