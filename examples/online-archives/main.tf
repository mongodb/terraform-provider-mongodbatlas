resource "mongodbatlas_online_archive" "users_archive" {
    project_id = var.project_id
    cluster_name = var.cluster_name
    coll_name = var.collection_name
    db_name = var.database_name

    criteria {
        type = "DATE"
        date_field = "date"
        date_format = "ISODATE"
        expire_after_days = 2
    }

    partition_fields {
        field_name = var.partition_field_one
        order = 0
    }

    partition_fields {
        field_name = var.partition_field_two
        order = 1
    }
}