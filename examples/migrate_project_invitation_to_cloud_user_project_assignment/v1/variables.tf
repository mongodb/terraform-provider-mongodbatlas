variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "username" {
  description = "MongoDB Atlas Username (email) for pending invitation"
  type        = string
}

variable "roles" {
  description = "Project roles to assign to the user"
  type        = list(string)
  default     = ["GROUP_READ_ONLY"]
}
