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
  type = string
}
variable "provider_instance_size_name" {
  type = string
}
variable "resource_group_name" {
  type = string
}
variable "vnet_name" {
  type = string
}
variable "atlas_cidr_block" {
  default = "192.168.248.0/21"
  type    = string
}
variable "location" {
  description = "The Azure region"
  type        = string
}
variable "provider_region_name" {
  description = "The Atlas region name"
  type        = string
}
variable "name" {
  description = "Atlas cluster name"
  type        = string
}
variable "address_space" {
  description = "Azure VNET CIDR"
  type        = string
}
variable "application_id" {
  default = "e90a1407-55c3-432d-9cb1-3638900a9d22"
  type    = string
}

