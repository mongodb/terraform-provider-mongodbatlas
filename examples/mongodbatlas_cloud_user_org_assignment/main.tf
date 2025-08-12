resource "mongodbatlas_cloud_user_org_assignment" "example" {
  org_id   = var.org_id
  username = var.user_email
  roles = {
    org_roles = ["ORG_MEMBER"]
  }
}

data "mongodbatlas_cloud_user_org_assignment" "example_username" {
  org_id     = var.org_id
  username   = var.user_email
  depends_on = [mongodbatlas_cloud_user_org_assignment.example]
}

data "mongodbatlas_cloud_user_org_assignment" "example_user_id" {
  org_id     = var.org_id
  user_id    = var.user_id
  depends_on = [mongodbatlas_cloud_user_org_assignment.example]
}
