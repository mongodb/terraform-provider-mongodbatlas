# Example: mongodbatlas_team_project_assignment  
  
This example demonstrates how to use the `mongodbatlas_team_project_assignment` resource to assign a team to an existing project with specified roles in MongoDB Atlas.  
  
## Usage  
  
```hcl  
provider "mongodbatlas" {  
  public_key  = var.public_key  
  private_key = var.private_key  
}  
  
resource "mongodbatlas_team_project_assignment" "example" {  
  project_id = var.project_id  
  team_id    = var.team_id  
  role_names = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]  
}  
  
data "mongodbatlas_team_project_assignment" "example_username" {  
  project_id = var.project_id  
  team_id    = var.team_id  
}  
```  
  
You must set the following variables:  
  
- `public_key`: Your MongoDB Atlas API public key.  
- `private_key`: Your MongoDB Atlas API private key.  
- `project_id`: The ID of the project to assign the team to.  
- `team_id`: The ID of the team to assign to the project.  
  
To learn more, see the [MongoDB Atlas API - Cloud Users](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-addallteamstoproject) Documentation.  

