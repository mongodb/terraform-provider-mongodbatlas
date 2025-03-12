output "connection_string_standard" {
  description = "Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb://` protocol."
  # value = mongodbatlas_cluster.this.connection_strings[0].standard # BEFORE
  value = mongodbatlas_advanced_cluster.this.connection_strings.standard # AFTER
}

# output "provider_name" { # BEFORE
#   value = mongodbatlas_advanced_cluster.this.provider_name # Unsupported attribute
# }

output "container_id" {
  description = "The Network Peering Container ID of the configuration specified in `region_configs`. The Container ID is the id of the container either created programmatically by the user before any clusters existed in a project or when the first cluster in the region (AWS/Azure) or project (GCP) was created. Example `AWS:US_EAST_1\" = \"61e0797dde08fb498ca11a71`"
  # value = mongodbatlas_cluster.this.container_id # BEFORE
  value = mongodbatlas_advanced_cluster.this.replication_specs[0].container_id["AWS:US_EAST_1"] # AFTER
}
