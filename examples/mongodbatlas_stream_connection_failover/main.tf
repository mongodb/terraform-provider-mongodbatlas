resource "mongodbatlas_stream_workspace" "example" {
  project_id     = var.project_id
  workspace_name = "workspace-with-failover"

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

# The primary (default-region) stream connection.
resource "mongodbatlas_stream_connection" "example" {
  project_id        = var.project_id
  workspace_name    = mongodbatlas_stream_workspace.example.workspace_name
  connection_name   = "KafkaConnection"
  type              = "Kafka"
  bootstrap_servers = var.bootstrap_servers
  authentication = {
    mechanism = "PLAIN"
    username  = var.kafka_username
    password  = var.kafka_password
  }
  security = {
    protocol = "SASL_SSL"
  }
}

# A failover connection shares the primary connection's name and is created for one of the
# workspace's failover regions, with its own regional configuration.
resource "mongodbatlas_stream_connection_failover" "example" {
  project_id        = var.project_id
  workspace_name    = mongodbatlas_stream_workspace.example.workspace_name
  connection_name   = mongodbatlas_stream_connection.example.connection_name
  region            = "DUBLIN_IRL"
  type              = "Kafka"
  bootstrap_servers = var.failover_bootstrap_servers
  authentication = {
    mechanism = "PLAIN"
    username  = var.failover_kafka_username
    password  = var.failover_kafka_password
  }
  security = {
    protocol = "SASL_SSL"
  }
}
