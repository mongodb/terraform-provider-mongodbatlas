provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
  base_url      = var.atlas_base_url
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs = [
    {
      region_configs = [
        {
          # Primary region: 2 nodes ensure the primary stays in this region
          # during normal operation.
          electable_specs = {
            instance_size = "M10"
            node_count    = 2
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        },
        {
          # Secondary region: 2 nodes provide a full failover target.
          electable_specs = {
            instance_size = "M10"
            node_count    = 2
          }
          provider_name = "AWS"
          priority      = 6
          region_name   = "US_WEST_2"
        },
        {
          # Tiebreaker region: 1 node provides the odd vote that prevents a
          # split vote if the primary and secondary regions are simultaneously
          # degraded. A smaller or geographically distant region is acceptable
          # here since it will rarely become primary.
          electable_specs = {
            instance_size = "M10"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 5
          region_name   = "EU_WEST_1"
        }
      ]
    }
  ]

  termination_protection_enabled = false

  tags = {
    environment = "production"
  }
}

resource "mongodbatlas_project" "project" {
  name   = var.project_name
  org_id = var.atlas_org_id
}
