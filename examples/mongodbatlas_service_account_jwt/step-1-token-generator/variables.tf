variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID with permissions to create Service Accounts."
  type        = string
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret."
  type        = string
  sensitive   = true
}

variable "org_id" {
  description = "MongoDB Atlas Organization ID."
  type        = string
}

variable "service_account_name" {
  description = "Name for the Service Account created by this step."
  type        = string
  default     = "jwt-example-sa"
}

variable "aws_region" {
  description = "AWS region where the Secrets Manager secret is created."
  type        = string
  default     = "us-east-1"
}

variable "secret_name" {
  description = "AWS Secrets Manager secret name for the Atlas JWT."
  type        = string
  default     = "atlas/ephemeral-jwt"
}

variable "token_version" {
  description = "Increment to rotate the stored token on subsequent applies."
  type        = number
  default     = 1
}
