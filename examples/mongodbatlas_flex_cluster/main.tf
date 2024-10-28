resource "mongodbatlas_flex_cluster" "flex_cluster" {
  project_id = var.project_id
  name       = "clusterName"
  provider_settings = {
    backing_provider_name = "AWS"
    region_name           = "US_EAST_1"
  }
  termination_protection_enabled = true
} # TODO: add plural and singular data source example
