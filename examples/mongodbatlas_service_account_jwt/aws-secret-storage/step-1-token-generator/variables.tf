variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID."
  type        = string
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret."
  type        = string
  sensitive   = true
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
