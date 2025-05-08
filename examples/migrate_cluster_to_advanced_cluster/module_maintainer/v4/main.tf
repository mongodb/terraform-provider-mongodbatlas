resource "mongodbatlas_advanced_cluster" "this" {
  project_id             = var.project_id
  name                   = var.cluster_name
  cluster_type           = var.cluster_type
  mongo_db_major_version = var.mongo_db_major_version
  replication_specs      = var.replication_specs
  tags                   = var.tags
}
