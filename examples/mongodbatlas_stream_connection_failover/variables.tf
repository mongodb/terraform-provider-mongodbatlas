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

variable "workspace_name" {
  description = "Name of the stream workspace. The workspace must have failover regions enabled"
  type        = string
}

variable "kafka_username" {
  description = "Username for connecting to your Kafka cluster"
  type        = string
}

variable "kafka_password" {
  description = "Password for connecting to your Kafka cluster"
  type        = string
  sensitive   = true
}

variable "failover_region" {
  description = "Failover region for the failover connection. Must be one of the workspace's failover regions"
  type        = string
}

variable "failover_bootstrap_servers" {
  description = "Comma separated list of server addresses for the failover region Kafka cluster"
  type        = string
}
