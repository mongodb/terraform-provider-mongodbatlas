variable "project_id" {
  description = "Atlas project id"
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

variable "fcv_expiration_date" {
  description = "Expiration date of the pinned FCV"
  type        = string
}
