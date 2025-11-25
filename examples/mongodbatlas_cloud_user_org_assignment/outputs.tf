output "user_from_username" {
  description = "User details retrieved by username"
  value       = data.mongodbatlas_cloud_user_org_assignment.example_username
}

output "user_from_user_id" {
  description = "User details retrieved by user_id"
  value       = data.mongodbatlas_cloud_user_org_assignment.example_user_id
}

output "created_user" {
  description = "Details of the created user"
  value       = mongodbatlas_cloud_user_org_assignment.example
}
