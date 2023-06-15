# mongo

variable "cloud_provider_access_name" {
  default = "AWS"
  type    = string
}

# aws
variable "access_key" {
  type = string
}
variable "secret_key" {
  type = string
}

# encryption at rest
variable "customer_master_key" {
  description = "The customer master secret key for AWS Account"
  type        = string
}
variable "atlas_region" {
  default     = "US_EAST_1"
  description = "Atlas Region"
  type        = string
}
variable "project_name" {
  description = "Atlas project name"
  type        = string
}
variable "org_id" {
  description = "The organization ID"
  type        = string
}
