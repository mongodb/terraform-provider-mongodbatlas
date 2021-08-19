// mongo 
variable project_id {
  default     = ""
}
variable cloud_provider_access_name {
    default = "AWS"
}
variable public_key {
  default     = ""
}
variable private_key {
  default     = ""
}

// aws
variable access_key {
  default     = ""
}
variable secret_key {
  default     = ""
}
variable aws_region {
  default     = ""
}

// encryption at rest
variable "customer_master_key" {
  description = "The customer master secret key for AWS Account"
  default     = ""
}

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
