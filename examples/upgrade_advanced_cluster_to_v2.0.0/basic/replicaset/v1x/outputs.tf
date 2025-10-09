output "connection_string_standard" {
  description = "Public connection string that you can use to connect to this cluster."
  value       = mongodbatlas_advanced_cluster.this.connection_strings[0].standard
}
