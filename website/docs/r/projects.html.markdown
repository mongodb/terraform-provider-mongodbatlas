---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: projects"
sidebar_current: "docs-mongodbatlas-resource-projects"
description: |-
		Provides a Projects resource.
---

# mongodb_atlas_projects

`mongodbatlas_projects` provides a Project resource. This allows projects to be created.

## Example Usage

```hcl
resource "mongodbatlas_projects" "my_project" {
	name   = "testacc-project"
	org_id = "5b93ff2f96e82120w0aaec19"
}
```

## Argument Reference

* `name` - (Required) The name of the project you want to create.
* `org_id` - (Required) The ID of the organization you want to create the project within.

~> **NOTE:** Projects created by API Keys must belong to an existing organization.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The project id.

## Import

Projects must be imported using project ID, e.g.

```
$ terraform import mongodbatlas_projects.my_project 5d09d6a59ccf6445652a444a
```
