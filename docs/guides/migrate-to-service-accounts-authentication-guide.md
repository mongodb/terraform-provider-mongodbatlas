---
page_title: "Migration Guide: Service Accounts Authentication"
---

# Migration Guide: Service Accounts Authentication

This guide helps you migrate from Programmatic Access Key (PAK) authentication to Service
Accounts (SA) authentication and viceversa without impacting your deployment. 

**Note:** For more information on SA, see [Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/)
in the MongoDB documentation.

## Procedure

1. Change your provider declaration variables.

    The following example declares PAK authentication:

    ```terraform
    provider "mongodbatlas" {
    public_key = var.mongodbatlas_public_key
    private_key  = var.mongodbatlas_private_key
    ```

    To change to SA, declare the variables as in the following example:

    ```terraform
    provider "mongodbatlas" {
    client_id = var.mongodbatlas_client_id
    client_secret  = var.mongodbatlas_client_secret
    }
    ```

2. Provide a valid JSON Web Token (JWT):

   ```terraform
   provider "mongodbatlas" { 
   access_token = var.mongodbatlas_access_token
   [is_mongodbgov_cloud = true // optional]
   }
   ```

    The JWT token is only valid during its set duration time. See [Generate Service Account Token](https://www.mongodb.com/docs/atlas/api/service-accounts/generate-oauth2-token/#std-label-generate-oauth2-token-atlas) for more details on creating an SA token.

   **IMPORTANT:**  Currently, the MongoDB Terraform provider does not support additional Token OAuth features.

   **NOTE:** You can not use ``mongodbatlas_event_trigger`` with Service Accounts as the authentication method.





