variable "atlas_client_id" {
  description = "The MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}

variable "atlas_client_secret" {
  description = "The MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "access_key" {
  description = "The access key for AWS Account"
  type        = string
}

variable "secret_key" {
  description = "The secret key for AWS Account"
  type        = string
  sensitive   = true
}

variable "aws_region" {
  description = "The AWS Region"
  default     = "us-east-1"
  type        = string
}

variable "atlas_org_id" {
  description = "The Atlas Organization ID"
  type        = string
}

variable "atlas_project_name" {
  description = "The Name of the Atlas project"
  default     = "tf-log-integration"
  type        = string
}

variable "s3_bucket_name" {
  description = "The name of the S3 bucket to which Atlas will send the logs"
  default     = "atlas-log-integration"
  type        = string
}

variable "aws_iam_role_name" {
  description = "The name of the IAM role to use to set up cloud provider access in Atlas"
  default     = "atlas-log-integration-role"
  type        = string
}

variable "log_types" {
  description = "The MongoDB log type to create"
  type        = string array
}

variable "prefix_path" {
  description = "The prefix to be added to the S3 Bucket path"
  type        = string
}