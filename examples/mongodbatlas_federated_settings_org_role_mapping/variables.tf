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
  description = "MongoDB Group/Project ID"
}

variable "name" {
  type        = string
  description = "MongoDB Identity Provider Name"
  default     = "mongodb_federation_test"
}
