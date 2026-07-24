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

variable "azure_service_principal_id" {
  description = "UUID that identifies the Azure Service Principal used to access the Azure Blob Storage account"
  type        = string
}

variable "azure_storage_account_name" {
  description = "Name of the Azure Storage account to use for the Azure Blob Storage connection"
  type        = string
}

variable "azure_region" {
  description = "Azure region where you locate the storage account"
  type        = string
  default     = null
}

variable "schema_registry_username" {
  description = "Username for connecting to your Schema Registry"
  type        = string
  default     = ""
}

variable "schema_registry_password" {
  description = "Password for connecting to your Schema Registry"
  type        = string
  sensitive   = true
  default     = ""
}

variable "kafka_iam_role_arn" {
  description = "ARN of the AWS IAM role that MongoDB Cloud assumes to authenticate to an Amazon MSK cluster (AWS_MSK_IAM)"
  type        = string
  default     = ""
}

variable "kafka_ssl_certificate" {
  description = "SSL certificate for client authentication to Kafka (mutual TLS)"
  type        = string
  default     = ""
}

variable "kafka_ssl_key" {
  description = "SSL key for client authentication to Kafka (mutual TLS)"
  type        = string
  sensitive   = true
  default     = ""
}

variable "kafka_ssl_key_password" {
  description = "Password for the SSL key, if it is password protected"
  type        = string
  sensitive   = true
  default     = ""
}
