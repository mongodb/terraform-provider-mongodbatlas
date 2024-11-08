resource "mongodbatlas_flex_cluster" "example-cluster" {
  project_id = var.project_id
  name       = var.cluster_name
  provider_settings = {
    backing_provider_name = "AWS"
    region_name           = "US_EAST_1"
  }
  termination_protection_enabled = true
}

data "mongodbatlas_flex_cluster" "example-cluster" {
  project_id = var.project_id
  name       = mongodbatlas_flex_cluster.example-cluster.name
}

data "mongodbatlas_flex_clusters" "example-clusters" {
  project_id = var.project_id
}

output "mongodbatlas_flex_cluster" {
  value = data.mongodbatlas_flex_cluster.example-cluster.name
}

output "mongodbatlas_flex_clusters_names" {
  value = [for cluster in data.mongodbatlas_flex_clusters.example-clusters.results : cluster.name]
}
