locals {
  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = var.instance_size
        node_count    = 3
      }
      provider_name = var.provider_name
      region_name   = var.region_name
      priority      = 7
    }]
  }]
}

resource "mongodbatlas_advanced_cluster" "this" {
  project_id        = var.project_id
  name              = var.name
  cluster_type      = "REPLICASET"
  replication_specs = local.replication_specs
  tags              = var.tags
}
