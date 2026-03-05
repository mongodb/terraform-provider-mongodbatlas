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

