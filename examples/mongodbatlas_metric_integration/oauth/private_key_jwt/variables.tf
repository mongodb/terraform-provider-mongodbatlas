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
  default     = "tf-metric-integration-oauth-private-key-jwt"
}

variable "otel_endpoint" {
  description = "OTLP-compatible endpoint URL for metric ingestion"
  type        = string
}

variable "metric_selection" {
  description = "Array of metric categories to export"
  type        = list(string)
  default     = ["ATLAS_STREAM_PROCESSING"]
}

variable "token_endpoint" {
  description = "OAuth 2.0 token endpoint URL"
  type        = string
}

variable "client_id" {
  description = "OAuth 2.0 client identifier registered with the token endpoint"
  type        = string
}

variable "token_request_params" {
  description = "Optional provider-specific parameters added to the token request"
  type        = map(string)
  default     = {}
}
