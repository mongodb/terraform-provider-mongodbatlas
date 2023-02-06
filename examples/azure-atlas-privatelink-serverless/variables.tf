
variable "project_id" {
  default = "PROJECT-ID"
}
variable "subscription_id" {
  default = "AZURE SUBSCRIPTION ID"
}
variable "client_id" {
  default = "AZURE CLIENT ID"
}
variable "client_secret" {
  default = "AZURE CLIENT SECRET"
}
variable "tenant_id" {
  default = "AZURE TENANT ID"
}
variable "resource_group_name" {
  default = "AZURE RESOURCE GROUP NAME"
}
variable "cluster_name" {
  description = "Cluster whose connection string to output"
  default     = "cluster-serverless"
}
