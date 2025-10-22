variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "team_name" {
  type        = string
  description = "Name of the Atlas team"
}

variable "usernames" {
  type        = list(string)
  description = "List of user emails to assign to the team"
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

