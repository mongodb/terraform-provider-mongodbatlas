// mongo 
variable "project_id" {
  default = ""
}
variable "cloud_provider_access_name" {
  default = "AWS"
}
variable "public_key" {
  default = ""
}
variable "private_key" {
  default = ""
}

// gcp
variable "service_account_key" {
  default = ""
}
variable "gcp_key_version_resource_id" {
  default = ""
}

// encryption at rest

variable "atlas_region" {
  default     = "US_EAST_1"
  description = "Atlas Region"
}

variable "project_name" {
  description = "Atlas project name"
  default     = ""
}
variable "org_id" {
  description = "The organization ID"
  default     = ""
}
