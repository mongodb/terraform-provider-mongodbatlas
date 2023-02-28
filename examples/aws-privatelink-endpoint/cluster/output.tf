locals {
  private_endpoints = flatten([for cs in mongodbatlas_cluster.aws_private_connection.connection_strings : cs.private_endpoint])

  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], aws_vpc_endpoint.vpce_east.id)
  ]
}

output "connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
