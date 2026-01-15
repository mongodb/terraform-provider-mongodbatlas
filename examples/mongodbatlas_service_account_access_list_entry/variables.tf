variable "org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}

variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID (for provider authentication)"
  type        = string
  default     = ""
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret (for provider authentication)"
  type        = string
  sensitive   = true
  default     = ""
}
