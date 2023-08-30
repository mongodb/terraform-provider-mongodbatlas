resource "mongodbatlas_project" "project1" {
  name   = "Atlas-DB-Scope"
  org_id = var.org_id
}
output "project_name" {
  value = mongodbatlas_project.project1.name
}
