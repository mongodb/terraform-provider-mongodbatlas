variable "public_key" {
  description = "MongoDB Atlas authentication Public API key"
  type        = string
}
variable "private_key" {
  description = "MongoDB Atlas authentication Private API key"
  type        = string
}
variable "project_id" {
  description = "MongoDB Atlas project id"
  type        = string
}
variable "datadog_api_key" {
  description = "Datadog api key"
  type        = string
}
variable "datadog_region" {
  description = "Datadog region"
  default     = "US5"
  type        = string
}
variable "cluster_name" {
  description = "Cluster to test regional mode"
  default     = "datadog-test-cluster"
  type        = string
}

variable "send_collection_latency_metrics" {
  description = "Send collection latency metrics (only for Datadog integrations)"
  default     = false
  type        = bool
}

variable "send_database_metrics" {
  description = "Send database metrics (only for Datadog integrations)"
  default     = false
  type        = bool
}

variable "send_user_provided_resource_tags" {
  description = "Send user provided resource tags (only for Datadog integrations)"
  default     = false
  type        = bool
}