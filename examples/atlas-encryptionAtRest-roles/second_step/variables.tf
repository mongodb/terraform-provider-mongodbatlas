variable "public_key" {
  description = "The public API key for MongoDB Atlas"
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
}
variable "project_id" {
  description = "Atlas project ID"
}
variable "customer_master_key" {
  description = "The customer master secret key for AWS Account"
}
variable "atlas_region" {
  default     = "US_EAST_1"
  description = "Atlas Region"
}

variable "cpa_role_id" {
  description = "AWS IAM ROLE ARN"
  default = ""
}
