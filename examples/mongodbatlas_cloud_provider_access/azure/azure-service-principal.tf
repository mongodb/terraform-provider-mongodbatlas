# Use the following block to create the Service Principal 
resource "azuread_service_principal" "example" {
  application_id               = var.atlas_azure_app_id
  app_role_assignment_required = false
}
