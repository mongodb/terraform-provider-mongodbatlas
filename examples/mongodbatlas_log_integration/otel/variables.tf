
variable "project_id" {
  description = "The MongoDB Project ID"
  type        = string
}

variable "otel_endpoint" {
  description = "The Open Telemetry endpoint to connect to"
  type        = string
}
variable "otel_supplied_headers" {
  description = "The Open Telemetry supplied headers"
  type        = list(object({
    name = string
    type = string
  }))
}