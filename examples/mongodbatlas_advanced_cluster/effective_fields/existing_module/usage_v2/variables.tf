variable "atlas_org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}

variable "project_name" {
  description = "MongoDB Atlas Project Name"
  type        = string
  default     = "effective-fields-v2-example"
}

variable "cluster_name" {
  description = "MongoDB Atlas Cluster Name"
  type        = string
  default     = "test-cluster-v2"
}
