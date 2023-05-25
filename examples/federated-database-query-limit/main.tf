resource "mongodbatlas_cluster" "atlas_cluster_1" {
  project_id                  = var.project_id
  provider_name               = var.provider_name
  name                        = var.atlas_cluster_name_1
  backing_provider_name       = var.backing_provider_name
  provider_region_name        = var.provider_region_name
  provider_instance_size_name = var.provider_instance_size_name
}


resource "mongodbatlas_cluster" "atlas_cluster_2" {
  project_id                  = var.project_id
  provider_name               = var.provider_name
  name                        = var.atlas_cluster_name_2
  backing_provider_name       = var.backing_provider_name
  provider_region_name        = var.provider_region_name
  provider_instance_size_name = var.provider_instance_size_name
}

resource "mongodbatlas_federated_database_instance" "test-instance" {
  project_id = var.project_id
  name       = var.federated_instance_name
  aws {
    role_id        = ""
    test_s3_bucket = ""
  }
  storage_databases {
    name = "VirtualDatabase0"
    collections {
      name = "VirtualCollection0"
      data_sources {
        collection = var.collection_1
        database   = var.database_1
        store_name = mongodbatlas_cluster.atlas_cluster_1.name
      }
      data_sources {
        collection = var.collection_2
        database   = var.database_2
        store_name = mongodbatlas_cluster.atlas_cluster_2.name
      }
    }
  }

  storage_stores {
    name         = mongodbatlas_cluster.atlas_cluster_1.name
    cluster_name = mongodbatlas_cluster.atlas_cluster_1.name
    project_id   = var.project_id
    provider     = "atlas"
    read_preference {
      mode = "secondary"
    }
  }

  storage_stores {
    name         = mongodbatlas_cluster.atlas_cluster_2.name
    cluster_name = mongodbatlas_cluster.atlas_cluster_2.name
    project_id   = var.project_id
    provider     = "atlas"
    read_preference {
      mode = "secondary"
    }
  }
}

resource "mongodbatlas_federated_query_limit" "query_limit" {
  project_id     = var.project_id
  tenant_name    = mongodbatlas_federated_database_instance.test-instance.name
  limit_name     = var.federated_query_limit
  overrun_policy = var.overrun_policy
  value          = var.limit_value
}
