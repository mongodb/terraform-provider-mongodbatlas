output "user_from_username" {
  description = "Project assignment details for the user retrieved by username"
  value       = data.mongodbatlas_cloud_user_project_assignment.example_username
}

output "user_from_user_id" {
  description = "Project assignment details for the user retrieved by user_id"
  value       = data.mongodbatlas_cloud_user_project_assignment.example_user_id
}

output "assigned_user" {
  description = "Details of the assigned user"
  value       = mongodbatlas_cloud_user_project_assignment.example
}
