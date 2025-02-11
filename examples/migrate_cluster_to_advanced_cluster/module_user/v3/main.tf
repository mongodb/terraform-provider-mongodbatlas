module "cluster" {
  source = "../../module_maintainer/v3"

  cluster_name           = var.cluster_name
  cluster_type           = var.cluster_type
  mongo_db_major_version = var.mongo_db_major_version
  project_id             = var.project_id
  replication_specs_new  = var.replication_specs_new
  tags                   = var.tags
}
