resource "mongodbatlas_project" "example" {
  name   = "project-name"
  org_id = var.org_id
}

resource "mongodbatlas_stream_instance" "example" {
  project_id = mongodbatlas_project.example
	instance_name = "InstanceName"
	data_process_region = {
		region = "VIRGINIA_USA"
		cloud_provider = "AWS"
  }
}