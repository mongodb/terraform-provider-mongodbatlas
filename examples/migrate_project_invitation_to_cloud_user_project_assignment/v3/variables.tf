variable "atlas_client_id" {
  type    = string
  description = "MongoDB Atlas Service Account Client ID"
  default = ""
}
variable "atlas_client_secret" {
  type    = string
  description = "MongoDB Atlas Service Account Client Secret"
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
