variable "atlas_client_id" {
  type = string
  description = "MongoDB Atlas Service Account Client ID"
}
variable "atlas_client_secret" {
  type = string
  description = "MongoDB Atlas Service Account Client Secret"
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

