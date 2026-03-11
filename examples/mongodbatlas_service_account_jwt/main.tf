# Obtain a short-lived JWT from the provider's Service Account credentials.
# The token is ephemeral: it exists only for the duration of the Terraform
# operation and is never written to state.
ephemeral "mongodbatlas_service_account_jwt" "token" {
  revoke_on_closure = true
}

# Use a null_resource with a local-exec provisioner to consume the ephemeral token.
# This prints the token_type and expires_in to confirm the token was issued.
resource "null_resource" "verify_token" {
  provisioner "local-exec" {
    command = "echo 'Token type: ${ephemeral.mongodbatlas_service_account_jwt.token.token_type}, Expires in: ${ephemeral.mongodbatlas_service_account_jwt.token.expires_in}s'"
  }
}
