resource "mongodbatlas_project" "atlas-project" {
  org_id = var.atlas_org_id
  name   = var.atlas_project_name
}

resource "mongodbatlas_cluster" "automated_backup_test" {
  project_id                  = mongodbatlas_project.atlas-project.id
  name                        = var.cluster_name
  provider_name               = "GCP"
  provider_region_name        = "US_EAST_4"
  provider_instance_size_name = "M30"
  cloud_backup                = true // enable cloud backup snapshots
  mongo_db_major_version      = "4.4"
  disk_size_gb                = "350"
}

resource "mongodbatlas_data_lake_pipeline" "test" {
  project_id = mongodbatlas_project.atlas-project.id
  name       = var.cluster_name
  sink {
    type = "DLS"
    partition_fields {
      name  = "access"
      order = 0
    }
  }

  source {
    type            = "ON_DEMAND_CPS"
    cluster_name    = mongodbatlas_cluster.automated_backup_test.name
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

data "mongodbatlas_data_lake_pipelines" "testDataSource" {
  project_id = mongodbatlas_data_lake_pipeline.test.project_id
}


