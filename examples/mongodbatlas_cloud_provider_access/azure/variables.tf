# mongo
variable "project_id" {
  type = string
}
variable "cloud_provider_access_name" {
  type    = string
  default = "AZURE"
}
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

variable "azure_tenant_id" {
  type = string
}

variable "atlas_azure_app_id" {
  type = string
}

