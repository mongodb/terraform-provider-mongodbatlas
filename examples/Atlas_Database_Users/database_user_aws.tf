resource "mongodbatlas_database_user" "user3" {
  username           = var.aws_iam_role_arn
  project_id         = mongodbatlas_project.project1.id
  auth_database_name = "$external"
  aws_iam_type       = "ROLE"

  roles {
    role_name     = "readAnyDatabase"
    database_name = "admin"
  }

  labels {
    key   = "Env"
    value = "Test"
  }

  scopes {
    name   = "My cluster name"
    type = "CLUSTER"
  }
}

data "mongodbatlas_database_user" "user3" {
  username           = mongodbatlas_database_user.user3.username
  project_id         = mongodbatlas_database_user.user3.project_id
  auth_database_name = "admin"
}

output "user3" {
  value = data.mongodbatlas_database_user.user3.username
}
