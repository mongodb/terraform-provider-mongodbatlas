# Original invitation outputs
output "invitation_id" {
  description = "ID of the pending project invitation"
  value       = mongodbatlas_project_invitation.pending_user.invitation_id
}

output "invited_username" {
  description = "Username of the invited user"
  value       = mongodbatlas_project_invitation.pending_user.username
}

output "assigned_roles" {
  description = "Roles assigned to the invited user"
  value       = mongodbatlas_project_invitation.pending_user.roles
}

output "invitation_details" {
  description = "Complete invitation details"
  value       = local.invitation_details
}

output "creation_timestamp" {
  description = "When the invitation was created"
  value       = mongodbatlas_project_invitation.pending_user.created_at
}

output "expiration_timestamp" {
  description = "When the invitation expires"
  value       = mongodbatlas_project_invitation.pending_user.expires_at
}
