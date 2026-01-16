# MongoDB Atlas Provider -- Project Service Account

This example shows how to create a Project Service Account in MongoDB Atlas.

## Important Notes

When you create a Project Service Account, Atlas automatically generates a secret. The secret value is returned only once, at creation time.

The example includes a sensitive output `service_account_first_secret` that captures this initial secret. 
You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw service_account_first_secret
```

For secret rotation, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

## Prerequisites
- Service Account with Project Owner permissions

## Variables Required to be set:
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret  
- `project_id`: Project ID where the Project Service Account will be created

## Outputs
- `service_account_client_id`: The Client ID of the created Project Service Account
- `service_account_name`: The name of the Project Service Account
- `service_account_first_secret` (sensitive): The initial secret value (only available at creation)
- `service_accounts_results`: All Project Service Accounts in the project
