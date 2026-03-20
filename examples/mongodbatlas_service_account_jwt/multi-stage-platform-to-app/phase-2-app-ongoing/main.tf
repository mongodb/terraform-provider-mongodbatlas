# Phase 2 (App Ongoing): Read SA credentials from Secrets Manager and provision
# infrastructure. This phase runs independently of the bootstrap JWT.

data "aws_secretsmanager_secret_version" "sa_creds" {
  secret_id = var.aws_secret_id
}

resource "mongodbatlas_advanced_cluster" "this" {
  project_id   = var.project_id
  name         = var.cluster_name
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M10"
            node_count    = 3
          }
          provider_name = "AWS"
          region_name   = "US_EAST_1"
          priority      = 7
        }
      ]
    }
  ]
}
