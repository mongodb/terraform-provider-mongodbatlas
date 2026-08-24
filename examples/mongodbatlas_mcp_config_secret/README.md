# MongoDB Atlas Provider -- MCP Config Secret

This example shows how to create an ingress secret for an organization-level MCP configuration.

## Important Notes

The secret value is returned only once, in the create response. The example includes a sensitive output `secret` that captures it.
You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw secret
```

## Variables Required to be set:
- `atlas_client_id`: The MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: The MongoDB Atlas Service Account Client Secret
- `org_id`: The organization ID where the MCP configuration will be created

## Outputs
- `secret_id`: The unique identifier of the created secret
- `secret` (sensitive): The plain-text secret value (only available at creation)
- `secret_expires_at`: The secret's expiry timestamp, retrieved via the singular data source
- `mcp_config_secrets_results`: All ingress secrets for the MCP configuration
