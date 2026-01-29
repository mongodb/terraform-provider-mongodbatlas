variable "gcp_project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "gcp_region" {
  description = "GCP Region"
  type        = string
  default     = "us-central1"
}

variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

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

variable "cluster_name" {
  description = "(Optional) Cluster whose connection string to output"
  type        = string
  default     = ""
}

variable "legacy_endpoint_count" {
  description = "Number of endpoints for legacy architecture (defaults to 50, matches Atlas project's privateServiceConnectionsPerRegionGroup setting)"
  type        = number
  default     = 50
}

variable "legacy_endpoint_service_id" {
  description = "Endpoint service ID for legacy architecture (can be any identifier string)"
  type        = string
  default     = "legacy-endpoint-group"
}

variable "port_mapped_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture (used as forwarding rule name and address name)"
  type        = string
  default     = "tf-test-port-mapped-endpoint"
}

variable "network_name" {
  description = "Name for the Google Compute Network"
  type        = string
  default     = "my-network"
}

variable "subnet_name" {
  description = "Name for the Google Compute Subnetwork"
  type        = string
  default     = "my-subnet"
}

variable "subnet_ip_cidr_range" {
  description = "IP CIDR range for the Google Compute Subnetwork"
  type        = string
  default     = "10.0.0.0/16"
}

variable "legacy_address_name_prefix" {
  description = "Prefix for Google Compute Address names (legacy architecture)"
  type        = string
  default     = "tf-test-legacy"
}

variable "legacy_address_base_ip" {
  description = "Base IP address for Google Compute Addresses (e.g., '10.0.42' will create addresses like 10.0.42.0, 10.0.42.1, etc.)"
  type        = string
  default     = "10.0.42"
}

variable "port_mapped_address_ip" {
  description = "IP address for the port-mapped Google Compute Address"
  type        = string
  default     = "10.0.42.100"
}
