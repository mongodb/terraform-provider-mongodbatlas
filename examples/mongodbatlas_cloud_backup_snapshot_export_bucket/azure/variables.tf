variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
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
