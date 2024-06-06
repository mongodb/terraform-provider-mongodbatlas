data "mongodbatlas_control_plane_ip_addresses" "test" {
}

output "outbound-aws-ip-addresses" {
  value = data.mongodbatlas_control_plane_ip_addresses.test.outbound.aws
}
