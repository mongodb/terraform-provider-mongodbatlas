variable "aws_secret_id" {
  description = "ARN of the AWS Secrets Manager secret containing the app SA credentials."
  type        = string
}

variable "aws_region" {
  description = "AWS region where the Secrets Manager secret is stored."
  type        = string
  default     = "us-east-1"
}

variable "project_id" {
  description = "Atlas project ID created by the platform phase."
  type        = string
}

variable "cluster_name" {
  description = "Name for the flex cluster."
  type        = string
  default     = "app-cluster"
}
