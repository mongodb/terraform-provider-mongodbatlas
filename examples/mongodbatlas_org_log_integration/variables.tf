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

variable "otel_endpoint" {
  description = "OpenTelemetry collector endpoint URL for event export (e.g. https://your-otel-collector.com:4318/v1/logs)"
  type        = string
}

variable "otel_supplied_headers" {
  description = "Custom headers to include in OTel export requests (e.g. authentication tokens)"
  type = list(object({
    name  = string
    value = string
  }))
  sensitive = true
}
