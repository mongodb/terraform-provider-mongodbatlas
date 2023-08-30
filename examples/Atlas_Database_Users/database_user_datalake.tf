resource "mongodbatlas_database_user" "user2" {
  username           = var.user[1]
  password           = var.password[1]
  project_id         = mongodbatlas_project.project1.id
  auth_database_name = "admin"

  roles {
    role_name     = "readWrite"
    database_name = var.database_name[1]
  }
  labels {
    key   = "Name"
    value = "DB User2"
  }

  scopes {
    name = var.data_lake
    type = "DATA_LAKE"
  }
}

data "mongodbatlas_database_user" "user2" {
  username           = mongodbatlas_database_user.user2.username
  project_id         = mongodbatlas_database_user.user2.project_id
  auth_database_name = "admin"
}

output "user2" {
  value = mongodbatlas_database_user.user2.username
}
