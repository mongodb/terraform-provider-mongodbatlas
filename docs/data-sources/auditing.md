# Data Source: mongodbatlas_auditing

`mongodbatlas_auditing` describes a Auditing.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

```terraform
resource "mongodbatlas_auditing" "test" {
	project_id                  = "<project-id>"
	audit_filter                = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
	audit_authorization_success = false
	enabled                     = true
}

data "mongodbatlas_auditing" "test" {
			project_id = mongodbatlas_auditing.test.id
		}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `configuration_type` - Denotes the configuration method for the audit filter. Possible values are: NONE - auditing not configured for the project.m FILTER_BUILDER - auditing configured via Atlas UI filter builderm FILTER_JSON - auditing configured via Atlas custom filter or API.
* `audit_filter` - Indicates whether the auditing system captures successful authentication attempts for audit filters using the "atype" : "authCheck" auditing event. For more information, see auditAuthorizationSuccess
* `audit_authorization_success` - JSON-formatted audit filter used by the project
* `enabled` - Denotes whether or not the project associated with the {GROUP-ID} has database auditing enabled.


See detailed information for arguments and attributes: [MongoDB API Auditing](https://docs.atlas.mongodb.com/reference/api/auditing/)