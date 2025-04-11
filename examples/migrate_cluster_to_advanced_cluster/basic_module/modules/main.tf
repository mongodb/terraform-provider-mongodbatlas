resource "mongodbatlas_advanced_cluster" "this" {
  lifecycle {
    precondition {
      condition     = !(var.auto_scaling_disk_gb_enabled && var.disk_size > 0)
      error_message = "Must use either auto_scaling_disk_gb_enabled or disk_size, not both."
    }
  }

  project_id             = var.project_id
  name                   = var.cluster_name
  cluster_type           = var.cluster_type
  mongo_db_major_version = var.mongo_db_major_version
  replication_specs      = var.replication_specs
  tags                   = var.tags
}
