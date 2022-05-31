data "mongodbatlas_cloud_federated_settings" "federated_settings" {
  org_id = "627a968f7f7f7f76f14"
}

data "mongodbatlas_cloud_federated_settings_identity_providers" "identitty_provider" {
  federation_settings_id = data.mongodbatlas_cloud_federated_settings.federated_settings.id
}

data "mongodbatlas_cloud_federated_settings_org_configs" "org_configs_ds" {
  federation_settings_id = data.mongodbatlas_cloud_federated_settings.federated_settings.id
}

data "mongodbatlas_cloud_federated_settings_org_role_mappings" "org_connections" {
  federation_settings_id = data.mongodbatlas_cloud_federated_settings.federated_settings.id
  org_id                 = "627a968f7f7f7f76f14"
}
resource "mongodbatlas_cloud_federated_settings_org_role_mapping" "org_role_mapping" {
  federation_settings_id = data.mongodbatlas_cloud_federated_settings.federated_settings.id
  org_id                 = "627a968f7f7f7f76f14"

  external_group_name = "newgroup"

  organization_roles = ["ORG_OWNER", "ORG_MEMBER"]
  group_id           = "628aa2f7f7f7f7b098b81b8"
  group_roles        = ["GROUP_OWNER", "GROUP_CLUSTER_MANAGER", "GROUP_DATA_ACCESS_ADMIN", "GROUP_DATA_ACCESS_READ_WRITE", "GROUP_SEARCH_INDEX_EDITOR", "GROUP_DATA_ACCESS_READ_ONLY", "GROUP_READ_ONLY"]

}
resource "mongodbatlas_cloud_federated_settings_org_config" "org_connections_import" {
  federation_settings_id     = data.mongodbatlas_cloud_federated_settings.federated_settings.id
  org_id                     = "627a968f7f7f7f76f14"
  domain_restriction_enabled = false
  domain_allow_list          = ["yourdomain.com"]
}

resource "mongodbatlas_cloud_federated_settings_identity_provider" "identity_provider" {
  federation_settings_id = data.mongodbatlas_cloud_federated_settings.federated_settings.id
  name                   = "mongodb_federation_test"
  associated_domains     = ["yourdomain.com"]
  sso_debug_enabled      = true
  status                 = "ACTIVE"
}
