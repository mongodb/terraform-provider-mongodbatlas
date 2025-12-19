variable "atlas_org_id" {
  description = "Atlas Organization ID"
  type        = string
}

variable "atlas_project_name" {
  description = "Name of the Atlas project"
  default     = "tf-log-migration"
  type        = string
}

variable "aws_region" {
  description = "AWS Region"
  default     = "us-east-1"
  type        = string
}

variable "s3_bucket_name" {
  description = "The name of the S3 bucket to which Atlas will send the logs"
  default     = "atlas-log-export-migration"
  type        = string
}

variable "aws_iam_role_name" {
  description = "The name of the IAM role to use to set up cloud provider access in Atlas"
  default     = "atlas-log-export-role"
  type        = string
}

