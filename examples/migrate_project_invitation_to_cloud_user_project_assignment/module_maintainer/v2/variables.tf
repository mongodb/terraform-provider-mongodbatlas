variable "project_id" {
  type        = string
  description = "MongoDB Atlas Project ID"
}

variable "username" {
  type        = string
  description = "Email/username of the Atlas user to invite"
}

variable "roles" {
  type        = list(string)
  description = "Project roles for the user"
}
