# mongo
variable "project_id" {
  type = string
}
variable "cloud_provider_access_name" {
  type    = string
  default = "AWS"
}
variable "atlas_client_id" {
  type = string
  description = "MongoDB Atlas Service Account Client ID"
}
variable "atlas_client_secret" {
  type = string
  description = "MongoDB Atlas Service Account Client Secret"
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
