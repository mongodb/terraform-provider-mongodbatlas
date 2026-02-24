
variable "project_id" {
  description = "MongoDB Project ID"
  type        = string
}

variable "otel_endpoint" {
  description = "The Open Telemetry endpoint to connect to"
  type        = string
}
variable " otel_supplied_headers" {
  description = "The Open Telemetry supplied headers"
  type        = array
    values {
      name = "The header name"
      type = string
}