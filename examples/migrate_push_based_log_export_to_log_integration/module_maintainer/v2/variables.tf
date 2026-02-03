variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  type        = string
}

variable "bucket_name" {
  description = "The name of the S3 bucket to which Atlas will send the logs"
  type        = string
}

variable "iam_role_id" {
  description = "ID of the AWS IAM role that is used to write to the S3 bucket"
  type        = string
}

variable "prefix_path" {
  description = "S3 directory path prefix for the old push_based_log_export resource"
  type        = string
  default     = "atlas-logs"
}

variable "new_prefix_path" {
  description = "S3 directory path prefix for the new log_integration resource. Use a distinct path during migration to avoid conflicts."
  type        = string
}

variable "log_types" {
  description = "Array of log types to export. Valid values: MONGOD, MONGOS, MONGOD_AUDIT, MONGOS_AUDIT"
  type        = set(string)
  default     = ["MONGOD", "MONGOS", "MONGOD_AUDIT", "MONGOS_AUDIT"]
}

variable "skip_push_based_log_export" {
  description = "Set to true to skip creating the push_based_log_export resource. Use this during migration: first set to false (both resources active), then set to true after validating the new configuration."
  type        = bool
  default     = false
}

