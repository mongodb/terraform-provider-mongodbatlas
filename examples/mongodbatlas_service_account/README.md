# MongoDB Atlas Provider -- Service Account

This example shows how to create a Service Account in MongoDB Atlas.

## Important Notes

When you create a Service Account, Atlas automatically generates a secret. The secret value returns only once, at creation time.

The example includes a sensitive output `service_account_first_secret` that captures this initial secret. 
You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw service_account_first_secret
```

For secret rotation, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

## Variables Required to be set:
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret  
- `org_id`: Organization ID where the Service Account will be created

## Outputs
- `service_account_client_id`: The Client ID of the created Service Account
- `service_account_name`: The name of the Service Account
- `service_account_first_secret` (sensitive): The initial secret value (only available at creation)
- `service_accounts_results`: All Service Accounts in the organization
