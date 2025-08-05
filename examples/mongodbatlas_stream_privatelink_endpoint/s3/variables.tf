variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  type        = string
}

variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}

variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
}

variable "region" {
  description = "AWS region where the S3 bucket is located"
  type        = string
}

variable "service_endpoint_id" {
  description = "service_endpoint_id should follow the format 'com.amazonaws.<region>.s3', for example 'com.amazonaws.us-east-1.s3'"
  type        = string
}

variable "s3_bucket_name" {
  description = "Name of the S3 bucket for stream data"
  type        = string
  default     = "mongodbatlas-stream-data"
}
