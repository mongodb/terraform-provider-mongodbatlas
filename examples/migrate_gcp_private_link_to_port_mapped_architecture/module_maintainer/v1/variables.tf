variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "gcp_project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "gcp_region" {
  description = "GCP Region"
  type        = string
}

variable "network_name" {
  description = "Name for the Google Compute Network"
  type        = string
}

variable "subnet_name" {
  description = "Name for the Google Compute Subnetwork"
  type        = string
}

variable "subnet_ip_cidr_range" {
  description = "IP CIDR range for the Google Compute Subnetwork"
  type        = string
}

variable "legacy_endpoint_count" {
  description = "Number of endpoints for legacy architecture (defaults to 50, matches Atlas project's privateServiceConnectionsPerRegionGroup setting)"
  type        = number
}

variable "legacy_address_name_prefix" {
  description = "Prefix for Google Compute Address names (legacy architecture)"
  type        = string
}

variable "legacy_address_base_ip" {
  description = "Base IP address for Google Compute Addresses (e.g., '10.0.42' will create addresses like 10.0.42.0, 10.0.42.1, etc.)"
  type        = string
}

variable "legacy_endpoint_service_id" {
  description = "Endpoint service ID for legacy architecture (can be any identifier string)"
  type        = string
}

variable "cluster_name" {
  description = "(Optional) Cluster whose connection string to output"
  type        = string
  default     = ""
}
