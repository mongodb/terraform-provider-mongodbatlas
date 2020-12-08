resource "mongodbatlas_cluster" "cluster" {
  project_id             = mongodbatlas_project.project.id
  name                   = "mongodb-atlas"
  mongo_db_major_version = var.mongodbversion
  cluster_type           = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = var.region
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  //Provider Settings "block"
  provider_backup_enabled      = true
  auto_scaling_disk_gb_enabled = true
  provider_name                = "AWS"
  disk_size_gb                 = 10
  provider_disk_iops           = 100
  provider_volume_type         = "STANDARD"
  provider_instance_size_name  = "M10"
  provider_encrypt_ebs_volume  = true
}
output "atlasclusterstring" {
  value = mongodbatlas_cluster.cluster.connection_strings
}
