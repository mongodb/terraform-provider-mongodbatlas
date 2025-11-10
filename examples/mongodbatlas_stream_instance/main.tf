# DEPRECATED: For new deployments, use mongodbatlas_stream_workspace instead.
# See ../mongodbatlas_stream_workspace/ for the updated example.

resource "mongodbatlas_project" "example" {
  name   = "project-name"
  org_id = var.org_id
}

# Add this moved block to migrate from stream_instance to stream_workspace:
# moved {
#   from = mongodbatlas_stream_instance.example
#   to   = mongodbatlas_stream_workspace.example
# }
