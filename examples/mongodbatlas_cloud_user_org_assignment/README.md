# Example: mongodbatlas_cloud_user_org_assignment

This example demonstrates how to use the `mongodbatlas_cloud_user_org_assignment` resource to assign a user to an existing organization with specified roles in MongoDB Atlas.

## Usage

```hcl
provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

resource "mongodbatlas_cloud_user_org_assignment" "example" {
  org_id   = var.org_id
  username = var.user_email
  roles = {
    org_roles = ["ORG_MEMBER"]
  }
}
```

You must set the following variables:

- `atlas_client_id`: Your MongoDB Atlas Service Account Client ID.
- `atlas_client_secret`: Your MongoDB Atlas Service Account Client Secret.
- `org_id`: The ID of the organization to assign the user to.
- `user_email`: The email address of the user to assign.

To learn more, see the [MongoDB Cloud Users Documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createorganizationuser).