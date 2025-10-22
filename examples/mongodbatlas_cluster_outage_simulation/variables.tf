variable "atlas_client_id" {
  type        = string
  description = "MongoDB Atlas Service Account Client ID"
}
variable "atlas_client_secret" {
  type        = string
  description = "MongoDB Atlas Service Account Client Secret"
}
variable "project_id" {
  type        = string
  description = "MongoDB Project ID"
}
variable "atlas_cluster_name" {
  type        = string
  description = "Atlas Cluster Name that will undergo outage simulation"
  default     = "Cluster0"
}
variable "atlas_cluster_type" {
  type        = string
  description = "Atlas Cluster Type"
  default     = "REPLICASET"
}
variable "provider_instance_size_name" {
  type        = string
  description = "Cluster tier. Default is M10"
  default     = "M10"
}
variable "provider_name" {
  type    = string
  default = "AWS"
}
