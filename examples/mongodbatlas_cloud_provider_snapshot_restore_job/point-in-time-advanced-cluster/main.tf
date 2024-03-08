# This will Create a Project,  Cluster, cloud backup snapshot and restore job

resource "mongodbatlas_project" "project_test" {
  name   = var.project_name
  org_id = var.org_id
}

resource "mongodbatlas_advanced_cluster" "advanced_cluster_test" {
  project_id   = mongodbatlas_project.project_test.id
  name         = var.cluster_name
  cluster_type = "REPLICASET"

  replication_specs {
    num_shards = 1

    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }

      provider_name = "AWS"
      region_name   = "US_EAST_1"
      priority      = 7
    }
  }

  backup_enabled         = true                       # enable cloud backup snapshots
  pit_enabled            = true                       # Flag that indicates whether the cluster uses continuous cloud backups
  retain_backups_enabled = var.retain_backups_enabled # keep the backup snapshots once the cluster is deleted
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_advanced_cluster.advanced_cluster_test.project_id
  cluster_name      = mongodbatlas_advanced_cluster.advanced_cluster_test.name
  description       = "My description"
  retention_in_days = "1"
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  count        = (var.point_in_time_utc_seconds == 0 ? 0 : 1)
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id # for provider versions > 1.15.1 the snapshot_id field can be omitted when delivery_type_config is "point_in_time"

  delivery_type_config {
    point_in_time             = true
    target_cluster_name       = mongodbatlas_advanced_cluster.advanced_cluster_test.name
    target_project_id         = mongodbatlas_advanced_cluster.advanced_cluster_test.project_id
    point_in_time_utc_seconds = var.point_in_time_utc_seconds
  }
}
