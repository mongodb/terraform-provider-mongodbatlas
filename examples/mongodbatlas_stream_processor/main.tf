resource "mongodbatlas_stream_instance" "example" {
  project_id    = var.project_id
  instance_name = "InstanceName"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

resource "mongodbatlas_stream_connection" "example-sample" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "sample_stream_solar"
  type            = "Sample"
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

resource "mongodbatlas_stream_connection" "example-kafka" {
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
}

resource "mongodbatlas_stream_processor" "stream-processor-sample-example" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = "sampleProcessorName"
  pipeline       = jsonencode([{ "$source" = { "connectionName" = resource.mongodbatlas_stream_connection.example-sample.connection_name } }, { "$emit" = { "connectionName" : "__testLog" } }])
  state          = "CREATED"
}

resource "mongodbatlas_stream_processor" "stream-processor-cluster-example" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = "clusterProcessorName"
  pipeline       = jsonencode([{ "$source" = { "connectionName" = resource.mongodbatlas_stream_connection.example-cluster.connection_name } }, { "$emit" = { "connectionName" : "__testLog" } }])
  state          = "STARTED"
}

resource "mongodbatlas_stream_processor" "stream-processor-kafka-example" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = "kafkaProcessorName"
  pipeline       = jsonencode([{ "$source" = { "connectionName" = resource.mongodbatlas_stream_connection.example-cluster.connection_name } }, { "$emit" = { "connectionName" : resource.mongodbatlas_stream_connection.example-kafka.connection_name, "topic" : "example_topic" } }])
  state          = "CREATED"
  options = {
    dlq = {
      coll            = "exampleColumn"
      connection_name = resource.mongodbatlas_stream_connection.example-cluster.connection_name
      db              = "exampleDb"
    }
  }
}

data "mongodbatlas_stream_processors" "example-stream-processors" {
  project_id    = var.project_id
  instance_name = mongodbatlas_stream_instance.example.instance_name
}

data "mongodbatlas_stream_processor" "example-stream-processor" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = mongodbatlas_stream_processor.stream-processor-sample-example.processor_name
}

# example making use of data sources
output "stream_processors_state" {
  value = data.mongodbatlas_stream_processor.example-stream-processor.state
}

output "stream_processors_results" {
  value = data.mongodbatlas_stream_processors.example-stream-processors.results
}
