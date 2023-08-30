resource "mongodbatlas_database_user" "user1" {
  username           = var.user[0]
  password           = var.password[0]
  project_id         = mongodbatlas_project.project1.id
  auth_database_name = "admin"

  roles {
    role_name     = "readWrite"
    database_name = var.database_name[0]
  }
  labels {
    key   = "Name"
    value = "DB User1"
  }

  scopes {
    name = mongodbatlas_cluster.cluster.name
    type = "CLUSTER"
  }

  scopes {
    name = mongodbatlas_cluster.cluster2.name
    type = "CLUSTER"
  }
}

data "mongodbatlas_database_user" "user1" {
  username           = mongodbatlas_database_user.user1.username
  project_id         = mongodbatlas_database_user.user1.project_id
  auth_database_name = "admin"
}

data "mongodbatlas_database_users" "allUsers" {
  project_id = mongodbatlas_database_user.user1.project_id
}

output "user1" {
  value = mongodbatlas_database_user.user1.username
}
