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
variable "cpa_role_id" {
  description = "AWS IAM ROLE ARN"
  default     = ""
}
variable "access_key" {
  description = "The access key for AWS Account"
  default     = ""
}
variable "secret_key" {
  description = "The secret key for AWS Account"
  default     = ""
}
variable "test_s3_bucket" {
  description = "The name of s3 bucket"
  default     = ""
}
variable "data_lake_name" {
  description = "The data lake name"
  default     = ""
}
variable "data_lake_region" {
  default     = "VIRGINIA_USA"
  description = "The data lake region"
}