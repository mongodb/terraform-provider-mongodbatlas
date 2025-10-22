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

variable "user_id" {
  description = "Atlas User ID"
  type        = string
}

variable "username" {
  description = "Atlas Username"
  type        = string
} 