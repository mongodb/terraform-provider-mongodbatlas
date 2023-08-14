resource "mongodbatlas_project" "atlas-project" {
  name   = "Test API Keys"
  org_id = var.org_id
}

resource "mongodbatlas_project_api_key" "api_1" {
  description = "test api_key multi"
  project_id  = mongodbatlas_project.atlas-project.id

  project_assignment {
    project_id = mongodbatlas_project.atlas-project.id
    role_names = ["ORG_OWNER", "GROUP_OWNER"]
  }
}

# Add IP Access List Entry to Programmatic API Key 
resource "mongodbatlas_access_list_api_key" "test3" {
  org_id     = var.org_id
  cidr_block = "0.0.0.0/1"
  api_key_id = mongodbatlas_project_api_key.api_1.api_key_id
}

