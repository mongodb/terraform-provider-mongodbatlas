# MongoDB Atlas Provider -- Project MCP Config

This example shows how to create a project-level Remote MCP configuration in MongoDB Atlas.

## Important Notes

This resource does not manage ingress secrets. Use [`mongodbatlas_project_mcp_config_secret`](../mongodbatlas_project_mcp_config_secret/README.md) to create one.

This resource only unlinks the MCP configuration from the project. the underlying MCP configuration is not deleted. To fully delete it, use the org-level [`mongodbatlas_mcp_config`](../mongodbatlas_mcp_config/README.md) resource.

## Variables Required to be set:
- `atlas_client_id`: The MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: The MongoDB Atlas Service Account Client Secret
- `project_id`: The project ID where the MCP configuration will be created

## Outputs
- `mcp_config_id`: The unique identifier of the created MCP configuration
- `mcp_config_name`: The name of the MCP configuration, retrieved via the singular data source
- `client_id`: The ingress Service Account client ID
- `egress_client_id`: The egress Service Account client ID, used for audit/compliance purposes
- `mcp_configs_results`: All MCP configurations in the project
