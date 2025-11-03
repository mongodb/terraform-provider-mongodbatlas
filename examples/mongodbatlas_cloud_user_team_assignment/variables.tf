variable "org_id" {
  description = "The MongoDB Atlas organization ID"
  type        = string
}

variable "team_id" {
  description = "The team ID"
  type        = string
}

variable "user_id" {
  description = "The user ID"
  type        = string
}

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
