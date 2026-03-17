# Step 1: Generate an ephemeral JWT and store it in AWS Secrets Manager.
# The JWT uses the provider's Service Account credentials by default.

resource "aws_secretsmanager_secret" "atlas_token" {
  name = var.secret_name
}

# Generate a short-lived Atlas JWT -- never written to state or plan.
# By default, this uses the provider's Service Account credentials.
ephemeral "mongodbatlas_service_account_jwt" "token" {}

# ---------------------------------------------------------------------------
# Using a specific Service Account for the JWT
# ---------------------------------------------------------------------------
# To generate a JWT with a different access level, create a dedicated
# Service Account and pass its credentials explicitly:
#
#   resource "mongodbatlas_service_account" "jwt_sa" {
#     org_id                     = var.org_id
#     name                       = "jwt-dedicated-sa"
#     description                = "SA used exclusively for ephemeral JWT generation."
#     roles                      = ["ORG_READ_ONLY"]
#     secret_expires_after_hours = 2160
#   }
#
#   resource "mongodbatlas_service_account_secret" "jwt_sa" {
#     org_id                     = var.org_id
#     client_id                  = mongodbatlas_service_account.jwt_sa.client_id
#     secret_expires_after_hours = 2160
#   }
#
#   ephemeral "mongodbatlas_service_account_jwt" "token" {
#     client_id     = mongodbatlas_service_account.jwt_sa.client_id
#     client_secret = mongodbatlas_service_account_secret.jwt_sa.secret
#   }
# ---------------------------------------------------------------------------

# ---------------------------------------------------------------------------
# Approach A: Write-only attribute (Terraform >= 1.11, recommended)
# ---------------------------------------------------------------------------
# Stores the JWT via secret_string_wo, which accepts ephemeral values.
# The token is sent to AWS but never stored in Terraform state.
# Increment var.token_version to rotate the stored token on subsequent applies.
# ---------------------------------------------------------------------------
resource "aws_secretsmanager_secret_version" "atlas_token" {
  secret_id                = aws_secretsmanager_secret.atlas_token.id
  secret_string_wo         = ephemeral.mongodbatlas_service_account_jwt.token.access_token
  secret_string_wo_version = var.token_version
}

# ---------------------------------------------------------------------------
# Approach B: local-exec provisioner (Terraform >= 1.10)
# ---------------------------------------------------------------------------
# You can use this approach instead of Approach A if you are on Terraform 1.10 or
# your cloud provider does not yet support write-only attributes.
#
# To switch to Approach B:
#
#   1. Set required_version = ">= 1.10" in versions.tf.
#
#   2. Replace:
#        resource "aws_secretsmanager_secret_version" "atlas_token" { ... }
#
#      With:
#        resource "terraform_data" "store_token" {
#          triggers_replace = [timestamp()]
#
#          provisioner "local-exec" {
#            command     = "aws secretsmanager put-secret-value --secret-id \"$SECRET_ID\" --secret-string \"$ATLAS_TOKEN\""
#            environment = {
#              SECRET_ID   = aws_secretsmanager_secret.atlas_token.id
#              ATLAS_TOKEN = ephemeral.mongodbatlas_service_account_jwt.token.access_token
#            }
#          }
#        }
#
# NOTE:
# - The token is passed via environment variable, never logged or stored in state.
# - Use triggers_replace = [timestamp()] to re-run on every apply and refresh the token.
# ---------------------------------------------------------------------------
