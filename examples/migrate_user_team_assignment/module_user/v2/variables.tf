variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "team_name" {
  type        = string
  description = "Name of the Atlas team"
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

