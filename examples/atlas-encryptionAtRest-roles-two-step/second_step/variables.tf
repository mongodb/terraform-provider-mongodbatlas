variable "public_key" {
  description = "The public API key for MongoDB Atlas"
  default     = ""
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
  default     = ""
}
variable "project_id" {
  description = "Atlas project ID"
  default     = ""
}
variable "customer_master_key" {
  description = "The customer master secret key for AWS Account"
  default     = ""
}
variable "atlas_region" {
  default     = "US_EAST_1"
  description = "Atlas Region"
}

variable "cpa_role_id" {
  description = "AWS IAM ROLE ARN"
  default     = ""
}
