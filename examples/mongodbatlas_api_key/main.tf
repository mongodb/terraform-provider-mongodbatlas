resource "mongodbatlas_api_key" "this" {
  org_id      = var.org_id
  description = "Test API Key"
  role_names  = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.org_id
}

resource "mongodbatlas_api_key_project_assignment" "this" {
  project_id = mongodbatlas_project.this.id
  api_key_id = mongodbatlas_api_key.this.api_key_id
  roles      = ["GROUP_OWNER"]
}

# Add IP Access List Entry to Programmatic API Key 
resource "mongodbatlas_access_list_api_key" "this" {
  org_id     = var.org_id
  cidr_block = "0.0.0.0/1"
  api_key_id = mongodbatlas_api_key.this.api_key_id
}
