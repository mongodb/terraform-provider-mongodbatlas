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

variable "new_endpoint_service_id" {
  description = "Endpoint service ID for port-mapped architecture (used as forwarding rule name and address name)"
  type        = string
  default     = "tf-test-port-mapped-endpoint"
}
