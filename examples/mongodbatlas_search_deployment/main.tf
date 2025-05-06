resource "mongodbatlas_project" "example" {
  name   = "project-name"
  org_id = var.org_id
}

resource "mongodbatlas_advanced_cluster" "example" {
  project_id   = mongodbatlas_project.example.id
  name         = "ClusterExample"
  cluster_type = "REPLICASET"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
  }
}

resource "mongodbatlas_search_deployment" "example" {
  project_id   = mongodbatlas_project.example.id
  cluster_name = mongodbatlas_advanced_cluster.example.name
  specs = [
    {
      instance_size = "S20_HIGHCPU_NVME"
      node_count    = 2
    }
  ]
}

data "mongodbatlas_search_deployment" "example" {
  project_id   = mongodbatlas_search_deployment.example.project_id
  cluster_name = mongodbatlas_search_deployment.example.cluster_name
}

output "mongodbatlas_search_deployment_id" {
  value = data.mongodbatlas_search_deployment.example.id
}

output "mongodbatlas_search_deployment_encryption_at_rest_provider" {
  value = data.mongodbatlas_search_deployment.example.encryption_at_rest_provider
}
