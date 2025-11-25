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
  type = string
}

variable "cluster_name" {
  type = string
}

variable "database_name" {
  type = string
}

variable "collection_name" {
  type = string
}

variable "partition_field_one" {
  type = string
}

variable "partition_field_two" {
  type = string
}

