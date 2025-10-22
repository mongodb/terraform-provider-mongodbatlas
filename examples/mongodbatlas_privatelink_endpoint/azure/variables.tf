
variable "project_id" {
  default = "PROJECT-ID"
  type    = string
}
variable "subscription_id" {
  default = "AZURE SUBSCRIPTION ID"
  type    = string
}
variable "atlas_client_id" {
  default = "AZURE CLIENT ID"
  description = "MongoDB Atlas Service Account Client ID"
  type    = string
}
variable "atlas_client_secret" {
  default = "AZURE CLIENT SECRET"
  description = "MongoDB Atlas Service Account Client Secret"
  type    = string
}
variable "tenant_id" {
  default = "AZURE TENANT ID"
  type    = string
}
variable "resource_group_name" {
  default = "AZURE RESOURCE GROUP NAME"
  type    = string
}
variable "cluster_name" {
  description = "(Optional) Cluster whose connection string to output"
  type        = string
}
