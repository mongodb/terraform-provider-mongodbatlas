provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

resource "mongodbatlas_cluster" "cluster" {
  project_id                  = mongodbatlas_project.project.id
  name                        = "NVMEToUpgrade"
  cluster_type                = "REPLICASET"
  provider_name               = var.provider_name
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = var.provider_instance_size_name
  provider_volume_type        = var.provider_volume_type
  provider_disk_iops          = var.provider_disk_iops
  cloud_backup                = true
}

resource "mongodbatlas_project" "project" {
  name   = "NVMEUpgradeTest"
  org_id = var.atlas_org_id
}
