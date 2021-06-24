variable "public_key" {
  description = "The public API key for MongoDB Atlas"
  default     = ""
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
  default     = ""
}
variable "project_name" {
  description = "Atlas project name"
  default     = ""
}
variable "region_name" {
  description = "AWS region "
  default     = ""
}
variable "route_table_cidr_block" {
  description = "CIDR Block"
  default     = ""
}
variable "vpc_id" {
  description = "AWS VPC ID"
  default     = ""
}
variable "aws_account_id" {
  description = "AWS account id"
  default     = ""
}


