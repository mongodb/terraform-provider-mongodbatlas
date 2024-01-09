provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

resource "mongodbatlas_cluster" "cluster" {
  project_id                  = mongodbatlas_project.project.id
  name                        = "ClusterToUpgrade"
  cluster_type                = "REPLICASET"
  provider_name               = var.provider_name
  backing_provider_name       = var.backing_provider_name
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = var.provider_instance_size_name
}

resource "mongodbatlas_project" "project" {
  name   = "TenantUpgradeTest"
  org_id = var.atlas_org_id
}
