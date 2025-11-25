resource "mongodbatlas_cloud_user_org_assignment" "example" {
  org_id   = var.org_id
  username = var.user_email
  roles = {
    org_roles = ["ORG_MEMBER"]
  }
}

data "mongodbatlas_cloud_user_org_assignment" "example_username" {
  org_id   = var.org_id
  username = mongodbatlas_cloud_user_org_assignment.example.username
}

data "mongodbatlas_cloud_user_org_assignment" "example_user_id" {
  org_id  = var.org_id
  user_id = mongodbatlas_cloud_user_org_assignment.example.user_id
}
