resource "mongodbatlas_project" "example" {
  name   = "project-name"
  org_id = var.org_id
}

resource "mongodbatlas_stream_workspace" "example" {
  project_id     = mongodbatlas_project.example.id
  workspace_name = "WorkspaceName"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
  stream_config = {
    tier = "SP30"
  }
}
