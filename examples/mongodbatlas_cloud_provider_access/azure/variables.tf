# mongo
variable "project_id" {
  type = string
}
variable "cloud_provider_access_name" {
  type    = string
  default = "AZURE"
}
variable "atlas_client_id" {
  type = string
  description = "MongoDB Atlas Service Account Client ID"
}
variable "atlas_client_secret" {
  type = string
  description = "MongoDB Atlas Service Account Client Secret"
}

variable "azure_tenant_id" {
  type = string
}

variable "atlas_azure_app_id" {
  type = string
}

