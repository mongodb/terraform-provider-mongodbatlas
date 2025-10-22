variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
}
variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  type        = string
}

variable "kafka_username" {
  description = "Username for connecting to your Kafka cluster"
  type        = string
}

variable "kafka_password" {
  description = "Password for connecting to your Kafka cluster"
  type        = string
}

variable "kafka_client_secret" {
  description = "Secret known only to the Kafka client and the authorization server"
  type        = string
}

variable "kafka_ssl_cert" {
  description = "Public certificate used for SASL_SSL configuration to connect to your Kafka cluster"
  type        = string
}

variable "cluster_name" {
  description = "Name of an existing cluster in your project that will be used to create a stream connection"
  type        = string
}

variable "other_project_id" {
  description = "Unique 24-hexadecimal digit string that identifies another project with a cluster that can be connected"
  type        = string
}

variable "other_cluster" {
  description = "Name of an existing cluster in another project within an organization that will be used to create a stream connection"
  type        = string
}