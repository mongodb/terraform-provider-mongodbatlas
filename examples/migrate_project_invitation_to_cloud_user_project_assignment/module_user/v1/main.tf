provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

module "project_membership" {
  source     = "../../module_maintainer/v1"
  project_id = var.project_id
  username   = var.username
  roles      = var.roles
}
