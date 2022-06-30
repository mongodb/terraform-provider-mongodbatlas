variable "public_key" {
  type        = string
  description = "Public Programmatic API key to authenticate to Atlas"
}
variable "private_key" {
  type        = string
  description = "Private Programmatic API key to authenticate to Atlas"
}
variable "org_id" {
  type        = string
  description = "MongoDB Organization ID"
}
variable "group_id" {
  type        = string
  description = "MongoDB Group ID"
}

variable "name" {
  type        = string
  description = "MongoDB Identity Provider Name"
  default     = "mongodb_federation_test"
}

variable "identity_provider_id" {
  type        = string
  description = "MongoDB Identity Provider ID"
  default     = "5754gdhgd758"
}
