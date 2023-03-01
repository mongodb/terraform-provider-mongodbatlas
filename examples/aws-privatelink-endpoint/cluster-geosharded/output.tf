locals {
  private_endpoints = flatten([for cs in mongodbatlas_cluster.geosharded.connection_strings : cs.private_endpoint])

  connection_strings_east = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], aws_vpc_endpoint.vpce_east.id)
  ]
  connection_strings_west = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], aws_vpc_endpoint.vpce_west.id)
  ]
}

output "connection_string_east" {
  value = length(local.connection_strings_east) > 0 ? local.connection_strings_east[0] : ""
}
output "connection_string_west" {
  value = length(local.connection_strings_west) > 0 ? local.connection_strings_west[0] : ""
}
