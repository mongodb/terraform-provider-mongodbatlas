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

variable "org_id" { type = string }

# Pending invite user (now managed via cloud_user_org_assignment)
variable "pending_username" { type = string }
variable "roles" {
  type    = set(string)
  default = ["ORG_MEMBER"]
}
variable "pending_team_ids" {
  type    = set(string)
  default = []
}

# Active user already in org
variable "active_username" { type = string }
