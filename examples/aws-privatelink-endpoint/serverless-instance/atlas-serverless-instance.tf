resource "mongodbatlas_serverless_instance" "aws_private_connection" {
  project_id                              = var.project_id
  name                                    = var.instance_name
  provider_settings_backing_provider_name = "AWS"
  provider_settings_provider_name         = "SERVERLESS"
  provider_settings_region_name           = "US_EAST_1"
  continuous_backup_enabled               = true
}