variable "atlas_client_id" {
  type        = string
  description = "MongoDB Atlas Service Account Client ID"
}
variable "atlas_client_secret" {
  type        = string
  description = "MongoDB Atlas Service Account Client Secret"
}
variable "org_id" {
  type        = string
  description = "MongoDB Organization ID"
}
variable "group_id" {
  type        = string
  description = "MongoDB Group/Project ID"
}

variable "name" {
  type        = string
  description = "MongoDB Identity Provider Name"
  default     = "mongodb_federation_test"
}
