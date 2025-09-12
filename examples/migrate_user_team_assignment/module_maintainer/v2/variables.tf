variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "team_name" {
  type        = string
  description = "Name of the Atlas team"
}

variable "user_ids" {
  description = "Set of user IDs to assign to the team"
  type        = set(string)
}
