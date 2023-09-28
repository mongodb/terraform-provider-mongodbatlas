
# NOTE:
# go through the sagemaker-example/README.md file to create prerequisites and pass the inputs for the below


module "mongodb-atlas-analytics-amazon-sagemaker-integration" {
  source = "../../terraform-mongodbatlas-amazon-sagemaker-integration"

  public_key = "<public_key>"
  private_key = "<private_key>"
  atlas_org_id = "<atlas_org_id>"

  atlas_project_id = "<atlas_project_id>"
  realm_app_id = "<realm_app_id>"
  database_name = "<database_name>"
  collection_name = "<collection_name>"
  service_id = "<service_id>"

  trigger_name = "<trigger_name>"

  model_ecr_image_uri = "<model_ecr_image_uri>"
  pull_lambda_ecr_image_uri = "<pull_lambda_ecr_image_uri>"
  model_data_s3_uri = "<model_data_s3_uri>"
  push_lambda_ecr_image_uri = "<push_lambda_ecr_image_uri>"
  mongo_endpoint = "<endpoint>"
}
