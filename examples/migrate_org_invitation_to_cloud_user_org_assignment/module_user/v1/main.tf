provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# Old module usage
module "org_membership" {
  source   = "../../module_maintainer/v1"
  org_id   = var.org_id
  username = var.username
  roles    = var.roles
  team_ids = var.team_ids
}

