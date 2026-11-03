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
  description = "Atlas Project ID"
  type        = string
}

variable "vpc_endpoint_id" {
  description = "AWS VPC interface endpoint ID, e.g. vpce-0abc123def4567890"
  type        = string
}

variable "vpc_endpoint_dns_name" {
  description = "DNS name of the AWS VPC interface endpoint"
  type        = string
}
