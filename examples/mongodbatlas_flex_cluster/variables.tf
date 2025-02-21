variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
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