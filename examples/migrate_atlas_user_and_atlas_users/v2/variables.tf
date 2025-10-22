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
