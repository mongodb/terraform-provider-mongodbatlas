provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
  base_url      = "https://cloud-dev.mongodb.com/"
}

# Old module usage (deprecated teams block)
module "project" {
  source       = "../../module_maintainer/v1"
  org_id       = var.org_id
  project_name = var.project_name
  team_map     = var.team_map
}
