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

variable "name" {
  type        = string
  description = "MongoDB DataLake Pipeline Name"
  default     = "mongodb_federation_test"
}

variable "cluster_name" {
  type        = string
  description = "Cluster Name"
}
