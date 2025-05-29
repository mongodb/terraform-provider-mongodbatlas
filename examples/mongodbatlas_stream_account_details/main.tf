data "mongodbatlas_stream_account_details" "account_details" {
  project_id = var.project_id
  cloud_provider= "aws"
  region_name = "US_EAST_1"
}

output "aws_account_id" {
  value = data.mongodbatlas_stream_account_details.account_details.aws_account_id
}

output "cidr_block" {
  value = data.mongodbatlas_stream_account_details.account_details.cidr_block
}

output "cloud_provider" {
  value = data.mongodbatlas_stream_account_details.account_details.cloud_provider
}

output "links" {
  value = data.mongodbatlas_stream_account_details.account_details.links
}

output "vpc_id" {
  value = data.mongodbatlas_stream_account_details.account_details.vpc_id
}