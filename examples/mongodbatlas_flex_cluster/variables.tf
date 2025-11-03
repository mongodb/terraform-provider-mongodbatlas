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
  description = "Atlas Project ID"
  type        = string
}
variable "cluster_name" {
  description = "Atlas cluster name"
  type        = string
  default     = "string"
}

variable "snapshot_id" {
  description = "Atlas snapshot ID"
  type        = string
}

variable "restore_job_id" {
  description = "Atlas restore job ID"
  type        = string
}