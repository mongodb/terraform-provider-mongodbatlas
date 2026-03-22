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
          electable_specs = {
            instance_size = "M10"
            node_count    = 3
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
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
