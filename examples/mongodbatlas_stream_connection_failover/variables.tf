variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}
variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  type        = string
}
variable "bootstrap_servers" {
  description = "Comma separated list of server addresses for the primary Kafka cluster"
  type        = string
}
variable "kafka_username" {
  description = "Username for the primary Kafka cluster"
  type        = string
}
variable "kafka_password" {
  description = "Password for the primary Kafka cluster"
  type        = string
  sensitive   = true
}
variable "failover_bootstrap_servers" {
  description = "Comma separated list of server addresses for the failover-region Kafka cluster"
  type        = string
}
variable "failover_kafka_username" {
  description = "Username for the failover-region Kafka cluster"
  type        = string
}
variable "failover_kafka_password" {
  description = "Password for the failover-region Kafka cluster"
  type        = string
  sensitive   = true
}
