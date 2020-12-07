resource "mongodbatlas_cluster" "cluster-atlas" {
  project_id                   = var.atlasprojectid
  name                         = "cluster-atlas"
  provider_backup_enabled      = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "4.2"
  cluster_type                 = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = var.atlas_region
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  //Provider settings
  provider_name               = "AWS"
  disk_size_gb                = 10
  provider_disk_iops          = 100
  provider_volume_type        = "STANDARD"
  provider_encrypt_ebs_volume = true
  provider_instance_size_name = "M10"
}
output "atlasclusterstring" {
  value = mongodbatlas_cluster.cluster-atlas.connection_strings
}
output "plstring" {
  value = lookup(mongodbatlas_cluster.cluster-atlas.connection_strings[0].aws_private_link_srv, aws_vpc_endpoint.ptfe_service.id)
}
