variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}
variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  type        = string
}
variable "bootstrap_servers" {
  description = "Comma separated list of server addresses for the primary-region Amazon MSK cluster"
  type        = string
}
variable "failover_bootstrap_servers" {
  description = "Comma separated list of server addresses for the failover-region Amazon MSK cluster"
  type        = string
}
variable "kafka_iam_role_arn" {
  description = "ARN of the AWS IAM role that MongoDB Cloud assumes to authenticate to the Amazon MSK clusters (AWS_MSK_IAM)"
  type        = string
}
