resource "mongodbatlas_advanced_cluster" "this" {
  project_id             = var.project_id
  name                   = var.cluster_name
  cluster_type           = "REPLICASET"
  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = var.instance_size
        node_count    = 3
      }
      provider_name = var.cloud_provider
      priority      = 7
      region_name   = var.region_name
    }]
  }]
}

resource "mongodbatlas_cluster_overload_simulation" "this" {
  project_id       = mongodbatlas_advanced_cluster.this.project_id
  cluster_name     = mongodbatlas_advanced_cluster.this.name
  duration_seconds = var.duration_seconds
}

data "mongodbatlas_cluster_overload_simulation" "this" {
  project_id    = mongodbatlas_cluster_overload_simulation.this.project_id
  cluster_name  = mongodbatlas_cluster_overload_simulation.this.cluster_name
  simulation_id = mongodbatlas_cluster_overload_simulation.this.simulation_id
}
