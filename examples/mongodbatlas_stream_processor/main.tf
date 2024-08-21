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

resource "mongodbatlas_stream_processor" "stream-processor-example" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = "processorName"
  pipeline       = jsonencode({ "$source" = { "connectionName" = "sample_stream_solar" }, "$emit" = { "connectionName" : "__testLog" } }) #"[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"
  state          = "STARTED"
}

data "mongodbatlas_stream_processors" "example-stream-processors" {
  project_id    = var.project_id
  instance_name = mongodbatlas_stream_instance.example.instance_name
}

data "mongodbatlas_stream_processor" "example-stream-processor" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = mongodbatlas_stream_processor.stream-processor-example.processor_name
}

# example making use of data sources
output "stream_processors_state" {
  value = data.mongodbatlas_stream_processor.example-stream-processor.state
}

output "stream_processors_results" {
  value = data.mongodbatlas_stream_processors.example-stream-processors.results
}
