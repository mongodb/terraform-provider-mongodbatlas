variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
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
  default     = "mongodb-federation-database-instance-test"
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
  type        = string
}
variable "secret_key" {
  description = "The secret key for an AWS Account"
  type        = string
}

variable "aws_region" {
  default     = "ap-southeast-1"
  description = "AWS Region"
  type        = string
}

variable "mongodb_aws_region" {
  description = "AWS Region used for the stores in mongodbatlas_federated_database_instance (e.g. EU_WEST_1)"
  type        = string
}
