############################################################
# v1: Original configuration using deprecated attribute
############################################################

# Map of team IDs to their roles
locals {
  team_map = {
    var.team_id_1 = var.team_1_roles
    var.team_id_2 = var.team_2_roles
  }
}

# Using deprecated team block inside mongodbatlas_project to assign teams to the project
resource "mongodbatlas_project" "this" {
  name   = "this"
  org_id = var.org_id

  dynamic "teams" {
    for_each = local.team_map
    content {
      team_id = teams.key
      role_names = teams.value
    }
  }
}