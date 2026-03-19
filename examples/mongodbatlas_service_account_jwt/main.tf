# Generate a short-lived Atlas JWT using the provider's Service Account credentials.
# The token is never written to Terraform state or plan.
ephemeral "mongodbatlas_service_account_jwt" "token" {
  revoke_on_closure = true
}

# Use the ephemeral token to configure a second provider instance.
provider "mongodbatlas" {
  alias        = "ephemeral"
  access_token = ephemeral.mongodbatlas_service_account_jwt.token.access_token
}
