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
