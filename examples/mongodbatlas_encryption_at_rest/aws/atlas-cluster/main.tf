resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.atlas_project_id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.atlas_project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  aws {
    iam_assumed_role_arn = aws_iam_role.test_role.arn
  }
}

resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.atlas_project_id

  aws_kms_config {
    enabled                = true
    customer_master_key_id = aws_kms_key.kms_key.id
    region                 = var.atlas_region
    role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  }
  
  enabled_for_search_nodes = true 
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id                  = mongodbatlas_encryption_at_rest.test.project_id
  name                        = "MyCluster"
  cluster_type                = "REPLICASET"
  backup_enabled              = true
  encryption_at_rest_provider = "AWS"

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}

data "mongodbatlas_encryption_at_rest" "test" {
  project_id = mongodbatlas_encryption_at_rest.test.project_id
}

output "is_aws_kms_encryption_at_rest_valid" {
  value = data.mongodbatlas_encryption_at_rest.test.aws_kms_config.valid
}
