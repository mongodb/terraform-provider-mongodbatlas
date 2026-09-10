resource "mongodbatlas_org_delegation_settings" "this" {
  org_id                         = var.atlas_org_id
  delegated_mcp_access           = "READ_ONLY"
  delegated_partner_access       = "DISALLOWED"
  idle_refresh_token_lifetime    = 3600
  maximum_refresh_token_lifetime = 86400
}

data "mongodbatlas_org_delegation_settings" "this" {
  org_id = mongodbatlas_org_delegation_settings.this.org_id
}

output "org_delegation_settings_delegated_mcp_access" {
  description = "The MCP delegated access policy of the organization."
  value       = mongodbatlas_org_delegation_settings.this.delegated_mcp_access
}

output "org_delegation_settings_delegated_partner_access" {
  description = "The partner delegated access policy of the organization."
  value       = mongodbatlas_org_delegation_settings.this.delegated_partner_access
}

output "org_delegation_settings_idle_refresh_token_lifetime" {
  description = "The maximum number of seconds a refresh token may be idle before it expires."
  value       = mongodbatlas_org_delegation_settings.this.idle_refresh_token_lifetime
}

output "org_delegation_settings_maximum_refresh_token_lifetime" {
  description = "The maximum lifetime of a refresh token in seconds, regardless of activity."
  value       = mongodbatlas_org_delegation_settings.this.maximum_refresh_token_lifetime
}

output "org_delegation_settings_delegated_mcp_access_ds" {
  description = "The MCP delegated access policy of the organization, from the data source."
  value       = data.mongodbatlas_org_delegation_settings.this.delegated_mcp_access
}
