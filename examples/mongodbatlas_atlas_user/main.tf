data "mongodbatlas_atlas_user" "test_user_by_username" {
  username = var.username
}

data "mongodbatlas_atlas_user" "test_user_by_id" {
  user_id = var.user_id
}

# example making use of data sources
output "user_firstname" {
  value = data.mongodbatlas_atlas_user.test_user_by_username.first_name
}

output "user_lastname" {
  value = data.mongodbatlas_atlas_user.test_user_by_id.last_name
}