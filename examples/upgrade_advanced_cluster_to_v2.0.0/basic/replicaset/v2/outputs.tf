output "connection_string_standard" {
  description = "Public connection string that you can use to connect to this cluster."
  # connection_strings is an object (no [0] index required anymomre)
  value = mongodbatlas_advanced_cluster.this.connection_strings.standard
}
