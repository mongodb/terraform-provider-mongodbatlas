# Project assignment outputs
output "assignment_id" {
  description = "ID of the user project assignment"
  value       = mongodbatlas_cloud_user_project_assignment.user_assignment.id
}

output "assigned_username" {
  description = "Username of the assigned user"
  value       = mongodbatlas_cloud_user_project_assignment.user_assignment.username
}

output "assigned_roles" {
  description = "Roles assigned to the user"
  value       = mongodbatlas_cloud_user_project_assignment.user_assignment.roles
}

output "user_id" {
  description = "MongoDB Atlas User ID (available with new resource)"
  value       = mongodbatlas_cloud_user_project_assignment.user_assignment.user_id
}

output "assignment_summary" {
  description = "Complete assignment summary"
  value       = local.assignment_summary
}

# Data source outputs (demonstrates read capability)
output "user_details_from_data_source" {
  description = "User details retrieved via data source"
  value       = local.user_details
}

# Demonstrates the advantages of the new resource
output "new_resource_advantages" {
  description = "Advantages of the new resource over deprecated project_invitation"
  value = {
    manages_active_membership = "Manages actual project membership, not just invitations"
    exposes_user_id           = "Provides user_id which wasn't available in project_invitation"
    supports_data_source      = "Has corresponding data source for reading assignments"
    import_capable            = "Can import existing project members"
    no_state_removal          = "Doesn't get removed from state when user accepts invitation"
  }
}

# Usage examples for common patterns
output "usage_examples" {
  description = "Common usage patterns with the new resource"
  value = {
    basic_assignment  = "Assign user to project with specific roles"
    read_assignment   = "Read existing user assignment from project"
    user_id_reference = "Use user_id for other resource references"
    role_management   = "Update user roles within project"
  }
}
