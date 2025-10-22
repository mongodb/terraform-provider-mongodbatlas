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
variable "project_id" {
  description = "Atlas project ID"
  type        = string
}
variable "cluster_name" {
  description = "Atlas cluster name"
  default     = "aws-private-connection"
  type        = string
}
