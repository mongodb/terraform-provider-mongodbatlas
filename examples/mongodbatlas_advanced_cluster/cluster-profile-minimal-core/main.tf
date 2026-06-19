provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# MINIMAL-CONFIG PROTOTYPE — CORE (cluster_profile omitted)
#
# Only project_id + name + provider_region are set. cluster_profile is omitted entirely, which
# behaves as CORE (the baseline) — an unset profile is treated the same as "CORE". The profile
# fills in the remaining previously-required inputs during `terraform plan`, resolving to a full cluster:
#   cluster_type      = "REPLICASET"            # static default
#   replication_specs = one AWS:US_EAST_1 shard with:
#     electable_specs.instance_size = "M10"     # CORE default base tier (entry-level dedicated)
#     electable_specs.node_count    = 3         # standard 3-node replica set
#     priority                      = 7
#     auto_scaling                  = (known after apply) — CORE injects NO auto-scaling defaults
resource "mongodbatlas_advanced_cluster" "example" {
  project_id      = "000000000000000000000000" # replace with your project_id
  name            = "my-cluster"
  provider_region = "AWS:US_EAST_1"
}
