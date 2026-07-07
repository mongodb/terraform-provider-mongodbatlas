resource "mongodbatlas_stream_connection" "example" {
  project_id        = var.project_id
  workspace_name    = var.workspace_name
  connection_name   = "KafkaConnection"
  type              = "Kafka"
  bootstrap_servers = "localhost:9092,localhost:9092"
  authentication = {
    mechanism = "PLAIN"
    username  = var.kafka_username
    password  = var.kafka_password
  }
  security = {
    protocol = "SASL_SSL"
  }
  config = {
    "auto.offset.reset" = "earliest"
  }
}

# A failover connection shares the primary connection's name and is created for one of the
# workspace's failover regions. It carries its own regional connection configuration.
resource "mongodbatlas_stream_connection_failover" "example" {
  project_id        = var.project_id
  workspace_name    = var.workspace_name
  connection_name   = mongodbatlas_stream_connection.example.connection_name
  region            = var.failover_region
  type              = "Kafka"
  bootstrap_servers = var.failover_bootstrap_servers
  authentication = {
    mechanism = "PLAIN"
    username  = var.kafka_username
    password  = var.kafka_password
  }
  security = {
    protocol = "SASL_SSL"
  }
  config = {
    "auto.offset.reset" = "earliest"
  }
}
