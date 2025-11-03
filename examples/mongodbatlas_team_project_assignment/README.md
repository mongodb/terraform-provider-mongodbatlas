# Example: mongodbatlas_team_project_assignment  
  
This example demonstrates how to use the `mongodbatlas_team_project_assignment` resource to assign a team to an existing project with specified roles in MongoDB Atlas.  
  
## Usage  
  
```hcl  
provider "mongodbatlas" {  
  client_id     = var.atlas_client_id  
  client_secret = var.atlas_client_secret  
}  
  
resource "mongodbatlas_team_project_assignment" "this" {  
  project_id = var.project_id  
  team_id    = var.team_id  
  role_names = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]  
}  
  
data "mongodbatlas_team_project_assignment" "this" {  
  project_id = mongodbatlas_team_project_assignment.this.project_id
  team_id    = mongodbatlas_team_project_assignment.this.team_id
}  
```  
  
You must set the following variables:  
  
- `atlas_client_id`: Your MongoDB Atlas Service Account Client ID.  
- `atlas_client_secret`: Your MongoDB Atlas Service Account Client Secret.  
- `project_id`: The ID of the project to assign the team to.  
- `team_id`: The ID of the team to assign to the project.  
  
To learn more, see the [MongoDB Atlas API - Cloud Users](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-addallteamstoproject) Documentation.  

