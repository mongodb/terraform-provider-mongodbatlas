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
  description = "MongoDB Atlas project ID that holds the cluster"
  type        = string
}

variable "cluster_name" {
  description = "MongoDB Atlas cluster with the collection to index"
  type        = string
}

variable "database" {
  description = "Database that holds the collection to index"
  type        = string
}

variable "collection_name" {
  description = "Collection to index."
  type        = string
}
