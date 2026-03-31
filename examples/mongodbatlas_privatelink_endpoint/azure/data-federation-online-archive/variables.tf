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

variable "project_id" {
  description = "Atlas project ID"
  type        = string
}

variable "atlas_region" {
  description = "Atlas region for the Azure private link endpoint (example: eastus2)."
  type        = string
  default     = "eastus2"
}

variable "azure_location" {
  description = "Azure location for the resource group and networking resources."
  type        = string
  default     = "East US 2"
}

variable "azure_resource_group_name" {
  description = "Azure resource group name for this example."
  type        = string
  default     = "mdb-atlas-df-oa-rg"
}

variable "vnet_name" {
  description = "Azure virtual network name for this example."
  type        = string
  default     = "mdb-atlas-df-oa-vnet"
}

variable "vnet_cidr" {
  description = "CIDR block for the Azure virtual network."
  type        = list(string)
  default     = ["10.0.0.0/16"]
}

variable "subnet_name" {
  description = "Azure subnet name where the private endpoint is created."
  type        = string
  default     = "mdb-atlas-df-oa-subnet"
}

variable "subnet_cidr" {
  description = "CIDR block for the Azure subnet."
  type        = list(string)
  default     = ["10.0.1.0/24"]
}
