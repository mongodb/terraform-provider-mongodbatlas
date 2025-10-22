provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

module "cluster" {
  source = "../../module_maintainer/v1"

  auto_scaling_disk_gb_enabled = var.auto_scaling_disk_gb_enabled
  cluster_name                 = var.cluster_name
  cluster_type                 = var.cluster_type
  disk_size                    = var.disk_size
  instance_size                = var.instance_size
  mongo_db_major_version       = var.mongo_db_major_version
  project_id                   = var.project_id
  provider_name                = var.provider_name
  replication_specs            = var.replication_specs
  tags                         = var.tags
}

output "mongodb_connection_strings" {
  description = "Collection of Uniform Resource Locators that point to the MongoDB database."
  value       = module.cluster.mongodb_connection_strings
}
