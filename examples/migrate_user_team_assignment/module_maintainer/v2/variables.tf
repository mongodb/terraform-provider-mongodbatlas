variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "team_name" {
  type        = string
  description = "Name of the Atlas team"
}

variable "team_assigments" {
  type = map(object({
    org_id  = string
    team_id = string
    user_id = string
  }))
}
