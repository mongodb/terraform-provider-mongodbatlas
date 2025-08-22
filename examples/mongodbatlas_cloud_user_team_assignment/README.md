# Example: mongodbatlas_cloud_user_team_assignment

This example demonstrates how to use the `mongodbatlas_cloud_user_team_assignment` resource to assign a user to a team within a MongoDB Atlas organization.

## Usage  
  
```hcl  
provider "mongodbatlas" {  
  public_key  = var.public_key  
  private_key = var.private_key  
}  
  
resource "mongodbatlas_cloud_user_team_assignment" "example" {  
  org_id  = var.org_id  
  team_id = var.team_id  
  user_id = var.user_id  
}  
```  

You must set the following variables:  
- `public_key`: Your MongoDB Atlas API public key.  
- `private_key`: Your MongoDB Atlas API private key.  
- `org_id`: The ID of the MongoDB Atlas organization.  
- `team_id`: The ID of the team to assign the user to.  
- `user_id`: The ID of the user to assign to the team.  


To learn more, see the [MongoDB Atlas API - Cloud Users](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-addusertoteam) Documentation.
