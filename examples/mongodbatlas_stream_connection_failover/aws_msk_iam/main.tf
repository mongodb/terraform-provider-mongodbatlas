# Create a stream workspace with a failover region and an Amazon MSK connection that authenticates with `AWS_MSK_IAM` in both the default and failover regions.
resource "mongodbatlas_stream_workspace" "example" {
  project_id     = var.project_id
  workspace_name = "workspace-with-failover-iam"

  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }

  # A failover connection can only be created for a region configured here.
  failover_regions = [
    {
      region         = "DUBLIN_IRL"
      cloud_provider = "AWS"
    }
  ]

  stream_config = {
    tier = "SP10"
  }
}

# The primary (default-region) stream connection. The IAM role must already be registered
# through Atlas Cloud Provider Access so that Atlas can assume it.
resource "mongodbatlas_stream_connection" "example" {
  project_id        = var.project_id
  workspace_name    = mongodbatlas_stream_workspace.example.workspace_name
  connection_name   = "KafkaIAMConnection"
  type              = "Kafka"
  bootstrap_servers = var.bootstrap_servers
  authentication = {
    mechanism = "AWS_MSK_IAM"
    aws = {
      role_arn = var.kafka_iam_role_arn
    }
  }
  security = {
    protocol = "SASL_SSL"
  }
}

# A failover connection shares the primary connection's name and authenticates to the MSK cluster
# in the failover region with the same IAM role.
resource "mongodbatlas_stream_connection_failover" "example" {
  project_id        = var.project_id
  workspace_name    = mongodbatlas_stream_workspace.example.workspace_name
  connection_name   = mongodbatlas_stream_connection.example.connection_name
  region            = "DUBLIN_IRL"
  type              = "Kafka"
  bootstrap_servers = var.failover_bootstrap_servers
  authentication = {
    mechanism = "AWS_MSK_IAM"
    aws = {
      role_arn = var.kafka_iam_role_arn
    }
  }
  security = {
    protocol = "SASL_SSL"
  }
}
