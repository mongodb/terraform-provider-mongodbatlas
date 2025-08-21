# New resource outputs
output "user_assignment_id" {
  description = "ID of the new user project assignment"
  value       = mongodbatlas_cloud_user_project_assignment.user_assignment.id
}

output "assigned_username" {
  description = "Username from the new assignment"
  value       = mongodbatlas_cloud_user_project_assignment.user_assignment.username
}

output "assigned_roles" {
  description = "Roles from the new assignment"
  value       = mongodbatlas_cloud_user_project_assignment.user_assignment.roles
}

output "user_id" {
  description = "User ID from the new assignment (not available in old resource)"
  value       = mongodbatlas_cloud_user_project_assignment.user_assignment.user_id
}

# Migration validation outputs
output "migration_validation" {
  description = "Validation results for the migration"
  value = {
    username_matches     = local.username_matches
    roles_match          = local.roles_match
    migration_successful = local.migration_successful
    ready_for_v3         = local.migration_successful
  }
}

output "migration_comparison" {
  description = "Compare configuration inputs vs actual assignment"
  value = {
    input_username  = var.username
    input_roles     = var.roles
    actual_username = local.new_assignment_user
    actual_roles    = local.new_assignment_roles
  }
}

# New capabilities not available in old resource
output "new_capabilities" {
  description = "New capabilities available with cloud_user_project_assignment"
  value = {
    user_id_available    = mongodbatlas_cloud_user_project_assignment.user_assignment.user_id != null
    manages_active_users = true
    supports_import      = true
  }
}
