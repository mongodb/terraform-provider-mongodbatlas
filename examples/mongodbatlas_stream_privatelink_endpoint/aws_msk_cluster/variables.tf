variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  type        = string
}

variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
}

variable "aws_account_id" {
  description = "The AWS Account ID (12 digits)"
  type        = string
}

variable "msk_cluster_name" {
  description = "The MSK cluster's desired name"
  type        = string
}

variable "aws_secret_arn" {
  description = "AWS Secrets Manager secret ARN. Must meet the criteria outlined in https://docs.aws.amazon.com/msk/latest/developerguide/msk-password-tutorial.html"
  type        = string
}
