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
  description = "Atlas project ID"
  type        = string
}
variable "cluster_name" {
  description = "Atlas cluster name"
  default     = "aws-cross-region-private-connection"
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
variable "aws_region_east" {
  default     = "us-east-1"
  description = "AWS Region East (primary endpoint service region)"
  type        = string
}
variable "aws_region_west" {
  default     = "us-west-2"
  description = "AWS Region West (remote region for cross-region connectivity)"
  type        = string
}
