# MongoDB Atlas Provider -- Project Service Account Secret

This example shows how to create a Project Service Account Secret.

## Important Notes

When you create a Project Service Account, Atlas automatically generates a secret. The secret value is returned only once, at creation time. This example creates a second secret for the same Project Service Account.

The example includes a sensitive output `secret` that captures this value. You can retrieve it using (**warning**: this prints the secret to your terminal):

```bash
terraform output -raw secret
```

For managing and rotating both secrets, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).

## Prerequisites
- Service Account with Project Owner permissions

## Variables Required to be set:
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_id`: Project ID where the Project Service Account will be created

## Outputs
- `secret_id`: The ID of the created secret
- `secret` (sensitive): The secret value
- `secret_expires_at`: The expiration date of the secret
