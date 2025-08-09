# This cluster is in GCP cloud-provider with VPC peering enabled

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = var.project_id
  name           = "cluster-test"
  cluster_type   = "REPLICASET"
  backup_enabled = true # enable cloud provider snapshots

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "GCP"
      region_name   = var.atlas_region
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      auto_scaling = {
        compute_enabled            = true
        compute_scale_down_enabled = true
        compute_min_instance_size  = "M10"
        compute_max_instance_size  = "M20"
        disk_gb_enabled            = true
      }
    }]
  }]
  tags = {
    environment = "prod"
  }
  advanced_configuration = {
    minimum_enabled_tls_protocol = "TLS1_2"
  }

  lifecycle {
    ignore_changes = [
      replication_specs[0].region_configs[0].electable_specs.instance_size,
    ]
  }
}

# The connection strings available for the GCP MognoDB Atlas cluster
output "connection_string" {
  value = mongodbatlas_advanced_cluster.cluster.connection_strings
}
