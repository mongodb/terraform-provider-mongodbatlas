variable "atlas_client_id" {
  description = "Atlas Service Account Client ID"
  type        = string
  default     = ""
}

variable "atlas_client_secret" {
  description = "Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "atlas_org_id" {
  description = "Atlas Organization ID"
  type        = string
}

variable "atlas_project_name" {
  description = "Name of the Atlas project"
  type        = string
  default     = "tf-log-integration-azure"
}

variable "azure_subscription_id" {
  description = "Azure ID that identifies your Azure subscription"
  type        = string
}

variable "azure_client_id" {
  description = "Azure ID identifies an Azure application associated with your Azure Active Directory tenant"
  type        = string
}

variable "azure_client_secret" {
  description = "Secret associated to the Azure application"
  type        = string
  sensitive   = true
}

variable "azure_tenant_id" {
  description = "Azure ID that identifies the Azure Active Directory tenant within your Azure subscription"
  type        = string
}

variable "atlas_azure_app_id" {
  description = "Azure Active Directory Application ID of Atlas"
  type        = string
}

variable "azure_service_principal_id" {
  description = "UUID identifying the Azure Service Principal representing Atlas"
  type        = string
}

variable "azure_resource_group_name" {
  description = "Name of the Azure resource group for log storage"
  type        = string
}

variable "azure_resource_group_location" {
  description = "Azure region where the resource group will be created"
  type        = string
  default     = "East US"
}

variable "azure_storage_account_name" {
  description = "Name of the Azure storage account"
  type        = string
}

variable "azure_storage_container_name" {
  description = "Name of the Azure storage container"
  type        = string
}
