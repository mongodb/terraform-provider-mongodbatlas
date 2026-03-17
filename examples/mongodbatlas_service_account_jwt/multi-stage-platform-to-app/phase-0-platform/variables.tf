variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID (org-level)."
  type        = string
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret (org-level)."
  type        = string
  sensitive   = true
}

variable "org_id" {
  description = "MongoDB Atlas Organization ID."
  type        = string
}

variable "project_name" {
  description = "Name for the Atlas project created for the app team."
  type        = string
  default     = "app-team-project"
}
