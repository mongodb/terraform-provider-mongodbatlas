resource "mongodbatlas_team" "this" {
  org_id = var.org_id  
  name   = var.team_name 
}