provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# cluster_profile PROTOTYPE — CORE (baseline)
#
# CORE keeps today's behavior. Because no `auto_scaling` block is set below, NO
# auto-scaling defaults are applied: the cluster stays a fixed M30. Omitting
# `cluster_profile` entirely behaves identically to setting it to "CORE".
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id      = mongodbatlas_project.project.id
  name            = var.cluster_name
  cluster_type    = "REPLICASET"
  cluster_profile = "CORE" # baseline; could also be omitted for the same effect

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          # CORE applies no auto-scaling defaults -> the cluster stays M30.
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        }
      ]
    }
  ]
}

resource "mongodbatlas_project" "project" {
  name   = var.project_name
  org_id = var.org_id
}
