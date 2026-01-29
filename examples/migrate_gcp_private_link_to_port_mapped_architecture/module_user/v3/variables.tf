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

variable "port_mapped_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture (used as forwarding rule name and address name)"
  type        = string
  default     = "tf-test-port-mapped-endpoint"
}

variable "port_mapped_address_ip" {
  description = "IP address for the port-mapped Google Compute Address"
  type        = string
  default     = "10.0.42.100"
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
