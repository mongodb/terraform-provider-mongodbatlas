variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
  default     = ""
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
  default     = ""
}

variable "project_id" {
  type        = string
  description = "Atlas Project ID"
}

variable "audit_filter_json" {
  type        = string
  description = "Path to the JSON file containing the audit filter configuration. Will use audit_filter.json as the default value."
  default     = ""
}
