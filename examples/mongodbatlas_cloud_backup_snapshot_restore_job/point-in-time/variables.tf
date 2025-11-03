variable "project_name" {
  description = "Atlas project name"
  type        = string
}
variable "org_id" {
  description = "The organization ID"
  type        = string
}
variable "cluster_name" {
  description = "Cluster name"
  type        = string
}
variable "point_in_time_utc_seconds" {
  description = "Point in time timestamp for snapshot_restore_job"
  default     = 0
  type        = number
}
