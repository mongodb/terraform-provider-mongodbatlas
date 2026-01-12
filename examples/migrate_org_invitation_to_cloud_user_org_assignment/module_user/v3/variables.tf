variable "atlas_client_id" {
  type        = string
  description = "Atlas Service Account Client ID"
}

variable "atlas_client_secret" {
  type        = string
  description = "Atlas Service Account Client Secret"
}

variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "username" {
  type        = string
  description = "Email/username of the Atlas user to assign"
}

variable "roles" {
  type        = list(string)
  description = "Org roles for the user"
  default     = ["ORG_MEMBER"]
}

variable "team_ids" {
  type        = set(string)
  description = "Team IDs to assign to the user"
  default     = []
}

