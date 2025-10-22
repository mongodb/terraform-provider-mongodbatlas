variable "project_id" {
  description = "Atlas project id"
  type        = string
}
variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
}

variable "fcv_expiration_date" {
  description = "Expiration date of the pinned FCV"
  type        = string
}
