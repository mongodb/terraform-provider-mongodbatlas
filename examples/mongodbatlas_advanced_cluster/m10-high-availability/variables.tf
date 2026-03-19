variable "atlas_base_url" {
  description = "MongoDB Atlas Base URL"
  type        = string
  default     = "https://cloud.mongodb.com/"
}
variable "atlas_org_id" {
  description = "Atlas organization id"
  type        = string
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
variable "project_name" {
  description = "Atlas project name"
  type        = string
  default     = "m10-ha-project"
}
variable "cluster_name" {
  description = "Atlas cluster name"
  type        = string
  default     = "m10-high-availability"
}
