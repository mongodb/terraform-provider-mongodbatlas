provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

module "user_info" {
  source     = "../../module_maintainer/v3"
  username   = var.username
  org_id     = var.org_id
  project_id = var.project_id
}

output "user_id" {
  value = module.user_info.user_id
}

output "username" {
  value = module.user_info.username
}

output "email_address" {
  value = module.user_info.email_address
}

output "first_name" {
  value = module.user_info.first_name
}

output "last_name" {
  value = module.user_info.last_name
}

output "org_roles" {
  value = module.user_info.org_roles
}

output "project_roles" {
  value = module.user_info.project_roles
}
