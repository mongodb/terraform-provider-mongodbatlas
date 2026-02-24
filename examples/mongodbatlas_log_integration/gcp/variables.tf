variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
}
variable "project_id" {
  description = "MongoDB Project ID"
  type        = string
}

variable "gcp_project_id" {
  description = "GCP Project ID"
  type        = string
}
variable "gcs_bucket_name" {
  description = "Name of the GCS bucket for log storage"
  type        = string
}