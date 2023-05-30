variable "public_key" {
  type        = string
  description = "Public Programmatic API key to authenticate to Atlas"
}
variable "private_key" {
  type        = string
  description = "Private Programmatic API key to authenticate to Atlas"
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
  description = "Atlas Cluster Name that will undergo outage simulation"
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
