# Example configuration for importing an existing MongoDB Atlas organization

resource "mongodbatlas_organization" "imported" {
  name = var.org_name

  # Optional settings - configure these to match your existing organization
  api_access_list_required     = var.api_access_list_required
  multi_factor_auth_required   = var.multi_factor_auth_required
  restrict_employee_access     = var.restrict_employee_access
  gen_ai_features_enabled      = var.gen_ai_features_enabled
  security_contact             = var.security_contact
  skip_default_alerts_settings = var.skip_default_alerts_settings
}

# Outputs for reference
output "org_name" {
  description = "The imported organization name"
  value       = mongodbatlas_organization.imported.name
}

# Note: public_key and private_key are not available when importing
# These are only generated when creating a new organization 
