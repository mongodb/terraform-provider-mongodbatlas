provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# Final module usage (no import blocks needed after migration)
module "org_membership" {
  source   = "../../module_maintainer/v3"
  org_id   = var.org_id
  username = var.username
  roles    = var.roles
  team_ids = var.team_ids
}

