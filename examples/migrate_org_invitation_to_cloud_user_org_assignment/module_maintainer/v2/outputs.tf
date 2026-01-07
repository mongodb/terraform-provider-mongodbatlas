output "user_id" {
  description = "User ID of the cloud user org assignment"
  value       = mongodbatlas_cloud_user_org_assignment.this.user_id
}

output "username" {
  description = "Username of the cloud user org assignment"
  value       = mongodbatlas_cloud_user_org_assignment.this.username
}

output "team_ids" {
  description = "Set of team IDs the user is assigned to"
  value       = var.team_ids
}

