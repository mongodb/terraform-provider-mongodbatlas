variable "atlas_client_id" {
  type        = string
  description = "Atlas Service Account Client ID"
}

variable "atlas_client_secret" {
  type        = string
  description = "Atlas Service Account Client Secret"
  sensitive   = true
}

variable "username" {
  type        = string
  description = "MongoDB Atlas Username (email)"
}

variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "project_id" {
  type        = string
  description = "MongoDB Atlas Project ID"
  default     = null
}
