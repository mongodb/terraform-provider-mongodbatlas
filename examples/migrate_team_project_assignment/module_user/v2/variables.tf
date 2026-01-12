variable "org_id" {
  description = "The ID of the MongoDB Atlas organization"
  type        = string
}

variable "project_name" {
  description = "Name of the project to create"
  type        = string
}

variable "team_map" {
  description = "Map of team_id to role_names"
  type        = map(list(string))
  default     = {}
}

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
