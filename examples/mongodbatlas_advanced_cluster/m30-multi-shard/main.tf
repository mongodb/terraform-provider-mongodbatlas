provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id           = mongodbatlas_project.project.id
  name                 = var.cluster_name
  cluster_type         = "SHARDED"
  backup_enabled       = true
  use_effective_fields = true

  replication_specs = [
    {
      # shard 1
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30" # Paid tier. Production-grade dedicated cluster
            node_count    = 3
            disk_size_gb  = 10
          }
          auto_scaling = {
            compute_enabled            = true
            compute_scale_down_enabled = true
            compute_min_instance_size  = "M30"
            compute_max_instance_size  = "M50"
            disk_gb_enabled            = true
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }
      ]
    },
    {
      # shard 2
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30" # Paid tier. Production-grade dedicated cluster
            node_count    = 3
            disk_size_gb  = 10
          }
          auto_scaling = {
            compute_enabled            = true
            compute_scale_down_enabled = true
            compute_min_instance_size  = "M30"
            compute_max_instance_size  = "M50"
            disk_gb_enabled            = true
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }
      ]
    }
  ]

  termination_protection_enabled = true

  tags = {
    environment = "production"
  }
}

resource "mongodbatlas_project" "project" {
  name   = var.project_name
  org_id = var.org_id
}
