provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

# Old module usage
module "user_team_assignment" {
  source    = "../../module_maintainer/v1"
  org_id    = var.org_id
  team_name = var.team_name
  usernames = var.usernames
}
