# Manage Adaptive Settings overrides for an existing MongoDB Atlas cluster and read the resulting effective settings with the corresponding data source.
resource "mongodbatlas_cluster_adaptive_settings" "this" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  adaptive_settings_overrides = jsonencode({
    LOAD_SHEDDING              = var.load_shedding_enabled
    SEARCH_OVERLOAD_PROTECTION = var.search_load_shedding_enabled
  })
}

data "mongodbatlas_cluster_adaptive_settings" "this" {
  project_id   = mongodbatlas_cluster_adaptive_settings.this.project_id
  cluster_name = mongodbatlas_cluster_adaptive_settings.this.cluster_name
}
