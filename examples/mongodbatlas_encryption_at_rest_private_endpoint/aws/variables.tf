variable "atlas_public_key" {
  description = "The public API key for MongoDB Atlas"
  type        = string
}
variable "atlas_private_key" {
  description = "The private API key for MongoDB Atlas"
  type        = string
  sensitive   = true
}
variable "atlas_project_id" {
  description = "Atlas Project ID"
  type        = string
}

variable "atlas_aws_region" {
  type        = string
  description = "Region in which the Encryption At Rest private endpoint is located."
}

variable "aws_kms_key_id" {
  type        = string
  description = "Region in which the Encryption At Rest private endpoint is located."
}

variable "access_key" {
  description = "The access key for AWS Account"
  type        = string
}
variable "secret_key" {
  description = "The secret key for AWS Account"
  type        = string
}
variable "aws_region" {
  type        = string
  description = "AWS Region"
}
