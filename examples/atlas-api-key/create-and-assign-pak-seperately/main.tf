
# Create Programmatic API Key
resource "mongodbatlas_api_key" "orgKey1" {
  description = "orgKey2"
  org_id      = var.org_id
  role_names  = ["ORG_OWNER"]
}

# Assign Newly Created Programmatic API Key to Project 
resource "mongodbatlas_project" "test2" {
  name   = "testorgapikey"
  org_id = var.org_id

  api_keys {
    api_key_id = mongodbatlas_api_key.orgKey1.api_key_id
    role_names = ["GROUP_OWNER"]
  }

  /* ensure this assignment gets cleaned up if the organization key created with api_key is deleted.  api_keys block would still need to be removed from the terraform config file */
  depends_on = [
    mongodbatlas_api_key.orgKey1
  ]

}

# Add IP Access List Entry to Programmatic API Key 
resource "mongodbatlas_access_list_api_key" "test3" {
  org_id     = var.org_id
  cidr_block = "0.0.0.0/1"
  api_key_id = mongodbatlas_api_key.orgKey1.api_key_id
}

