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
