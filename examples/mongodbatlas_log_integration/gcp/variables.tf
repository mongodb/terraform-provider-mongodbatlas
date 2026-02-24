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