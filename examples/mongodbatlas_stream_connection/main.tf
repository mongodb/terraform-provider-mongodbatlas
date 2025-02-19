resource "mongodbatlas_stream_instance" "example" {
  project_id    = var.project_id
  instance_name = "InstanceName"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

resource "mongodbatlas_stream_connection" "example-cluster" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "ClusterConnection"
  type            = "Cluster"
  cluster_name    = var.cluster_name
  db_role_to_execute = {
    role = "atlasAdmin"
    type = "BUILT_IN"
  }
}

resource "mongodbatlas_stream_connection" "example-kafka-plaintext" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "KafkaPlaintextConnection"
  type            = "Kafka"
  authentication = {
    mechanism = "PLAIN"
    username  = var.kafka_username
    password  = var.kafka_password
  }
  bootstrap_servers = "localhost:9092,localhost:9092"
  config = {
    "auto.offset.reset" : "earliest"
  }
  security = {
    protocol = "PLAINTEXT"
  }
  networking = {
    access = {
      type = "PUBLIC"
    }
  }
}

resource "mongodbatlas_stream_connection" "example-kafka-ssl" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "KafkaSSLConnection"
  type            = "Kafka"
  authentication = {
    mechanism = "PLAIN"
    username  = var.kafka_username
    password  = var.kafka_password
  }
  bootstrap_servers = "localhost:9092,localhost:9092"
  config = {
    "auto.offset.reset" : "earliest"
  }
  security = {
    broker_public_certificate = var.kafka_ssl_cert
    protocol                  = "SSL"
  }
}

resource "mongodbatlas_stream_connection" "example-sample" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "sample_stream_solar"
  type            = "Sample"
}

resource "mongodbatlas_stream_connection" "example-aws-lambda" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "AWSLambdaConnection"
  type            = "AWSLambda"
  aws             = {
    role_arn = "arn:aws:iam::123456789123:role/lambdaRole"
  }
}

data "mongodbatlas_stream_connection" "example-kafka-ssl" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = mongodbatlas_stream_connection.example-kafka-ssl.connection_name
}

data "mongodbatlas_stream_connections" "example" {
  project_id    = var.project_id
  instance_name = mongodbatlas_stream_instance.example.instance_name
}

# example making use of data sources
output "stream_connection_bootstrap_servers" {
  value = data.mongodbatlas_stream_connection.example-kafka-ssl.bootstrap_servers
}

output "stream_connection_total_count" {
  value = data.mongodbatlas_stream_connections.example.total_count
}