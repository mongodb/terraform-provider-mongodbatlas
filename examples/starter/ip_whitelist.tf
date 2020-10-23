resource "mongodbatlas_project_ip_whitelist" "ip_whitelist" {
  project_id = mongodbatlas_project.project.id
  cidr_block = "77.107.233.162/32"
  comment    = "cidr block for accessing the cluster"
}
output "ipwhitelist" {
  value = mongodbatlas_project_ip_whitelist.ip_whitelist.cidr_block
}
