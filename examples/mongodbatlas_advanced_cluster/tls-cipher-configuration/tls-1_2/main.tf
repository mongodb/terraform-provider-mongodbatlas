provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = mongodbatlas_project.project.id
  name         = "TLS12Configuration"
  cluster_type = "SHARDED"

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = var.provider_instance_size_name
        disk_iops     = 3000
        node_count    = 3
        disk_size_gb  = 10
      }
      provider_name = var.provider_name
      priority      = 7
      region_name   = "US_EAST_1"
    }]
    },
    {
      region_configs = [{
        electable_specs = {
          instance_size = var.provider_instance_size_name
          disk_iops     = 3000
          node_count    = 3
          disk_size_gb  = 10
        }
        provider_name = var.provider_name
        priority      = 7
        region_name   = "US_EAST_1"
      }]
  }]
  advanced_configuration = {
    custom_openssl_cipher_config_tls12 = [
      "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
    ]
    minimum_enabled_tls_protocol = "TLS1_2"
    tls_cipher_config_mode       = "CUSTOM"
  }

}

resource "mongodbatlas_project" "project" {
  name   = "TLS12Configuration"
  org_id = var.atlas_org_id
}
