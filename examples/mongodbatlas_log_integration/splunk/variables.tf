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
  default     = "tf-log-integration-splunk"
}

variable "splunk_hec_token" {
  description = "Splunk HTTP Event Collector (HEC) token"
  type        = string
  sensitive   = true
}

variable "splunk_hec_url" {
  description = "Splunk HTTP Event Collector (HEC) endpoint URL including port (e.g. https://your-splunk-instance.com:8088)"
  type        = string
}
