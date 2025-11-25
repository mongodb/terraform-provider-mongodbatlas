# mongo
variable "project_id" {
  type = string
}
variable "cloud_provider_access_name" {
  type    = string
  default = "AWS"
}
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

# aws
variable "access_key" {
  type = string
}
variable "secret_key" {
  type = string
}
variable "aws_region" {
  type = string
}
