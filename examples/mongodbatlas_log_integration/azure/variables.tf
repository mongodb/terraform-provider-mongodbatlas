#MongoDB authentication variables
variable "atlas_client_id" {
  description = "The MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}

variable "atlas_client_secret" {
  description = "The MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "public_key" {
  description = "The Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "The Private API key to authenticate to Atlas"
  type        = string
}
variable "project_id" {
  description = "The MongoDB Project ID"
  type        = string
}
variable "log_types" {
  description = "The MongoDB log type to create"
  type        = list(string)
}

# Azure authentication variables
variable "azure_subscription_id" {
  description = "The Azure Subscription ID"
  type        = string
}
variable "azure_client_id" {
  description = "The Azure Principal Application (Client) ID"
  type        = string
}
variable "azure_client_secret" {
  description = "The Azure Service Principal Client Secret"
  type        = string
  sensitive   = true
}

# Azure variables
variable "atlas_azure_app_id" {
  description = "The Azure Active Directory Application ID of Atlas"
  type        = string
}
variable "azure_role_id" {
  description = "The UUID identifying the Azure Service Principal"
  type        = string
}
variable "azure_tenant_id" {
  description = "The Azure Active Directory Tenant ID"
  type        = string
}
variable "azure_resource_group_name" {
  description = "The Name of the Azure resource group for log storage"
  type        = string
}
variable "azure_storage_account_name" {
  description = "The Name of the Azure storage account"
  type        = string
}
variable "azure_storage_container_name" {
  description = "The Name of the Azure storage container"
  type        = string
}

variable "access_key" {
  description = "The access key for GCP Account"
  type        = string
}

variable "secret_key" {
  description = "The secret key for GCP Account"
  type        = string
  sensitive   = true
}

variable "azure_region" {
  description = "The name of the Azure container region"
  type        = string
}

