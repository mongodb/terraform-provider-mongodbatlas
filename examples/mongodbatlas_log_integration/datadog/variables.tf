
variable "project_id" {
  description = "The MongoDB Project ID"
  type        = string
}

variable "datadog_log_types" {
  description = "The MongoDB log type to create"
  type        = string array
}

variable "datadog_api_key" {
  description = "The Datadog Project API Key"
  type        = string
}
variable "datadog_region" {
  description = "The Datadog Project storage region"
  type        = string
}