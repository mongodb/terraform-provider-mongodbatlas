
module "aws-fargate" {
  source = "../../aws-fargate"

  public_key = "ghewvngy"
  password = ["dbtestuser", "root"]
  private_key = "e0702d6b-b062-4a70-bbd0-7044c4f50f75"
  database_name = ["test1","test2"]
  atlas_org_id = "63350255419cf25e3d511c95"
  region = "US_EAST_1"

  database_user_name = "test1"
  database_password = "test2"
  org_id = "63350255419cf25e3d511c95"
  cluster_region = "US_EAST_1"
  server_service_ecr_image_uri = "816546967292.dkr.ecr.us-east-1.amazonaws.com/partner-meanstack-atlas-fargate-server:latest"
  web_access_cidr = "10.0.0.0/16"
  availability_zones = ["us-east-1a" , "us-east-1b"]
  environmentId = "dev"
  client_service_ecr_image_uri = "816546967292.dkr.ecr.us-east-1.amazonaws.com/partner-meanstack-atlas-fargate-client:latest"

}