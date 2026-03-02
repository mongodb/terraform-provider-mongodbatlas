
#MongoDB authentication variables
variable "atlas_client_id" {
  description = "The MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}

variable "atlas_client_secret" {
  description = "The MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

#Datadog variables
variable "project_id" {
  description = "The MongoDB Project ID"
  type        = string
}

variable "datadog_log_types" {
  description = "The MongoDB log type to create"
  type        = list(string)
}

variable "datadog_api_key" {
  description = "The Datadog Project API Key"
  type        = string
}
variable "datadog_region" {
  description = "The Datadog Project storage region"
  type        = string
}