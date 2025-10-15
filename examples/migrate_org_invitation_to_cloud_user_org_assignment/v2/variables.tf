variable "public_key" {
  type    = string
  default = ""
}
variable "private_key" {
  type    = string
  default = ""
}

variable "org_id" { type = string }
variable "pending_username" { type = string }
variable "roles" {
  type    = set(string)
  default = ["ORG_MEMBER"]
}
variable "pending_team_ids" {
  type    = set(string)
  default = []
}
variable "active_username" { type = string }
