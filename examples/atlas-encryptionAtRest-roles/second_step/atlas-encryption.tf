resource "mongodbatlas_encryption_at_rest" "test" {
  project_id = var.project_id

  aws_kms = {
    enabled                = true
    customer_master_key_id = var.customer_master_key
    region                 = var.atlas_region
    role_id = var.cpa_role_id
  }
}
