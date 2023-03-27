variable "public_key" {
  description = "MongoDB Atlas authentication Public API key"
}
variable "private_key" {
  description = "MongoDB Atlas authentication Private API key"
}
variable "project_id" {
  description = "MongoDB Atlas project id"
}
variable "datadog_api_key" {
  description = "Datadog api key"
}
variable "datadog_region" {
  description = "Datadog region"
  default     = "US5"
}
variable "cluster_name" {
  description = "Cluster to test regional mode"
  default     = "datadog-test-cluster"
}