variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "cluster_name" {
  description = "Name of the cluster"
  type        = string
  default     = "example-cluster"
}

variable "instance_size" {
  description = "Instance size for the cluster"
  type        = string
  default     = "M10"
}
