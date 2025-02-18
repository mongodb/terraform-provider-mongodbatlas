resource "mongodbatlas_advanced_cluster" "this" {
  project_id             = var.project_id
  name                   = var.cluster_name
  cluster_type           = "REPLICASET"
  mongo_db_major_version = var.mongo_db_major_version

  advanced_configuration = {
    javascript_enabled = true
  }
  replication_specs = [{
    region_configs = [{
      provider_name = "AWS"
      region_name   = "US_WEST_1"
      priority      = 7
      electable_specs = {
        node_count    = 2
        instance_size = var.instance_size
        disk_size_gb  = 30
      }
      }, {
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      priority      = 6
      electable_specs = {
        node_count    = 3
        instance_size = var.instance_size
        disk_size_gb  = 30
      }
      read_only_specs = {
        node_count    = 1
        instance_size = var.instance_size
        disk_size_gb  = 30
      }
      analytics_specs = {
        node_count    = 1
        instance_size = var.instance_size
        disk_size_gb  = 30
      }
    }]
  }]

  # Generated by atlas-cli-plugin-terraform.
  # Please confirm that all references to this resource are updated.
}
# moved {
#   from = mongodbatlas_cluster.this          # change `this` to your specific resource identifier
#   to   = mongodbatlas_advanced_cluster.this # change `this` to your specific resource identifier
# }
