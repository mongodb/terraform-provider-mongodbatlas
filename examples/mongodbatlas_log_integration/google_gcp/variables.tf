variable "atlas_client_id" {
  description = "The MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}

variable "atlas_client_secret" {
  description = "The MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "access_key" {
  description = "The access key for GCP Account"
  type        = string
}

variable "secret_key" {
  description = "The secret key for GCP Account"
  type        = string
  sensitive   = true
}

variable "GCP_region" {
  description = "The GCP Region"
  default     = "us-east1-b"
  type        = string


variable "public_key" {
  description = "The Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "The Private API key to authenticate to Atlas"
  type        = string
}
variable "project_id" {
  description = "The MongoDB Project ID"
  type        = string
}

variable "gcs_project_id" {
  description = "The GCS Project ID"
  type        = string
}
variable "gcs_bucket_name" {
  description = "The Name of the GCS bucket for log storage"
  type        = string
}