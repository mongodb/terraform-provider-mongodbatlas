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
  default     = "tf-log-integration-otel"
}

variable "otel_endpoint" {
  description = "OpenTelemetry collector endpoint URL for log ingestion (e.g. https://your-otel-collector.com:4318/v1/logs)"
  type        = string
}

variable "otel_supplied_headers" {
  description = "Custom headers to include in OTel log export requests (e.g. authentication tokens)"
  type = list(object({
    name  = string
    value = string
  }))
  sensitive = true
}
