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