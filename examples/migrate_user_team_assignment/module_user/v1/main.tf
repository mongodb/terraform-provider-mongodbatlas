provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

# Old module usage
module "user_team_assignment" {
  source    = "../../module_maintainer/v1"
  org_id    = var.org_id
  team_name = var.team_name
  usernames = var.usernames
}
