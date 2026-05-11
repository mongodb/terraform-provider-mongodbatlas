# Minimal example: manage an Atlas Service Account via the generic resource.
# A typed `mongodbatlas_service_account` resource exists — prefer it for production.

resource "mongodbatlas_api_resource" "this" {
  path                  = "/api/atlas/v2/orgs/${var.org_id}/serviceAccounts"
  id_attribute          = ["clientId"]
  create_only_body_keys = ["secretExpiresAfterHours"]

  body = {
    name                    = "example-generic-sa"
    description             = "Service account managed via mongodbatlas_api_resource"
    roles                   = ["ORG_READ_ONLY"]
    secretExpiresAfterHours = 24
  }
}

data "mongodbatlas_api_resource" "this" {
  path = mongodbatlas_api_resource.this.id
}

output "service_account_client_id" {
  value = mongodbatlas_api_resource.this.output.clientId
}

output "service_account_first_secret" {
  description = "Returned only on Create. After the first refresh this is null."
  value       = try(mongodbatlas_api_resource.this.output.secrets[0].secret, null)
  sensitive   = true
}
