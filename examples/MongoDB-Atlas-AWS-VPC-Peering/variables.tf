variable "public_key" {
  description = "The public API key for MongoDB Atlas"
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
}

variable "access_key" {
  description = "The access key for AWS Account"
}
variable "secret_key" {
  description = "The secret key for AWS Account"
}
variable "atlas_region" {
  default     = "US_EAST_1"
  description = "Atlas Region"
}
variable "aws_region" {
  default     = "ap-southeast-1"
  description = "AWS Region"
}
variable "atlas_dbuser" {
  description = "The db user for Atlas"
}
variable "atlas_dbpassword" {
  description = "The db user passwd for Atlas"
}
variable "aws_account_id" {
  description = "My AWS Account ID"
  default     = "208629369896"
}
variable "atlasorgid" {
  description = "Atlas Org ID"
  default     = "5c98a80fc56c98ef210b8633"
}
variable "atlas_vpc_cidr" {
  description = "Atlas CIDR"
  default     = "192.168.232.0/21"
}
