# MongoDB Atlas Cloud Provider Access Setup for Azure
resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = var.project_id
  provider_name = "AZURE"
  
  azure_config {
    atlas_azure_app_id   = var.azure_atlas_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

# MongoDB Atlas Cloud Provider Access Authorization
resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = var.project_id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

  azure {
    atlas_azure_app_id   = var.azure_atlas_app_id
    service_principal_id = var.azure_service_principal_id
    tenant_id            = var.azure_tenant_id
  }
}

# MongoDB Atlas Federated Database Instance with Azure
resource "mongodbatlas_federated_database_instance" "azure_example" {
  project_id = var.project_id
  name       = var.federated_instance_name

  # Azure cloud provider configuration
  cloud_provider_config {
    azure {
      role_id = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
    }
  }

  # Minimal storage configuration using only Atlas cluster
  storage_databases {
    name = "VirtualDatabase0"
    collections {
      name = "VirtualCollection0"
      data_sources {
        store_name = "azure_cluster_store"
        database   = var.database_name
        collection = var.collection_name
      }
    }
  }

  storage_stores {
    name         = "azure_cluster_store"
    provider     = "atlas"
    cluster_name = var.cluster_name
    project_id   = var.project_id
    read_preference {
      mode = "secondary"
    }
  }
}