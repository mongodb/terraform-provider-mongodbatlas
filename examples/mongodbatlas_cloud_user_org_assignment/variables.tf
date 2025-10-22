variable "org_id" {
  description = "The MongoDB Atlas organization ID"
  type        = string
}

variable "user_email" {
  description = "The email address of the user"
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
  default     = ""
}
