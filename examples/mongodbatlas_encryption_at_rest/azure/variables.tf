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
variable "atlas_project_id" {
  description = "Atlas Project ID"
  type        = string
}
variable "azure_subscription_id" {
  type        = string
  description = "Azure ID that identifies your Azure subscription"
}

variable "azure_resource_group_name" {
  type        = string
  description = "Name of the Azure resource group that contains your Azure Key Vault"
}

variable "azure_key_vault_name" {
  type        = string
  description = "Unique string that identifies the Azure Key Vault that contains your key"
}

variable "azure_key_identifier" {
  type        = string
  description = "Web address with a unique key that identifies for your Azure Key Vault"
}

variable "role_id" {
  type        = string
  description = "Azure ID that identifies the role with permissions to access your Azure Key Vault"
}


