variable "project_id" {
  description = "Atlas project ID of the backup source cluster."
  type        = string
}

variable "cluster_name" {
  description = "Name of the backup source cluster."
  type        = string
}

variable "snapshot_id" {
  description = "Cloud Backup snapshot ID to inspect. When unset, uses the latest completed snapshot for the cluster."
  type        = string
  default     = null
}

variable "available_snapshots_limit" {
  description = "Maximum number of newest completed snapshots to include in the available_snapshots output."
  type        = number
  default     = 5
}
