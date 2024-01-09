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
variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
}
variable "cluster_name" {
  description = "(Optional) Cluster whose connection string to output"
  type        = string
}
