resource "mongodbatlas_cluster_adaptive_settings" "this" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  adaptive_settings_overrides = jsonencode({
    OVERLOAD_PROTECTION        = var.overload_protection_enabled
    SEARCH_OVERLOAD_PROTECTION = var.search_overload_protection_enabled
  })
}

data "mongodbatlas_cluster_adaptive_settings" "this" {
  project_id   = mongodbatlas_cluster_adaptive_settings.this.project_id
  cluster_name = mongodbatlas_cluster_adaptive_settings.this.cluster_name
}
