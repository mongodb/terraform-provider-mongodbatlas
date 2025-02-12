output "connection_string_standard" {
  # value = mongodbatlas_cluster.this.connection_strings[0].standard # BEFORE
  value = mongodbatlas_advanced_cluster.this.connection_strings.standard # AFTER
}

# output "provider_name" { # BEFORE
#   value = mongodbatlas_advanced_cluster.this.provider_name # Unsupported attribute
# }

output "container_id" {
  # value = mongodbatlas_cluster.this.container_id # BEFORE
  value = mongodbatlas_advanced_cluster.this.replication_specs[0].container_id["AWS:US_EAST_1"] # AFTER
}
