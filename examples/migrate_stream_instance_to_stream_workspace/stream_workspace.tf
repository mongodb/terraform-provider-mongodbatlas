# AFTER: New stream workspace resource
resource "mongodbatlas_stream_workspace" "example" {
  project_id     = var.project_id
  workspace_name = var.workspace_name # Note: instance_name -> workspace_name
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
  stream_config = {
    tier = "SP30"
  }
}

# Moved block to migrate from stream_instance to stream_workspace
moved {
  from = mongodbatlas_stream_instance.example
  to   = mongodbatlas_stream_workspace.example
}
