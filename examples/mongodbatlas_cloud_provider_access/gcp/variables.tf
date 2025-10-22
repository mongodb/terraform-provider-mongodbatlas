variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
}

variable "atlas_project_id" {
  description = "MongoDB Atlas project ID where the cloud provider access will be configured"
  type        = string
}
