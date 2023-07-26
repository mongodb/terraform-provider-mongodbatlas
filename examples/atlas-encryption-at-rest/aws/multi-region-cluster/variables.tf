variable "atlas_org_id" {
  description = "Atlas Organization ID (e.g. 5beae24579358e0ae95492af)"
  type        = string
}

variable "atlas_project_name" {
  description = "Name of the project to create"
  type        = string
  default     = "TestProject"
}

variable "atlas_region" {
  description = "The AWS region in which the AWS customer master key exists: CA_CENTRAL_1, US_EAST_1, US_EAST_2, US_WEST_1, US_WEST_2, SA_EAST_1"
  type        = string
  default     = "US_EAST_1"
}

variable "aws_region" {
  description = "AWS region (e.g. us-east-1) where "
  type        = string
  default     = "us-east-1"
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

variable "aws_region_shard_3" {
  description = "Region of the third shard"
  type        = string
  default     = "US_WEST_2"
}
