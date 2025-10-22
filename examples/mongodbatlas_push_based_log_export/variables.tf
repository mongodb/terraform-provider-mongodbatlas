variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
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
  description = "AWS Region"
  default     = "us-east-1"
  type        = string
}
variable "atlas_org_id" {
  description = "Atlas Organization ID"
  type        = string
}
variable "atlas_project_name" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  default     = "tf-push-based-log"
  type        = string
}
variable "s3_bucket_name" {
  description = "The name of the bucket to which Atlas will send the logs to"
  default     = "atlas-log-export"
  type        = string
}
variable "s3_bucket_policy_name" {
  description = "The name of the IAM role policy to configure for the S3 bucket"
  default     = "atlas-log-export-s3-policy"
  type        = string
}
variable "aws_iam_role_name" {
  description = "The name of the IAM role to use to set up cloud provider access in Atlas"
  default     = "atlas-log-export-role"
  type        = string
}
variable "aws_iam_role_policy_name" {
  description = "The name of the IAM role policy for the configured aws_iam_role_name"
  default     = "atlas-log-export-role-policy"
  type        = string
}

