variable "gcp_project_id" {
  default = "GCP-PROJECT"
}
variable "gcp_region" {
  default = "us-central1"
}
variable "project_id" {
  default = "PROJECT-ID"
}
variable "public_key" {
  description = "Public API key to authenticate to Atlas"
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
}
variable "cluster_name" {
  description = "(Optional) Cluster whose connection string to output"
  default     = ""
}
