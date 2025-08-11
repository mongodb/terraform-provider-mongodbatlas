output "assigned_team" {
  description = "Details of the assigned team"
  value       = mongodbatlas_team_project_assignment.example
}
output "team_from_team_id" {
  description = "Project assignment details for the team retrieved by team_id"
  value       = data.mongodbatlas_team_project_assignment.example_username
}
