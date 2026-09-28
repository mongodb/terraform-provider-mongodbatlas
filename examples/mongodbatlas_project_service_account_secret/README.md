# MongoDB Atlas Provider -- Project Service Account Secret

This example shows how to create a Project Service Account Secret.

## Important Notes

Atlas returns each Project Service Account secret value only once, at creation time. This example creates a secret through `mongodbatlas_project_service_account_secret`, so Terraform manages the secret and you can read its value after the first apply.

The example includes a sensitive output `secret` that captures the value of that managed secret. You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw secret
```

For managing and rotating both secrets, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

## Prerequisites
- Service Account with Project Owner permissions used for Provider Authentication

## Variables Required to be set:
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_id`: Project ID where the Project Service Account will be created

## Outputs
- `secret_id`: The ID of the created secret
- `secret` (sensitive): The secret value
- `secret_expires_at`: The expiration date of the secret
