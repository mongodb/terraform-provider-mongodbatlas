resource "mongodbatlas_serverless_instance" "cluster_atlas" {
  project_id                              = var.atlasprojectid
  name                                    = "cluster-atlas"
  provider_settings_backing_provider_name = "AWS"
  provider_settings_provider_name         = "SERVERLESS"
  provider_settings_region_name           = "US_EAST_1"
  continuous_backup_enabled               = true
}

data "mongodbatlas_serverless_instance" "cluster_atlas" {
  project_id = var.atlasprojectid
  name       = mongodbatlas_serverless_instance.cluster_atlas.name
  depends_on = [mongodbatlas_privatelink_endpoint_service_serverless.atlaseplink]
}


output "atlasclusterstring" {
  value = data.mongodbatlas_serverless_instance.cluster_atlas.connection_strings_standard_srv
}

/* Note Value not available until second apply*/
/*
output "plstring" {
 value = mongodbatlas_serverless_instance.cluster_atlas.connection_strings_private_endpoint_srv[0]
}
*/
