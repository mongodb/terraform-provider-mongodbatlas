# mongo
variable "project_id" {
  type    = string
}
variable "cloud_provider_access_name" {
  type    = string
  default = "AZURE"
}
variable "public_key" {
  type    = string
}
variable "private_key" {
  type    = string
}

variable "azure_tenant_id" {
  type    = string
}

variable "atlas_azure_app_id" {
  type = string
}

