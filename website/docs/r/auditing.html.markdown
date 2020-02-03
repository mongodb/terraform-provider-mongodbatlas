---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: auditing"
sidebar_current: "docs-mongodbatlas-resource-auditing"
description: |-
    Provides a Auditing resource.
---

# mongodbatlas_auditing

`mongodbatlas_auditing` provides an Auditing resource. This allows auditing to be created.

## Example Usage

```hcl
resource "mongodbatlas_auditing" "test" {
		project_id                  = "<project-id>"
		audit_filter                = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		audit_authorization_success = false
		enabled                     = true
	}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to configure auditing.
* `audit_filter` - Indicates whether the auditing system captures successful authentication attempts for audit filters using the "atype" : "authCheck" auditing event. For more information, see auditAuthorizationSuccess
* `audit_authorization_success` - JSON-formatted audit filter used by the project
* `enabled` - Denotes whether or not the project associated with the {project_id} has database auditing enabled.

~> **NOTE:** Auditing created by API Keys must belong to an existing organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `configuration_type` - Denotes the configuration method for the audit filter. Possible values are: 
	* NONE - auditing not configured for the project.
	* FILTER_BUILDER - auditing configured via Atlas UI filter builder.
	* FILTER_JSON - auditing configured via Atlas custom filter or API.

## Import

Auditing must be imported using auditing ID, e.g.

```
$ terraform import mongodbatlas_auditing.my_auditing 5d09d6a59ccf6445652a444a
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/auditing/)