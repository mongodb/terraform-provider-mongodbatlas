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

// azure
variable client_id {
  default     = ""
}
variable subscription_id {
  default     = ""
}
variable resource_group_name {
  default     = ""
}
variable client_secret {
  default     = ""
}
variable tenant_id {
  default     = ""
}
variable key_vault_name {
  default     = ""
}
variable key_identifier {
  default     = ""
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
