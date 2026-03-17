variable "atlas_client_id" {
  description = "Atlas Service Account Client ID"
  type        = string
  default     = ""
}

variable "atlas_client_secret" {
  description = "Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "atlas_org_id" {
  description = "Atlas Organization ID"
  type        = string
}

variable "atlas_project_name" {
  description = "Name of the Atlas project"
  type        = string
  default     = "tf-log-integration-datadog"
}

variable "datadog_api_key" {
  description = "Datadog API key for authentication"
  type        = string
  sensitive   = true
}

variable "datadog_region" {
  description = "Datadog region"
  type        = string
  default     = "US1"
}
