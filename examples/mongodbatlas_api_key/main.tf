resource "mongodbatlas_api_key" "test" {
  org_id      = var.org_id
  description = "Test API Key"

  role_names = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_project" "test1" {
  name   = var.project_name
  org_id = var.org_id
}

resource "mongodbatlas_api_key_project_assignment" "test1" {
  project_id = mongodbatlas_project.test1.id
  api_key_id = mongodbatlas_api_key.test.api_key_id
  roles      = ["GROUP_OWNER"]
}
