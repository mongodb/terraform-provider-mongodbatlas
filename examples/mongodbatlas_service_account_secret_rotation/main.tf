# Rotate a Service Account secret while a consumer keeps authenticating: the account carries two
# secrets, and rotating one leaves the other valid through the switch.

resource "mongodbatlas_service_account" "this" {
  org_id                 = var.org_id
  name                   = var.service_account_name
  description            = "Service Account rotated by this example"
  roles                  = var.service_account_roles
  without_initial_secret = true
}

resource "mongodbatlas_service_account_secret" "secret_1" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = var.secret_expires_after_hours
}

resource "mongodbatlas_service_account_secret" "secret_2" {
  org_id                     = var.org_id
  client_id                  = mongodbatlas_service_account.this.client_id
  secret_expires_after_hours = var.secret_expires_after_hours

  # Serialize the creates to avoid DATA_CONCURRENCY_ERROR (HTTP 409).
  depends_on = [mongodbatlas_service_account_secret.secret_1]
}

# Read the live secret metadata from Atlas. A resource expires_at refreshes only when that
# resource is read, so the slot that is not rotating keeps a stale value in state. This data
# source always returns the current expiry for both slots.
data "mongodbatlas_service_account" "this" {
  org_id    = var.org_id
  client_id = mongodbatlas_service_account.this.client_id

  depends_on = [
    mongodbatlas_service_account_secret.secret_1,
    mongodbatlas_service_account_secret.secret_2,
  ]
}

locals {
  # Creating a secret cuts every existing secret to min(remaining lifetime, 7 days), so the slot
  # with the largest expires_at is the one created or rotated last. ISO 8601 UTC strings are
  # fixed width, so lexical order matches chronological order.
  live_secrets    = data.mongodbatlas_service_account.this.secrets
  by_secret_id    = { for s in local.live_secrets : s.secret_id => s }
  freshest_expiry = sort([for s in local.live_secrets : s.expires_at])[length(local.live_secrets) - 1]
  freshest_secret_id = [
    for s in local.live_secrets : s.secret_id if s.expires_at == local.freshest_expiry
  ][0]

  slots = {
    secret_1 = mongodbatlas_service_account_secret.secret_1.secret_id
    secret_2 = mongodbatlas_service_account_secret.secret_2.secret_id
  }

  # The freshest slot is the credential consumers adopt after a rotation.
  current_credentials = {
    client_id     = mongodbatlas_service_account.this.client_id
    client_secret = local.freshest_secret_id == local.slots.secret_1 ? mongodbatlas_service_account_secret.secret_1.secret : mongodbatlas_service_account_secret.secret_2.secret
  }
}

output "client_id" {
  description = "The Client ID of the Service Account. Constant across rotations."
  value       = mongodbatlas_service_account.this.client_id
}

output "current_credentials" {
  description = "The credential consumers adopt after a rotation, resolved from the slot with the largest expires_at."
  sensitive   = true
  value       = local.current_credentials
}

output "expires_at" {
  description = "Live Atlas expiry for each slot, read through the data source. The non-rotated slot shows the 7-day cut."
  value = {
    secret_1 = local.by_secret_id[local.slots.secret_1].expires_at
    secret_2 = local.by_secret_id[local.slots.secret_2].expires_at
  }
}

output "secret_1" {
  description = "The secret value of slot 1. Returned only when slot 1 is created."
  sensitive   = true
  value       = mongodbatlas_service_account_secret.secret_1.secret
}

output "secret_2" {
  description = "The secret value of slot 2. Returned only when slot 2 is created."
  sensitive   = true
  value       = mongodbatlas_service_account_secret.secret_2.secret
}
