variable "public_key" {
  type        = string
  description = "Public Programmatic API key to authenticate to Atlas"
}
variable "private_key" {
  type        = string
  description = "Private Programmatic API key to authenticate to Atlas"
}
variable "project_id" {
  type        = string
  description = "MongoDB Project ID"
}

variable "test_s3_bucket" {
  type        = string
  description = "Name of the S3 data bucket that the provided role ID is authorized to access"
}

variable "collection" {
  type        = string
  description = "Human-readable label that identifies the collection in the database"
}

variable "database" {
  type        = string
  description = "Human-readable label that identifies the database, which contains the collection in the cluster"
}

variable "path" {
  type        = string
  description = "File path that controls how MongoDB Cloud searches for and parses files in the storeName before mapping them to a collection"
}

variable "prefix" {
  type        = string
  description = "Prefix that controls how MongoDB Cloud searches for and parses files in the storeName before mapping them to a collection"
}

variable "name" {
  type        = string
  description = "MongoDB Federated Database Instance Name"
  default     = "mongodb_federation_database_instance_test"
}

variable "policy_name" {
  type        = string
  description = "AWS Policy Name"
  default     = "mongodb_federation_database_instance_policy"
}

variable "role_name" {
  type        = string
  description = "AWS Role Name"
  default     = "mongodb_federation_database_instance_role"
}

variable "atlas_cluster_name" {
  type        = string
  description = "Atlas Cluster Name"
  default     = "ClusterFederatedTest"
}

variable "access_key" {
  description = "The access key for an AWS Account"
}
variable "secret_key" {
  description = "The secret key for an AWS Account"
}

variable "aws_region" {
  default     = "ap-southeast-1"
  description = "AWS Region"
}
