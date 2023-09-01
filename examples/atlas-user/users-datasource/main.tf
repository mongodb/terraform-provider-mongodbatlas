data "mongodbatlas_atlas_users" "test_org_users" {
  org_id = var.org_id
}

data "mongodbatlas_atlas_users" "test_project_users" {
  project_id = var.project_id
}

data "mongodbatlas_atlas_users" "test_team_users" {
  team_id = var.team_id
  org_id  = var.org_id
}

# example making use of data sources
output "org_user_count" {
  value = data.mongodbatlas_atlas_users.test_org_users.total_count
}

output "project_user_count" {
  value = data.mongodbatlas_atlas_users.test_project_users.total_count
}

output "team_user_count" {
  value = data.mongodbatlas_atlas_users.test_team_users.total_count
}