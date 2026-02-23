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

variable "access_key" {
  description = "The access key for your Microsfot Azure Account"
  type        = string
}

variable "secret_key" {
  description = "The secret key for your Microsoft Azure Account"
  type        = string
  sensitive   = true
}

variable "azure_region" {
  description = "Azure Region"
  default     = "westus2"
  type        = string
}

variable "atlas_org_id" {
  description = "Atlas Organization ID"
  type        = string
}

variable "atlas_project_name" {
  description = "Name of the Atlas project"
  default     = "tf-log-integration"
  type        = string
}

variable "bucketname" {
  description = "The name of the Azure Blob to which Atlas will send the logs"
  default     = "my-log-bucket"
  type        = string
}

variable "iamRoleId" {
  description = "The name of the IAM role to use to set up cloud provider access in Atlas"
  default     = "atlas-log-integration-role"
  type        = string
}
