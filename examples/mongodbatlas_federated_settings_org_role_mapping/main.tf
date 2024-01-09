data "mongodbatlas_federated_settings" "federated_settings" {
  org_id = var.org_id
}
data "mongodbatlas_federated_settings_identity_providers" "identity_provider" {
  federation_settings_id = data.mongodbatlas_federated_settings.federated_settings.id
}

data "mongodbatlas_federated_settings_org_configs" "org_configs_ds" {
  federation_settings_id = data.mongodbatlas_federated_settings.federated_settings.id
}

data "mongodbatlas_federated_settings_org_role_mappings" "org_role_mapping" {
  federation_settings_id = data.mongodbatlas_federated_settings.federated_settings.id
  org_id                 = var.org_id
}
resource "mongodbatlas_federated_settings_org_role_mapping" "org_role_mapping" {
  federation_settings_id = data.mongodbatlas_federated_settings.federated_settings.id
  org_id                 = var.org_id
  external_group_name    = "newgroup"

  role_assignments {
    group_id = var.group_id
    roles    = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN", "GROUP_SEARCH_INDEX_EDITOR", "GROUP_DATA_ACCESS_READ_ONLY"]
  }

  role_assignments {
    org_id = var.org_id
    roles  = ["ORG_OWNER", "ORG_MEMBER"]
  }

}
resource "mongodbatlas_federated_settings_org_config" "org_connections_import" {
  federation_settings_id     = data.mongodbatlas_federated_settings.federated_settings.id
  org_id                     = var.org_id
  identity_provider_id       = var.identity_provider_id
  domain_restriction_enabled = false
  domain_allow_list          = ["yourdomain.com"]
  post_auth_role_grants      = ["ORG_MEMBER"]
}

resource "mongodbatlas_federated_settings_identity_provider" "identity_provider" {
  federation_settings_id       = data.mongodbatlas_federated_settings.federated_settings.id
  name                         = var.name
  associated_domains           = ["yourdomain.com"]
  sso_debug_enabled            = true
  status                       = "ACTIVE"
  sso_url                      = "https://mysso.oktapreview.com/app/mysso_terrafssotesdev_1/exk1f7f7f7fk5wp50h8/sso/saml"
  issuer_uri                   = "http://www.okta.com/exk1f716hf7f750h8"
  request_binding              = "HTTP-POST"
  response_signature_algorithm = "SHA-256"
}
