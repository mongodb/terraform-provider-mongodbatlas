variable "public_key" {
  description = "The public API key for MongoDB Atlas"
  type        = string
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
  type        = string
}
variable "project_id" {
  description = "Atlas project ID"
  type        = string
}
variable "cluster_name" {
  description = "Atlas cluster name"
  default     = "geosharded"
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
variable "atlas_region_east" {
  default     = "US_EAST_1"
  description = "Atlas Region East"
  type        = string
}
variable "atlas_region_west" {
  default     = "US_WEST_1"
  description = "Atlas Region West"
  type        = string
}
variable "aws_region_east" {
  default     = "us-east-1"
  description = "AWS Region East"
  type        = string
}

variable "aws_region_west" {
  default     = "us-west-1"
  description = "AWS Region West"
  type        = string
}
