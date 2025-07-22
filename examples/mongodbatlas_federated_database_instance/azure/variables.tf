variable "public_key" {
  type        = string
  description = "Public Programmatic API key to authenticate to Atlas"
}
variable "private_key" {
  type        = string
  description = "Private Programmatic API key to authenticate to Atlas"
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