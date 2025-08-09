resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = var.cloud_provider
      region_name   = var.region
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
}

output "connection_strings" {
  value = mongodbatlas_advanced_cluster.cluster.connection_strings.standard_srv
}
