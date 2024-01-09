
variable "project_id" {
  default = "PROJECT-ID"
  type    = string
}
variable "subscription_id" {
  default = "AZURE SUBSCRIPTION ID"
  type    = string
}
variable "client_id" {
  default = "AZURE CLIENT ID"
  type    = string
}
variable "client_secret" {
  default = "AZURE CLIENT SECRET"
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
  description = "Cluster whose connection string to output"
  default     = "cluster-serverless"
  type        = string
}
