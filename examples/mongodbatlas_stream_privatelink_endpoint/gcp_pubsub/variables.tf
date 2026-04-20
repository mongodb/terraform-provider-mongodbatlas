variable "project_id" {
  description = "The MongoDB Atlas project ID"
  type        = string
}

variable "cluster_name" {
  description = "The name of the GCP cluster to provision in the same region"
  type        = string
  default     = "gcp-cluster"
}

variable "gcp_region" {
  description = "The GCP region where the Pub/Sub private link will be created"
  type        = string
  default     = "us-east4"
}
