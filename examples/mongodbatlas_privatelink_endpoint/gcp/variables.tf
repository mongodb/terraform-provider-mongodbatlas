variable "gcp_project_id" {
  default = "GCP-PROJECT"
  type    = string
}
variable "gcp_region" {
  default = "us-central1"
  type    = string
}
variable "project_id" {
  default = "PROJECT-ID"
  type    = string
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
}
