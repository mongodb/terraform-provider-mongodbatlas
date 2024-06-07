// Global Vars
variable "project_name" {
  type        = string
  default     = "tf-example-oidc-Azure"
  description = "Both for the resource group in Azure and project name in MongoDB Atlas"
}

variable "owner" {
  type        = string
  default     = "apix-espen"
  description = "Used in tags.Owner for all resources supporting tags"
}

# Azure Vars
variable "ssh_public_key" {
  type        = string
  description = "Azure VM instance supports connection with ssh, see README.md for help. Tip: `TF_VAR_ssh_public_key=$(cat ~/.ssh/id_rsa.pub)`"
}

variable "location" {
  type        = string
  default     = "eastus"
  description = "Azure location, e.g., eastus"
}
variable "vm_admin_username" {
  type        = string
  default     = "adminuser"
  description = "Username used for the Azure VM instance"
}

variable "token_audience" {
  type        = string
  default     = "https://management.azure.com/"
  description = "Used as `resource` when getting the access token. See more in the [Azure documentation](https://learn.microsoft.com/en-us/entra/identity/managed-identities-azure-resources/how-to-use-vm-token#get-a-token-using-http)"
}

# MongoDB Atlas vars

variable "region" {
  type        = string
  description = "MongoDB Atlas Cluster Region, e.g., US_EAST_1"
}

variable "org_id" {
  type        = string
  description = "MongoDB Organization ID"
}

# Insert record vars

variable "insert_record_fields" {
  type = map(string)
  default = {
    "hello" = "world"
  }
}
variable "insert_record_database" {
  type    = string
  default = "test"
}

variable "insert_record_collection" {
  type    = string
  default = "test"
}

