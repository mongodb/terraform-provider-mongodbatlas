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

variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "federated_instance_name" {
  description = "Name for the Federated Database Instance"
  type        = string
  default     = "gcp-federated-instance"
}

variable "gcp_project_id" {
  description = "Google Cloud project ID that hosts the Google Cloud Storage bucket"
  type        = string
}

variable "gcs_bucket_name" {
  description = "Name of the Google Cloud Storage bucket the Federated Database Instance queries. Must be globally unique."
  type        = string
}

variable "gcs_bucket_location" {
  description = "Location where the Google Cloud Storage bucket is created"
  type        = string
  default     = "us-central1"
}

variable "data_process_region" {
  description = "Atlas Data Federation region where the Federated Database Instance processes queries. Google Cloud supports IOWA_USA and BELGIUM_EU."
  type        = string
  default     = "IOWA_USA"
}
