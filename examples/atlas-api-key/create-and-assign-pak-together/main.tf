
# Create Programmatic API Key + Assign Programmatic API Key to Project
resource "mongodbatlas_project_api_key" "test" {
  description = "test create and assign"
  project_id  = var.project_id
  role_names  = ["GROUP_OWNER"]
}

# Add IP Access List Entry to Programmatic API Key 
resource "mongodbatlas_access_list_api_key" "test3" {
  org_id     = var.org_id
  cidr_block = "0.0.0.0/1"
  api_key_id = mongodbatlas_project_api_key.test.api_key_id
}
