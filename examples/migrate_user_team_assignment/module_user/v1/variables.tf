variable "org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}

variable "team_name" {
  description = "Name of the Atlas team"
  type        = string
}

variable "usernames" {
  description = "List of user emails to assign to the team"
  type        = list(string)
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

