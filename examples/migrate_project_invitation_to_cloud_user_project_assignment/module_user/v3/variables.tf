variable "atlas_client_id" {
  type        = string
  description = "Atlas Service Account Client ID"
}

variable "atlas_client_secret" {
  type        = string
  description = "Atlas Service Account Client Secret"
}

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
  default     = ["GROUP_READ_ONLY"]
}
