variable "public_key" { type = string }
variable "private_key" { type = string }

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
