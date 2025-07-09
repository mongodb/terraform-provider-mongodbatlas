resource "mongodbatlas_api_key" "this" {
  org_id      = var.org_id
  description = "Test API Key"
  role_names  = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_project" "first_project" {
  name   = "First Project"
  org_id = var.org_id
}

resource "mongodbatlas_project" "second_project" {
  name   = "Second Project"
  org_id = var.org_id
}

resource "mongodbatlas_api_key_project_assignment" "first_assignment" {
  project_id = mongodbatlas_project.first_project.id
  api_key_id = mongodbatlas_api_key.this.api_key_id
  roles      = ["GROUP_OWNER"]
}

resource "mongodbatlas_api_key_project_assignment" "second_assignment" {
  project_id = mongodbatlas_project.second_project.id
  api_key_id = mongodbatlas_api_key.this.api_key_id
  roles      = ["GROUP_OWNER"]
}

# Add IP Access List Entry to Programmatic API Key 
resource "mongodbatlas_access_list_api_key" "this" {
  org_id     = var.org_id
  cidr_block = "0.0.0.0/1"
  api_key_id = mongodbatlas_api_key.this.api_key_id
}
