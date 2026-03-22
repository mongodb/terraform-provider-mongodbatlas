provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
  base_url      = var.atlas_base_url
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = mongodbatlas_project.project.id
  name         = var.cluster_name
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M0"
          }
          provider_name         = "TENANT"
          backing_provider_name = "AWS"
          region_name           = "US_EAST_1"
          priority              = 7
        }
      ]
    }
  ]

  # Termination protection is disabled by default for free-tier clusters used in
  # development and evaluation. Enable it before moving to production.
  termination_protection_enabled = false

  tags = {
    environment = "dev"
  }
}

resource "mongodbatlas_project" "project" {
  name   = var.project_name
  org_id = var.atlas_org_id
}
