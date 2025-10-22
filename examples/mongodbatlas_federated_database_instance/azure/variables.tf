variable "atlas_client_id" {
  type        = string
  description = "MongoDB Atlas Service Account Client ID"
}
variable "atlas_client_secret" {
  type        = string
  description = "MongoDB Atlas Service Account Client Secret"
}

variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "azure_atlas_app_id" {
  description = "Azure Atlas Application ID"
  type        = string
}

variable "azure_service_principal_id" {
  description = "Azure Service Principal ID"
  type        = string
}

variable "azure_tenant_id" {
  description = "Azure Tenant ID"
  type        = string
}

variable "federated_instance_name" {
  description = "Name for the Federated Database Instance"
  type        = string
  default     = "azure-federated-instance"
}
