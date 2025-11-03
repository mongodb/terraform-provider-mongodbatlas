resource "mongodbatlas_advanced_cluster" "aws_private_connection" {
  project_id     = var.project_id
  name           = var.cluster_name
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
  depends_on = [mongodbatlas_privatelink_endpoint_service.pe_east_service]
}
