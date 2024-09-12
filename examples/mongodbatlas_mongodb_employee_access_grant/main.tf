resource "mongodbatlas_mongodb_employee_access_grant" "example" {
  project_id      = var.project_id
  cluster_name    = var.cluster_name
  grant_type      = "CLUSTER_INFRASTRUCTURE_AND_APP_SERVICES_SYNC_DATA"
  expiration_time = "2025-01-01T12:00:00Z"
}


data "mongodbatlas_mongodb_employee_access_grant" "ds_example" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
}

output "grant_type" {
  value = data.mongodbatlas_mongodb_employee_access_grant.ds_example.grant_type
}

output "expiration_time" {
  value = data.mongodbatlas_mongodb_employee_access_grant.ds_example.expiration_time
}
