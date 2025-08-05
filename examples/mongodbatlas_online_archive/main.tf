resource "mongodbatlas_online_archive" "users_archive" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  coll_name    = var.collection_name
  db_name      = var.database_name

  criteria {
    type              = "DATE"
    date_field        = "created"
    date_format       = "ISODATE"
    expire_after_days = 2
  }

  data_expiration_rule {
    expire_after_days = 90
  }

  partition_fields {
    field_name = "created"
    order      = 0
  }

  partition_fields {
    field_name = var.partition_field_one
    order      = 1
  }

  partition_fields {
    field_name = var.partition_field_two
    order      = 2
  }

  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}

data "mongodbatlas_online_archive" "read_archive" {
  project_id   = mongodbatlas_online_archive.users_archive.project_id
  cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
  archive_id   = mongodbatlas_online_archive.users_archive.archive_id
}

data "mongodbatlas_online_archives" "all" {
  project_id   = mongodbatlas_online_archive.users_archive.project_id
  cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
}

output "online_archive_state" {
  value = data.mongodbatlas_online_archive.read_archive.state
}

output "online_archives_results" {
  value = data.mongodbatlas_online_archives.all.results
}

