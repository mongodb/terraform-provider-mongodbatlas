provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# MINIMAL-CONFIG PROTOTYPE — FULL EXPLICIT (reverse-compatibility check)
#
# An existing-style config that sets every (now-optional) input explicitly and uses NEITHER
# cluster_profile NOR provider_region. It must plan exactly as it does today: because the user
# specified everything, applyMinimalConfigDefaults is a no-op and no profile defaults are
# injected. This demonstrates that Required -> Optional did not change behavior for users who
# specify everything.
resource "mongodbatlas_advanced_cluster" "example" {
  project_id   = "000000000000000000000000" # replace with a real project_id to apply
  name         = "full-explicit"
  cluster_type = "REPLICASET"

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
}
