variable "public_key" {
  description = "The public API key for MongoDB Atlas"
  type        = string
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
  type        = string
}
variable "project_id" {
  description = "Atlas project ID"
  type        = string
}
variable "azure_tenant_id" {
  type = string
}
variable "subscription_id" {
  default = "Azure Subscription ID"
  type    = string
}
variable "client_id" {
  default = "Azure Client ID"
  type    = string
}
variable "client_secret" {
  default = "Azure Client Secret"
  type    = string
}
variable "tenant_id" {
  default = "Azure Tenant ID"
  type    = string
}
variable "azure_atlas_app_id" {
  type = string
}
variable "azure_resource_group_location" {
  type = string
}
variable "storage_account_name" {
  type = string
}
