variable "public_key" {
  type    = string
  default = ""
}
variable "private_key" {
  type    = string
  default = ""
}

variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "username" {
  description = "MongoDB Atlas Username (email) for user assignment"
  type        = string
}

variable "roles" {
  description = "Project roles to assign to the user"
  type        = list(string)
  default     = ["GROUP_READ_ONLY"]
}
