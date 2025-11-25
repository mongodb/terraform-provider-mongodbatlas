# azure variables
variable "azure_region" {
  description = "The Azure region where resources will be created."
  type        = string
}

variable "azure_resource_group" {
  description = "Name for the Azure resource group."
  type        = string
}

variable "vnet_name" {
  description = "Name for the Azure Virtual Network."
  type        = string
}

variable "subnet_name" {
  description = "Name for the Azure Subnet that will host the Private Endpoint."
  type        = string
}

variable "eventhub_namespace_name" {
  description = "Globally unique name for the Azure Event Hubs Namespace."
  type        = string
}

variable "eventhub_name" {
  description = "Name for the Azure Event Hub within the namespace."
  type        = string
}

variable "vnet_address_space" {
  description = "The address space for the Azure Virtual Network."
  type        = list(string)
}

variable "subnet_address_prefix" {
  description = "The address prefix for the Azure Subnet."
  type        = list(string)
}

# MongoDB Atlas variables

variable "project_id" {
  description = "The ID of the MongoDB Atlas project."
  type        = string
}

variable "atlas_region" {
  description = "The atlas region of the Provider’s cluster. See [AZURE](https://www.mongodb.com/docs/atlas/reference/microsoft-azure/#stream-processing-instances)"
  type        = string
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