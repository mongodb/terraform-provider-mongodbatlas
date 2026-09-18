variable "atlas_org_id" {
  description = "Atlas Organization ID"
  type        = string
}

variable "atlas_project_name" {
  description = "Name of the Atlas project"
  type        = string
  default     = "tf-metric-integration-oauth-client-secret"
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

variable "client_secret" {
  description = "OAuth 2.0 client secret"
  type        = string
  sensitive   = true
}

variable "oauth_scopes" {
  description = "OAuth 2.0 scopes requested on the token"
  type        = list(string)
  default     = []
}
