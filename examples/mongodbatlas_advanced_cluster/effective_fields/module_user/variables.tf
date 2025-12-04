variable "atlas_org_id" {
  description = "Atlas Organization ID"
  type        = string
}

variable "project_name" {
  description = "Atlas Project Name"
  type        = string
  default     = "EffectiveFieldsExample"
}

variable "cluster_name" {
  description = "Atlas Cluster Name"
  type        = string
  default     = "example-cluster"
}
