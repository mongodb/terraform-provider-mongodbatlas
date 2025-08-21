# Comparison outputs to validate migration
output "email_mapping_comparison" {
  description = "Compare old vs new email retrieval"
  value = {
    old_email = local.old_user_email_by_id
    new_email = local.new_user_email_by_id
    matches   = local.email_mapping_matches
  }
}

output "org_roles_comparison" {
  description = "Compare old vs new organization roles"
  value = {
    old_roles = local.old_user_org_roles
    new_roles = local.new_user_org_roles
  }
}

output "project_roles_comparison" {
  description = "Compare old vs new project roles"
  value = {
    old_roles = local.old_user_project_roles
    new_roles = local.new_user_project_roles
  }
}

output "org_users_comparison" {
  description = "Compare old vs new organization user lists"
  value = {
    old_emails    = local.old_org_user_emails
    new_emails    = local.new_org_user_emails
    old_count     = length(local.old_org_user_emails)
    new_count     = length(local.new_org_user_emails)
    count_matches = local.org_users_count_matches
  }
}

output "project_users_comparison" {
  description = "Compare old vs new project user lists"
  value = {
    old_emails = local.old_project_user_emails
    new_emails = local.new_project_user_emails
    old_count  = length(local.old_project_user_emails)
    new_count  = length(local.new_project_user_emails)
  }
}

output "team_users_comparison" {
  description = "Compare old vs new team user lists"
  value = {
    old_emails = local.old_team_user_emails
    new_emails = local.new_team_user_emails
    old_count  = length(local.old_team_user_emails)
    new_count  = length(local.new_team_user_emails)
  }
}

# Migration validation
output "migration_validation" {
  description = "Overall migration validation results"
  value = {
    email_mapping_works     = local.email_mapping_matches
    org_users_count_matches = local.org_users_count_matches
    ready_for_v3            = local.email_mapping_matches && local.org_users_count_matches
  }
}
