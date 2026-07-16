variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "Atlas project ID"
  type        = string
}

variable "azure_resource_group_name" {
  description = "Existing Azure resource group name."
  type        = string
}

variable "vnet_name" {
  description = "Existing Azure virtual network name."
  type        = string
}

variable "subnet_name" {
  description = "Existing Azure subnet name where the private endpoint is created."
  type        = string
}

variable "atlas_data_federation_private_link_service_resource_id" {
  description = "Azure Resource ID of the Atlas-managed Data Federation Private Link Service for your region. See https://www.mongodb.com/docs/atlas/data-federation/tutorial/config-private-endpoint/ to find the value for your region."
  type        = string
}
