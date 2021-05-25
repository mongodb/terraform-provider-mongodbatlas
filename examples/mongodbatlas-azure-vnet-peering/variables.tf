variable "public_key" {
}
variable "private_key" {
}
variable "project_id" {
}
variable "provider_instance_size_name" {
}
variable "provider_disk_type_name" {
}
variable "resource_group_name" {
}
variable "vnet_name" {
}
variable "atlas_cidr_block" {
  default = "192.168.248.0/21"
}
variable "location" {
  description = "The Azure region"
}
variable "provider_region_name" {
  description = "The Atlas region name"
}
variable "name" {
  description = "Atlas cluster name"
}
variable "address_space" {
  description = "Azure VNET CIDR"
}
variable "application_id" {
  default = "e90a1407-55c3-432d-9cb1-3638900a9d22"
}

