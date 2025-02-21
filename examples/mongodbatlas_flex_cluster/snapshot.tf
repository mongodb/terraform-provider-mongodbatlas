data "mongodbatlas_flex_snapshot" "snapshot" {
  project_id  = var.project_id
  name        = mongodbatlas_flex_cluster.example-cluster.name
  snapshot_id = var.snapshot_id
}

data "mongodbatlas_flex_snapshots" "snapshots" {
  project_id = var.project_id
  name       = mongodbatlas_flex_cluster.example-cluster.name
}

output "mongodbatlas_flex_snapshot" {
  value = data.mongodbatlas_flex_snapshot.snapshot.name
}

output "mongodbatlas_flex_snapshots" {
  value = [for snapshot in data.mongodbatlas_flex_snapshots.snapshots.results : snapshot.snapshot_id]
}
