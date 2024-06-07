locals {
  test_user_username = "user_${random_password.username.result}"
  test_user_password = random_password.password.result
}

resource "random_password" "password" {
  length  = 12
  special = false
}

resource "random_password" "username" {
  length  = 12
  special = false
}

resource "mongodbatlas_database_user" "test_user" {
  username           = local.test_user_username
  password           = local.test_user_password
  project_id         = mongodbatlas_project.this.id
  auth_database_name = "admin"

  roles {
    role_name     = "atlasAdmin"
    database_name = "admin"
  }
}
