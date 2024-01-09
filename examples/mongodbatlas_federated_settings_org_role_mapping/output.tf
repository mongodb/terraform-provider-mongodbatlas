output "federated_settings_ds" {
  value = data.mongodbatlas_federated_settings.federated_settings.id
}

output "identity_provider" {
  value = data.mongodbatlas_federated_settings_identity_providers.identity_provider.id
}

output "org_configs_ds" {
  value = data.mongodbatlas_federated_settings_org_configs.org_configs_ds.id
}

output "org_role_mapping" {
  value = data.mongodbatlas_federated_settings_org_role_mappings.org_role_mapping.id
}
