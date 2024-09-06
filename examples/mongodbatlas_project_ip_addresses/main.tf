data "mongodbatlas_project_ip_addresses" "test" {
  project_id = var.project_id
}

output "project_services" {
  value = data.mongodbatlas_project_ip_addresses.test.services
}
