# MongoDB Atlas Provider -- Delegation Settings

This example shows how to configure the delegation settings of a MongoDB Atlas organization with the `mongodbatlas_delegation_settings` resource. Delegation settings control how MCP (Model Context Protocol) and partner delegated access are permitted within the organization, as well as the refresh token lifetimes.

## Important Notes

Delegation settings are a singleton at the organization level: the settings always exist, so creating the resource updates the existing settings and destroying it only removes the resource from the Terraform state without changing the settings in Atlas.

## Variables Required to be set:
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `atlas_org_id`: Organization ID where the Delegation Settings will be configured

## Outputs
- `delegation_settings_delegated_mcp_access`: The MCP delegated access policy of the organization
- `delegation_settings_delegated_partner_access`: The partner delegated access policy of the organization
- `delegation_settings_idle_refresh_token_lifetime`: The maximum number of seconds a refresh token may be idle before it expires
- `delegation_settings_maximum_refresh_token_lifetime`: The maximum lifetime of a refresh token in seconds, regardless of activity
