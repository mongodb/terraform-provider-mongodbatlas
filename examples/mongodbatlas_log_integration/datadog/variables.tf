
variable "project_id" {
  description = "MongoDB Project ID"
  type        = string
}

variable "log_types" {
  description = "The MongoDB log type to create"
  type        = string array
}

variable "api_key" {
  description = "The Datadog Project API Key"
  type        = string
}
variable "region" {
  description = "The Datadog Project storage region"
  type        = string
}