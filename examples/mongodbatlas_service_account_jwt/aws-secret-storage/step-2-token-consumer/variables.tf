variable "aws_region" {
  description = "AWS region where the Secrets Manager secret is stored."
  type        = string
  default     = "us-east-1"
}

variable "aws_secret_id" {
  description = "ARN of the AWS Secrets Manager secret containing the Atlas JWT (from step1 output)."
  type        = string
}

variable "org_id" {
  description = "MongoDB Atlas Organization ID."
  type        = string
}

variable "project_name" {
  description = "Name for the Atlas project created using the JWT."
  type        = string
  default     = "jwt-example-project"
}
