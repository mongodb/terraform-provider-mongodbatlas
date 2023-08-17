resource "mongodbatlas_project" "atlas-project" {
  name   = "Test"
  org_id = "60ddf55c27a5a20955a707d7"
}


resource "mongodbatlas_database_user" "test" {
  username           = "test"
  project_id         = mongodbatlas_project.atlas-project.id
  auth_database_name = "admin"
  password           = "testPassword"

  roles {
    role_name     = "readWriteAnyDatabase"
    database_name = "admin"
  }

  roles {
    role_name     = "atlasAdmin"
    database_name = "admin"
  }

  labels {
    key   = "test1"
    value = "test1"
  }
  labels {
    key   = "test2"
    value = "test2"
  }
}

resource "mongodbatlas_database_user" "test2" {
  username           = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
  aws_iam_type       = "USER"
  project_id         = mongodbatlas_project.atlas-project.id
  auth_database_name = "$external"

  roles {
    role_name     = "readWriteAnyDatabase"
    database_name = "admin"
  }

  labels {
    key   = "test"
    value = "testsss"
  }
}

data "mongodbatlas_database_user" "test" {
  username           = mongodbatlas_database_user.test.username
  project_id         = mongodbatlas_database_user.test.project_id
  auth_database_name = "admin"
}

data "mongodbatlas_database_users" "test" {
  project_id = mongodbatlas_database_user.test.project_id
}