resource "mongodbatlas_project" "atlas-project" {
  org_id = var.atlas_org_id
  name   = var.atlas_project_name
}

resource "mongodbatlas_advanced_cluster" "automated_backup_test" {
  project_id   = mongodbatlas_project.atlas-project.id
  name         = var.cluster_name
  cluster_type = "REPLICASET"

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
      }

      provider_name = "GCP"
      region_name   = "US_EAST_1"
      priority      = 7
    }]
  }]

  backup_enabled = true # enable cloud backup snapshots
}

resource "mongodbatlas_data_lake_pipeline" "test" {
  project_id = mongodbatlas_project.atlas-project.id
  name       = var.name
  sink {
    type = "DLS"
    partition_fields {
      field_name = "access"
      order      = 0
    }
  }

  source {
    type            = "ON_DEMAND_CPS"
    cluster_name    = mongodbatlas_advanced_cluster.automated_backup_test.name
    database_name   = "sample_airbnb"
    collection_name = "listingsAndReviews"
  }

  transformations {
    field = "testField"
    type  = "EXCLUDE"
  }

  transformations {
    field = "testField2"
    type  = "EXCLUDE"
  }

}
