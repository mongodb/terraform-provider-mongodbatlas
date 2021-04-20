resource "mongodbatlas_project" "test" {
  name   = var.project_name
  org_id = var.org_id
}

resource "mongodbatlas_database_user" "test" {
  username           = var.username
  ldap_auth_type     = "USER"
  project_id         = mongodbatlas_project.test.id
  auth_database_name = "$external"

  roles {
    role_name     = var.role_name
    database_name = "admin"
  }

  labels {
    key   = "First Key"
    value = "First Value"
  }
}
