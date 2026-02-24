
variable "project_id" {
  description = "The MongoDB Project ID"
  type        = string
}

variable "log_types" {
  description = "The MongoDB log type to create"
  type        = string array
}

variable "hec_token" {
  description = "The Splunk HTTP Event Collector (HEC) token"
  type        = string
}
variable "hec_url" {
  description = "The Splunk HTTP Event Collector (HEC) URL"
  type        = string
}