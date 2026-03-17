variable "atlas_client_id" {
  description = "Atlas Service Account Client ID"
  type        = string
  default     = ""
}

variable "atlas_client_secret" {
  description = "Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "atlas_org_id" {
  description = "Atlas Organization ID"
  type        = string
}

variable "atlas_project_name" {
  description = "Name of the Atlas project"
  type        = string
  default     = "tf-log-integration-gcp"
}

variable "gcp_project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "gcs_bucket_name" {
  description = "Name of the GCS bucket for storing Atlas logs (must be globally unique)"
  type        = string
  default     = "atlas-log-integration"
}

variable "gcs_bucket_location" {
  description = "Location where the GCS bucket is created"
  type        = string
  default     = "US"
}
