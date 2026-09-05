# MongoDB Atlas Provider -- MCP Config

This example shows how to create an organization-level Remote MCP configuration in MongoDB Atlas.

## Important Notes

This resource does not manage ingress secrets. Use `mongodbatlas_mcp_config_secret` to create one.

## Required Variables
- `atlas_client_id`: The MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: The MongoDB Atlas Service Account Client Secret
- `org_id`: The organization ID where the MCP configuration will be created

## Outputs
- `mcp_config_id`: The unique identifier of the created MCP configuration
- `mcp_config_name`: The name of the MCP configuration, retrieved via the singular data source
- `client_id`: The ingress Service Account client ID
- `egress_client_id`: The egress Service Account client ID, used for audit/compliance purposes
- `mcp_configs_results`: All MCP configurations in the organization
