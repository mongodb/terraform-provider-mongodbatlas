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
  default     = "us-central1"
}

variable "network_name" {
  description = "Name of the Google Compute Network"
  type        = string
  default     = "my-network"
}

variable "subnet_name" {
  description = "Name of the Google Compute Subnetwork"
  type        = string
  default     = "my-subnet"
}

variable "subnet_ip_cidr_range" {
  description = "IP CIDR range for the subnet"
  type        = string
  default     = "10.0.0.0/16"
}

variable "legacy_endpoint_count" {
  description = "Number of endpoints for legacy architecture (defaults to 50, matches Atlas project's privateServiceConnectionsPerRegionGroup setting)"
  type        = number
  default     = 50
}

variable "endpoint_base_name" {
  description = "Base name for endpoint resources"
  type        = string
  default     = "tf-test-legacy"
}

variable "endpoint_base_ip" {
  description = "Base IP address for endpoints (e.g., '10.0.42')"
  type        = string
  default     = "10.0.42"
}

variable "legacy_endpoint_service_id" {
  description = "Endpoint service ID for legacy architecture (can be any identifier string)"
  type        = string
  default     = "legacy-endpoint-group"
}

variable "cluster_name" {
  description = "(Optional) Cluster whose connection string to output"
  type        = string
  default     = ""
}
