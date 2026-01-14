# MongoDB Atlas Provider -- Project Service Account

This example shows how to create a Project Service Account in MongoDB Atlas.

## Important Notes

When a Project Service Account is created, a secret is automatically generated. The secret value is only returned once at creation time.

The example includes a sensitive output `service_account_first_secret` that captures this initial secret. 
You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw service_account_first_secret
```

For secret rotation, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

~> **IMPORTANT WARNING:** Managing Service Accounts with Terraform **exposes sensitive organizational secrets** in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data). You may also want to consider managing your Service Accounts via a more secure method, such as the [HashiCorp Vault MongoDB Atlas Secrets Engine](https://developer.hashicorp.com/vault/docs/secrets/mongodbatlas).

~> **IMPORTANT NOTE** The use of `mongodbatlas_project_service_account` resource is no longer the recommended approach for users with Organization Owner permissions. For new configurations, we recommend using the `mongodbatlas_service_account` resource and the `mongodbatlas_service_account_project_assignment` resource to assign the Service Account to projects. This approach is more flexible and aligns with best practices.

~> **IMPORTANT:** When you delete a `mongodbatlas_project_service_account` resource, it deassigns the Service Account from the project but does not delete the Service Account itself. The Service Account will remain in your organization and can be reassigned to projects or deleted separately if needed.

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
