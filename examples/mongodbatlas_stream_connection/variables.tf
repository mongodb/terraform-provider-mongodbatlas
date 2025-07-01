variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
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

variable "s3_test_bucket_name" {
  description = "Name of the S3 bucket to use for testing AWS S3 connections"
  type        = string
}
