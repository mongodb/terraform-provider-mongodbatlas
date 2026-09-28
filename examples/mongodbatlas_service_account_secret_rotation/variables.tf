variable "org_id" {
  description = "Atlas Organization ID where this configuration creates the Service Account."
  type        = string
}

variable "service_account_name" {
  description = "Name of the Service Account to create and rotate."
  type        = string
  default     = "example-rotation-service-account"
}

variable "service_account_roles" {
  description = "Roles for the Service Account. It needs a role that can manage its own secrets to rotate while authenticated as itself."
  type        = list(string)
  default     = ["ORG_OWNER"]
}

variable "secret_expires_after_hours" {
  description = "Requested lifetime of each secret in hours. The overlap between the two slots is at most 7 days regardless of this value."
  type        = number
  default     = 2160 # 90 days
}
