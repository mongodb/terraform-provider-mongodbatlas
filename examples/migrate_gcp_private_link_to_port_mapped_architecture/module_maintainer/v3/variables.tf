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

variable "new_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture (used as forwarding rule name and address name)"
  type        = string
  default     = "tf-test-port-mapped-endpoint"
}

variable "port_mapped_endpoint_ip" {
  description = "IP address for port-mapped endpoint"
  type        = string
  default     = "10.0.42.100"
}

variable "cluster_name" {
  description = "(Optional) Cluster whose connection string to output"
  type        = string
  default     = ""
}
