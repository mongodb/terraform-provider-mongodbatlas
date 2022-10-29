variable "project_name" {
  description = "Atlas project name"
  default     = ""
}
variable "org_id" {
  description = "The organization ID"
  default     = ""
}
variable "cluster_name" {
  description = "Cluster name"
  default     = ""
}
variable "point_in_time_utc_seconds" {
  description = "Point in time timestamp for snapshot_restore_job"
  default     = 0
}
