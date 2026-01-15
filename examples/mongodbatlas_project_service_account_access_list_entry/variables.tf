variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID (for provider authentication)"
  type        = string
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret (for provider authentication)"
  type        = string
  sensitive   = true
}
