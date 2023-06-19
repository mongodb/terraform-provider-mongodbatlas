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

variable "collection_1" {
  type        = string
  description = "Human-readable label that identifies the collection in the database in first cluster"
}

variable "database_1" {
  type        = string
  description = "Human-readable label that identifies the database, which contains the collection in the first cluster"
}

variable "collection_2" {
  type        = string
  description = "Human-readable label that identifies the collection in the database in second cluster"
}

variable "database_2" {
  type        = string
  description = "Human-readable label that identifies the database, which contains the collection in the second cluster"
}

variable "federated_instance_name" {
  type        = string
  description = "MongoDB Federated Database Instance Name."
  default     = "FederatedDatabaseInstance0"
}

variable "atlas_cluster_name_1" {
  type        = string
  description = "First Atlas Cluster Name."
  default     = "ClusterFederatedTest1"
}

variable "atlas_cluster_name_2" {
  type        = string
  description = "Second Atlas Cluster Name."
  default     = "ClusterFederatedTest2"
}

variable "provider_region_name" {
  type        = string
  description = "Physical location where MongoDB Cloud deploys your AWS-hosted MongoDB cluster nodes."
  default     = "US_EAST_1"
}

variable "provider_instance_size_name" {
  type        = string
  description = "Cluster tier. Default is M10"
  default     = "M10"
}

variable "backing_provider_name" {
  type    = string
  default = "AWS"
}

variable "provider_name" {
  type    = string
  default = "AWS"
}

variable "federated_query_limit" {
  type        = string
  description = "Human-readable label that identifies the user-managed limit to modify."
  default     = "bytesProcessed.monthly"
}

variable "overrun_policy" {
  type        = string
  description = "Action to take when the usage limit is exceeded."
  default     = "BLOCK"
}

variable "limit_value" {
  type        = number
  description = "Amount to set the federated query limit to."
  default     = 5147483648
}
