variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "username" {
  type        = string
  description = "Email/username of the Atlas user to invite"
}

variable "roles" {
  type        = list(string)
  description = "Org roles for the user"
}

variable "team_ids" {
  type        = set(string)
  description = "Team IDs to include on the invitation)"
  default     = []
}

