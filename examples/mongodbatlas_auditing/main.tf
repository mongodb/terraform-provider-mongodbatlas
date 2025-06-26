# Specify an auditing resource and enable auditing for a project.
# To configure auditing, specify the unique project ID. If you change
# this value to a different "project_id", this deletes the current audit
# settings for the original project.

# "audit_authorization_success" indicates whether the auditing system
# captures successful authentication attempts for audit filters using
# the "atype" : "authCheck" auditing event. Warning! If you set
# "audit_authorization_success" to "true", this can severely impact
# cluster performance. Enable this option with caution.

# "audit_filter" is the JSON-formatted audit filter.
# "enabled" denotes whether or not the project associated with the
# specified "{project_id}"" has database auditing enabled. Defaults to "false".

# Auditing created by API Keys must belong to an existing organization.

# In addition to arguments listed previously, the following attributes
# are exported:

# "configuration_type" denotes the configuration method for the audit filter.
# Possible values are:
# - "NONE" - auditing is not configured for the project.
# - "FILTER_BUILDER" - auditing is configured via the Atlas UI filter builder.
# - "FILTER_JSON" - auditing is configured via a custom filter in Atlas or API.

locals {
  audit_filter_json = var.audit_filter_json != "" ? var.audit_filter_json : "${path.module}/audit_filter.json"
}
resource "mongodbatlas_auditing" "this" {
  project_id   = var.project_id
  audit_filter = file(local.audit_filter_json)

  audit_authorization_success = false
  enabled                     = true
}
