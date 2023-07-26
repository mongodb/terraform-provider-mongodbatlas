variable "atlas_project_id" {
  description = "Atlas Project ID"
  type        = string
}

variable "cluster_name" {
  description = "Atlas cluster name"
  type        = string
  default     = "MultiCloudCluster"
}

variable "aws_region_shard_1" {
  description = "Region of the first shard"
  type        = string
  default     = "US_EAST_1"
}

variable "aws_region_shard_2" {
  description = "Region of the second shard"
  type        = string
  default     = "US_EAST_2"
}

variable "provider_name" {
  description = "Name of the provider to use for the cluster"
  type        = string
  default     = "AWS"
}

variable "instance_size" {
  description = "Instance Size of the cluster"
  type        = string
  default     = "M10"
}