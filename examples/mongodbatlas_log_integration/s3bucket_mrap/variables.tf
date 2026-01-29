variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "name_prefix" {
  description = "Prefix for naming AWS resources (must be globally unique for S3)"
  type        = string
  default     = "atlas-logs"
}

variable "prefix_path" {
  description = "S3 directory path prefix for log files"
  type        = string
  default     = "atlas-logs/"
}

variable "log_types" {
  description = "Array of log types to export"
  type        = list(string)
  default     = ["MONGOD", "MONGOS"]
}
