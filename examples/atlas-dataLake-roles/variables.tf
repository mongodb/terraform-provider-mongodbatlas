variable "public_key" {
  description = "The public API key for MongoDB Atlas"
  default     = ""
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
  default     = ""
}
variable "project_name" {
  description = "Atlas project name"
  default     = ""
}
variable "org_id" {
  description = "Atlas organization id"
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
variable "customer_master_key" {
  description = "The customer master secret key for AWS Account"
  default     = ""
}
variable "atlas_region" {
  default     = "US_EAST_1"
  description = "Atlas Region"
}
variable "aws_region" {
  default     = "us-east-1"
  description = "AWS Region"
}
variable "aws_iam_role_arn" {
  description = "AWS IAM ROLE ARN"
  default     = ""
}
