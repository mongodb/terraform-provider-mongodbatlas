variable "public_key" {
  description = "The public API key for MongoDB Atlas"
  type        = string
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
  type        = string
}
variable "base_url" {
  type = string
}
variable "project_name" {
  description = "Atlas project name"
  type        = string
}
variable "org_id" {
  description = "Atlas organization id"
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
variable "aws_region" {
  default     = "us-east-1"
  description = "AWS Region"
  type        = string
}
variable "test_s3_bucket" {
  description = "The name of s3 bucket"
  type        = string
}
variable "data_lake_name" {
  description = "The data lake name"
  type        = string
}
