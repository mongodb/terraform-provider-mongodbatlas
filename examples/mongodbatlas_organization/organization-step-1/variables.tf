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
variable "org_owner_id" {
  type        = string
  description = "MongoDB Organization Owner ID"
}
variable "security_contact" {
  type        = string
  description = "Email address for the organization to receive security-related notifications"
}




