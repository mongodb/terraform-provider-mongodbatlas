resource "mongodbatlas_data_lake" "test" {
  project_id         = var.project_id
  name = var.data_lake_name
  aws_role_id = var.cpa_role_id
  aws_test_s3_bucket = var.test_s3_bucket
  data_process_region = {
    cloud_provider = "AWS"
    region = var.data_lake_region
  }
}
