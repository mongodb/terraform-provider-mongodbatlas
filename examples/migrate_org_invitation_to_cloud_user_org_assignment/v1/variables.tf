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

variable "org_id" { type = string }
# Pending invite user
variable "pending_username" { type = string }
variable "roles" {
  type    = set(string)
  default = ["ORG_MEMBER"]
}
# Teams for pending invite
variable "pending_team_ids" {
  type    = set(string)
  default = []
}

# Active user already in org (no invitation resource remains in state)
variable "active_username" { type = string }
