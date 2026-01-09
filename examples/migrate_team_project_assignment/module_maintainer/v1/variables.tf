variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "project_name" {
  type        = string
  description = "Name of the project to create"
}

variable "team_map" {
  type        = map(list(string))
  description = "Map of team_id to role_names"
  default     = {}
}
