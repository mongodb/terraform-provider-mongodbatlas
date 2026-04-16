variable "project_id" {
  description = "The MongoDB Atlas project ID"
  type        = string
}

variable "gcp_region" {
  description = "The GCP region where the Pub/Sub private link will be created"
  type        = string
  default     = "us-east1"
}
