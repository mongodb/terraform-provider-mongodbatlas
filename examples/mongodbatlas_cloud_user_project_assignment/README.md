# Example: mongodbatlas_cloud_user_project_assignment  
  
This example demonstrates how to use the `mongodbatlas_cloud_user_project_assignment` resource to assign a user to a MongoDB Atlas project with specified roles.
  
## Usage  
  
```hcl  
provider "mongodbatlas" {  
  public_key  = var.public_key  
  private_key = var.private_key  
}  
  
resource "mongodbatlas_cloud_user_project_assignment" "example" {  
  project_id = var.project_id  
  username   = var.user_email  
  roles      = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]  
}  
```  
  
You must set the following variables:  
  
- `public_key`: Your MongoDB Atlas API public key.  
- `private_key`: Your MongoDB Atlas API private key.  
- `project_id`: The ID of the MongoDB Atlas project to assign the user to.  
- `user_email`: The email address of the user to assign to the project.   

## Considerations  
  
- To use this resource, the requesting Service Account or API Key must have the Project Owner role.  
- If the user has a pending invitation to join the project's organization, MongoDB Cloud modifies it and grants project access.  
- If the user doesn't have an invitation to join the organization, MongoDB Cloud sends a new invitation that grants the user organization and project access.  
- If the user is already active in the project's organization, MongoDB Cloud grants access to the project.  
  

To learn more, see the [MongoDB Atlas API - Cloud Users](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-addprojectuser) Documentation.  
