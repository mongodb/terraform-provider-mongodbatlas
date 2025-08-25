provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = var.project_id
  name         = "cluster"
  cluster_type = "REPLICASET"

  mongo_db_major_version = "7.0"

  pinned_fcv = {
    expiration_date = var.fcv_expiration_date # e.g. format: "2024-11-22T10:50:00Z". Hashicorp time provider https://registry.terraform.io/providers/hashicorp/time/latest/docs/resources/offset can be used to compute this string value.
  }

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10"
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    }
  ]
}

output "feature_compatibility_version" {
  value = mongodbatlas_advanced_cluster.cluster.pinned_fcv.version
}

