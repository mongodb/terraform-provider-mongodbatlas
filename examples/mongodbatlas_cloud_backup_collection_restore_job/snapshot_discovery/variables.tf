variable "project_id" {
  description = "Atlas project ID of the source cluster"
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label of the source cluster"
  type        = string
}

variable "snapshot_id" {
  description = "Snapshot ID to browse. Defaults to the latest completed snapshot when unset."
  type        = string
  default     = null
}

variable "discovery_database_names" {
  description = "Database names to list collections from. Leave empty to list only databases."
  type        = list(string)
  default     = []
}
