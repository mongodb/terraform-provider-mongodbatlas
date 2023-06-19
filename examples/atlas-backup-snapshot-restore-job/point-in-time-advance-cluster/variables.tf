variable "project_name" {
  description = "Atlas project name"
  default     = "ProjectTest0"
  type        = string
}
variable "org_id" {
  description = "The organization ID"
  type        = string
}
variable "cluster_name" {
  description = "Cluster name"
  default     = "ClusterTest0"
  type        = string
}
variable "point_in_time_utc_seconds" {
  description = "Point in time timestamp for snapshot_restore_job"
  default     = 0
  type        = number
}

variable "retain_backups_enabled" {
  description = "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster"
  default     = true
  type        = bool
}
