# Example: mongodbatlas_cloud_user_team_assignment

This example demonstrates how to use the `mongodbatlas_cloud_user_team_assignment` resource to assign a user to a team within a MongoDB Atlas organization.

## Usage  
  
```hcl  
provider "mongodbatlas" {  
  client_id     = var.atlas_client_id  
  client_secret = var.atlas_client_secret  
}  
  
resource "mongodbatlas_cloud_user_team_assignment" "example" {  
  org_id  = var.org_id  
  team_id = var.team_id  
  user_id = var.user_id  
}  
```  

You must set the following variables:  
- `atlas_client_id`: Your MongoDB Atlas Service Account Client ID.  
- `atlas_client_secret`: Your MongoDB Atlas Service Account Client Secret.  
- `org_id`: The ID of the MongoDB Atlas organization.  
- `team_id`: The ID of the team to assign the user to.  
- `user_id`: The ID of the user to assign to the team.  


To learn more, see the [MongoDB Atlas API - Cloud Users](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-addusertoteam) Documentation.
