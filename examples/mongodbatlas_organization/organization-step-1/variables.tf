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
variable "operations_contact" {
  type        = string
  description = "Email address, typically a distribution list, for the organization to receive proactive notifications about its infrastructure"
  default     = null
}
variable "absolute_session_timeout_in_seconds" {
  type        = number
  description = "Maximum session duration in seconds for the Atlas UI. Accepted values range between 3,600 (1 hour) and 43,200 (12 hours)"
  default     = null
}
variable "idle_session_timeout_in_seconds" {
  type        = number
  description = "Maximum idle session duration in seconds for the Atlas UI. Accepted values start at 300 (5 minutes)"
  default     = null
}




