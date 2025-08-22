############################################################
# v2: New resource usage
############################################################

# Map of team IDs to their roles
locals {
  team_map = {
    var.team_id_1 = var.team_1_roles
    var.team_id_2 = var.team_2_roles
  }
}

# Ignore the deprecated teams block in mongodbatlas_project
resource "mongodbatlas_project" "this" {
  name   = "this"
  org_id = var.org_id
  lifecycle {
    ignore_changes = [teams]
  }
}

# Use the new mongodbatlas_team_project_assignment resource
resource "mongodbatlas_team_project_assignment" "this" {  
  for_each = local.team_map  
    
  project_id = mongodbatlas_project.this.id  
  team_id    = each.key  
  role_names = each.value  
} 

# Import existing team-project relationships into the new resource
import {
  for_each = local.team_map
  to = mongodbatlas_team_project_assignment.this[each.key]
  id = "${mongodbatlas_project.this.id}/${each.key}"
}

# Example outputs showing team assignments in various formats
output "team_project_assignments" {  
  description = "List of all team assignments for the MongoDB Atlas project"  
  value = [  
    for assignment in mongodbatlas_team_project_assignment.this :  
    {  
      team_id    = assignment.team_id  
      role_names = assignment.role_names  
    }  
  ]  
}  

output "team_project_assignments_map" {  
  description = "Map of team_id to role_names for the MongoDB Atlas project"  
  value = {  
    for k, assignment in mongodbatlas_team_project_assignment.this :  
    assignment.team_id => assignment.role_names  
  }  
}

# Data source to read current team assignments for the project
data "mongodbatlas_team_project_assignment" "this" {  
  project_id = mongodbatlas_project.this.id
  team_id = var.team_id_1 # Example for one team; repeat for others as needed
}

output "data_team_project_assignment" {  
  description = "Data source output for team assignment"  
  value = {  
    team_id    = data.mongodbatlas_team_project_assignment.this.team_id  
    project_id = data.mongodbatlas_team_project_assignment.this.project_id  
    role_names = data.mongodbatlas_team_project_assignment.this.role_names  
  }  
}
