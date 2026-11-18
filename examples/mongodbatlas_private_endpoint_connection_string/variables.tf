variable "project_id" {
  description = "Unique identifier of the MongoDB Atlas project that contains the cluster."
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label of the existing dedicated AWS cluster that the connection string points to. The cluster must run in the same AWS region as the private endpoint."
  type        = string
}

variable "connection_string_name" {
  description = "Human-readable label of at most 64 characters that identifies the private endpoint connection string."
  type        = string
  default     = "my-connection-string"
}

variable "optimized_mode_enabled" {
  description = "Flag that enables the optimized load-balanced mode. Requires an AWS sharded cluster and private endpoints that meet the optimized connection string requirements. This value cannot be changed after creation."
  type        = bool
  default     = false
}

variable "atlas_region_name" {
  description = "Atlas region name for the private endpoint, for example `US_EAST_1`. It must match the region of the cluster and resolve to the same AWS region as `aws_region`."
  type        = string
  default     = "US_EAST_1"
}

variable "aws_region" {
  description = "AWS region for the VPC endpoint, for example `us-east-1`. It must match the AWS region of the Atlas cluster."
  type        = string
  default     = "us-east-1"
}

variable "vpc_id" {
  description = "Identifier of the VPC that hosts the interface endpoint."
  type        = string
}

variable "subnet_ids" {
  description = "Identifiers of the subnets that host the interface endpoint's elastic network interfaces."
  type        = list(string)
}

variable "security_group_ids" {
  description = "Identifiers of the security groups that the interface endpoint's elastic network interfaces use."
  type        = list(string)
}
