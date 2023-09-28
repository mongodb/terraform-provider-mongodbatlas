variable "atlas_org_id" {
  description = "Atlas organization id"
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


variable profile {
  description = "A secret with name cfn/atlas/profile/{Profile}"
  default = "default"
  type = string
}

variable atlas_project_id {
  description = "Atlas Project ID."
  type = string
}

variable database_name {
  description = "Database name for the trigger."
  type = string
}

variable collection_name {
  description = "Collection name for the trigger."
  type = string
}

variable service_id {
  description = "Service ID."
  type = string
}

variable realm_app_id {
  description = "Realm App ID."
  type = string
}

variable model_data_s3_uri {
  description = "The S3 path where the model artifacts, which result from model training, are stored. This path must point to a single gzip compressed tar archive (.tar.gz suffix)."
  type = string
}

variable model_ecr_image_uri {
  description = "AWS managed Deep Learning Container Image URI or your custom Image URI from ECR to deploy and run the model."
  type = string
}

variable pull_lambda_ecr_image_uri {
  description = "ECR image URI of the Lambda function to read MongoDB events from EventBridge."
  type = string
}

variable push_lambda_ecr_image_uri {
  description = "ECR image URI of the Lambda function to write results back to MongoDB."
  type = string
}

variable mongo_endpoint {
  description = "Your MongoDB endpoint to push results by Lambda function."
  type = string
}

variable "trigger_name" {
  description = "value of trigger name"
  type = string
  
}
