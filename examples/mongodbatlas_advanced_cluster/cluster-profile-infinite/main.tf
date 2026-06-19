provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# cluster_profile PROTOTYPE — INFINITE
#
# INFINITE injects compute auto-scaling defaults for every region config where you
# did NOT set an `auto_scaling` block. Here the configured instance size is M30 and
# no auto_scaling is set, so the provider automatically applies during plan:
#   compute_enabled            = true
#   compute_scale_down_enabled = true
#   compute_min_instance_size  = "M30"   # = configured instance size
#   compute_max_instance_size  = "M50"   # = two tiers above (M30 -> M40 -> M50)
#   disk_gb_enabled            = false
# Run `terraform plan` to see these values appear on auto_scaling automatically.
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id      = mongodbatlas_project.project.id
  name            = var.cluster_name
  cluster_type    = "REPLICASET"
  cluster_profile = "INFINITE"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M30"
            node_count    = 3
          }
          # INFINITE: auto-scaling defaults applied automatically — min=M30, max=M50.
          # To opt out for this region, set your own auto_scaling block here:
          # explicit input always wins over the profile default.
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
