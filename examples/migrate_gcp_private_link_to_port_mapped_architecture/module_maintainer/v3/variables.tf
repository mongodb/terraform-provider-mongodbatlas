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

variable "port_mapped_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture (used as forwarding rule name and address name)"
  type        = string
}

variable "port_mapped_address_ip" {
  description = "IP address for the port-mapped Google Compute Address"
  type        = string
}

variable "cluster_name" {
  description = "(Optional) Cluster whose connection string to output"
  type        = string
  default     = ""
}
