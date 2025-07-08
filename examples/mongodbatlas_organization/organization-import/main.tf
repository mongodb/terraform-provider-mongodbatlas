# Example configuration for importing an existing MongoDB Atlas organization
# IMPORTANT: Do NOT include org_owner_id, description, or role_names when importing
# These are creation-only attributes and will cause errors if specified during import

resource "mongodbatlas_organization" "imported" {
  name = var.org_name

  # Optional settings - configure these to match your existing organization
  api_access_list_required     = var.api_access_list_required
  multi_factor_auth_required   = var.multi_factor_auth_required
  restrict_employee_access     = var.restrict_employee_access
  gen_ai_features_enabled      = var.gen_ai_features_enabled
  security_contact             = var.security_contact
  skip_default_alerts_settings = var.skip_default_alerts_settings

  # federation_settings_id can be set if your organization is federated
  # federation_settings_id = var.federation_settings_id
}

# Outputs for reference
output "org_id" {
  description = "The imported organization ID"
  value       = mongodbatlas_organization.imported.org_id
}

output "org_name" {
  description = "The imported organization name"
  value       = mongodbatlas_organization.imported.name
}

# Note: public_key and private_key are not available when importing
# These are only generated when creating a new organization 
