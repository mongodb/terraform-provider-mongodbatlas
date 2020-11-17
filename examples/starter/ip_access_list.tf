resource "mongodbatlas_project_ip_access_list" "ip" {
  project_id = mongodbatlas_project.project.id
  ip_address = "77.107.233.162"
  comment    = "cidr block for accessing the cluster"
}
output "ipaccesslist" {
  value = mongodbatlas_project_ip_access_list.ip.ip_address
}
