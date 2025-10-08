---
page_title: "Migration Guide: Service Accounts"
---

# Migration Guide: Service Accounts Authentication

This guide helps you migrate from Programmatic Access Key (PAK) authentication to Service
Accounts (SA) authentication and viceversa without impacting your deployment.

**Note:** For more information on SA, see [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/)
in the MongoDB documentation.

## Procedure

To migrate from Programmatic Access Key (PAK) authentication to Service
Accounts (SA) authentication, change your provider declaration variables. You can implement
this change by either:

- Providing a client ID and secret

- Providing a valid access token

### Provide a Client ID and Secret

The following example shows the variables for PAK authentication:

```terraform
provider "mongodbatlas" {
  public_key  = var.mongodbatlas_public_key
  private_key = var.mongodbatlas_private_key
}
```

To change to SA, declare the `client_id` and `client_secret` variables as in the following example:

```terraform
provider "mongodbatlas" {
client_id = var.mongodbatlas_client_id
client_secret = var.mongodbatlas_client_secret
}
```

### Provide a Valid Access Token

The following example shows SA authentication set up through the ``access_token`` attribute:

```terraform
provider "mongodbatlas" { 
access_token = var.mongodbatlas_access_token
[is_mongodbgov_cloud = true // optional]
}
```

Consider that the access token is **valid for one hour only**.

See [Generate Service Account Token](https://www.mongodb.com/docs/atlas/api/service-accounts/generate-oauth2-token/#std-label-generate-oauth2-token-atlas) for more details on creating an SA token. 

**IMPORTANT:**  Currently, the MongoDB Terraform provider does not support additional Token OAuth features like scopes.

**NOTE:** You can't use ``mongodbatlas_event_trigger`` with Service Accounts as the authentication method.
