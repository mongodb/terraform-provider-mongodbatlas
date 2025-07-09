# Atlas API credentials
variable "public_key" {
  type        = string
  description = "Public Programmatic API key to authenticate to Atlas"
}

variable "private_key" {
  type        = string
  description = "Private Programmatic API key to authenticate to Atlas"
}

# Organization configuration
variable "org_name" {
  type        = string
  description = "Name of the existing organization to import"
}

variable "api_access_list_required" {
  type        = bool
  description = "Flag that indicates whether to require API operations to originate from an IP Address added to the API access list"
  default     = false
}

variable "multi_factor_auth_required" {
  type        = bool
  description = "Flag that indicates whether to require users to set up Multi-Factor Authentication (MFA)"
  default     = false
}

variable "restrict_employee_access" {
  type        = bool
  description = "Flag that indicates whether to block MongoDB Support from accessing Atlas infrastructure"
  default     = false
}

variable "gen_ai_features_enabled" {
  type        = bool
  description = "Flag that indicates whether this organization has access to generative AI features"
  default     = true
}

variable "security_contact" {
  type        = string
  description = "Email address for the organization to receive security-related notifications"
  default     = ""
}

variable "skip_default_alerts_settings" {
  type        = bool
  description = "Flag that indicates whether to prevent Atlas from automatically creating organization-level alerts"
  default     = true
}

# Only needed if the organization is federated
variable "federation_settings_id" {
  type        = string
  description = "Unique identifier for the federation settings"
  default     = ""
}
