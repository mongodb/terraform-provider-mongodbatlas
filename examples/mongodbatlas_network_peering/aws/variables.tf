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

variable "access_key" {
  description = "The access key for AWS Account"
  type        = string
}
variable "secret_key" {
  description = "The secret key for AWS Account"
  type        = string
}
variable "atlas_region" {
  description = "Atlas Region"
  default     = "US_EAST_1"
  type        = string
}
variable "aws_region" {
  description = "AWS Region"
  default     = "ap-southeast-1"
  type        = string
}
variable "atlas_dbuser" {
  description = "The db user for Atlas"
  type        = string
}
variable "atlas_dbpassword" {
  description = "The db user passwd for Atlas"
  type        = string
}
variable "aws_account_id" {
  description = "My AWS Account ID"
  default     = "208629369896"
  type        = string
}
variable "atlas_org_id" {
  description = "Atlas Org ID"
  default     = "5c98a80fc56c98ef210b8633"
  type        = string
}
variable "atlas_vpc_cidr" {
  description = "Atlas CIDR"
  default     = "192.168.232.0/21"
  type        = string
}
