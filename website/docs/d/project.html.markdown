---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: projects"
sidebar_current: "docs-mongodbatlas-datasource-project"
description: |-
    Describes a Projects
---

# mongodb_atlas_project

`mongodb_atlas_project` describe a project. This represents a project which will be applied to a organization.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
	resource "mongodbatlas_projects" "test" {
		name   = "test-datasource-project"
		org_id = "5b71ff2f96e82120d0aaec14"
	}

	data "mongodbatlas_project" "test" {
			name = "${mongodbatlas_projects.test.name}"
		}
```

## Argument Reference

* `name` - (Required) The name of the project you want to create.
* `org_id` - (Required) The ID of the organization you want to create the project within.

~> **NOTE:** Projects created by API Keys must belong to an existing organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The project id.


See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/projects/) Documentation for more information.