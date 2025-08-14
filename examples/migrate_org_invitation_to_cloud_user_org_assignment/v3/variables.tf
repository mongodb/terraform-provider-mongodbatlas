variable "public_key" { type = string }
variable "private_key" { type = string }

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
