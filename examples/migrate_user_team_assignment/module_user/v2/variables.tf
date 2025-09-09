variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "team_name" {
  type        = string
  description = "Name of the Atlas team"
}

variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
  default     = ""
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
  default     = ""
}

