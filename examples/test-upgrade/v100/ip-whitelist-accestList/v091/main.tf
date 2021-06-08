resource "mongodbatlas_project" "test" {
  name   = var.project_name
  org_id = var.org_id
}


resource "mongodbatlas_project_ip_whitelist" "test" {
  project_id = mongodbatlas_project.test.id
  ip_address = var.ip_address
  comment    = var.comment
}
output "project_id" {
  value = mongodbatlas_project.test.id
}
output "entry" {
  value = mongodbatlas_project_ip_whitelist.test.ip_address
}
