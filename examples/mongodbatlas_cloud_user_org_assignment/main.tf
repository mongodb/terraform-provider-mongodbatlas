resource "mongodbatlas_cloud_user_org_assignment" "example" {
  org_id   = var.org_id
  username = var.user_email
  roles = {
    org_roles = ["ORG_MEMBER"]
  }
} 