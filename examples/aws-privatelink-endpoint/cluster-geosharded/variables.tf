variable "public_key" {
  description = "The public API key for MongoDB Atlas"
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
}
variable "project_id" {
  description = "Atlas project ID"
}
variable "cluster_name" {
  description = "Atlas cluster name"
  default     = "geosharded"
}
variable "access_key" {
  description = "The access key for AWS Account"
}
variable "secret_key" {
  description = "The secret key for AWS Account"
}
variable "atlas_region_east" {
  default     = "US_EAST_1"
  description = "Atlas Region East"
}
variable "atlas_region_west" {
  default     = "US_WEST_1"
  description = "Atlas Region West"
}
variable "aws_region_east" {
  default     = "us-east-1"
  description = "AWS Region East"
}

variable "aws_region_west" {
  default     = "us-west-1"
  description = "AWS Region West"
}
