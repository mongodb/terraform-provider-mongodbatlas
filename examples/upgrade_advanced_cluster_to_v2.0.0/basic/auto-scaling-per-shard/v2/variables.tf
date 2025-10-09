variable "project_id" {
  description = "Atlas project id"
  type        = string
}

variable "cluster_name" {
  description = "Atlas cluster name"
  type        = string
  default     = "AutoScalingCluster"
}
