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

variable "org_id" {
  description = "Atlas Organization ID where the project is created"
  type        = string
}

variable "database" {
  description = "Database that holds the collection to index"
  type        = string
}

variable "collection_name" {
  description = "Collection to index. Must contain the fields referenced by each index below"
  type        = string
}
