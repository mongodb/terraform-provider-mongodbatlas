variable "atlas_client_id" {
  type        = string
  description = "MongoDB Atlas Client ID (Service Account ID)"
}

variable "atlas_client_secret" {
  type        = string
  description = "MongoDB Atlas Client Secret (Service Account Secret)"
  sensitive   = true
}

variable "atlas_org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}
