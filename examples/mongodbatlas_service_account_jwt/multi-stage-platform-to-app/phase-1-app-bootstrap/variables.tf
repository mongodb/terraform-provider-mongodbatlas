variable "access_token" {
  description = "Short-lived JWT from the platform phase."
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "Atlas project ID created by the platform phase."
  type        = string
}

variable "org_id" {
  description = "MongoDB Atlas Organization ID."
  type        = string
}

variable "aws_region" {
  description = "AWS region for the Secrets Manager secret."
  type        = string
  default     = "us-east-1"
}

variable "sa_name" {
  description = "Name for the app team's project-scoped Service Account."
  type        = string
  default     = "app-team-sa"
}

variable "secret_name" {
  description = "AWS Secrets Manager secret name for the app SA credentials."
  type        = string
  default     = "atlas/app-team-sa-creds"
}

variable "token_version" {
  description = "Increment to rotate the stored credentials on subsequent applies."
  type        = number
  default     = 1
}
