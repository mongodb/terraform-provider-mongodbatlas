variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  default     = ""
}

variable "org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}

variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "team_id" {
  description = "MongoDB Atlas Team ID"
  type        = string
}

variable "user_id" {
  description = "MongoDB Atlas User ID"
  type        = string
}

variable "username" {
  description = "MongoDB Atlas Username (email)"
  type        = string
}
