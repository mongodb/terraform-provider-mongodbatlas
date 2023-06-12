
# Create Programmatic API Key + Assign Programmatic API Key to Project
resource "mongodbatlas_project_api_key" "test" {
  description = "test create and assign"
  project_id  = var.project_id

  project_assignment {
    project_id = var.project_id
    role_names = ["GROUP_READ_ONLY", "GROUP_OWNER"]
  }
  project_assignment {
    project_id = var.additional_project_id
    role_names = ["GROUP_READ_ONLY"]
  }
}


