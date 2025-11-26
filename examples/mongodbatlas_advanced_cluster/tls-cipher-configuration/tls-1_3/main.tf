provider "mongodbatlas" {
  client_id     = var.atlas_client_id
  client_secret = var.atlas_client_secret
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = mongodbatlas_project.project.id
  name         = "TLS13Configuration"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = var.provider_instance_size_name
            node_count    = 3
          }
          provider_name = var.provider_name
          region_name   = "US_EAST_1"
          priority      = 7
        }
      ]
    }
  ]
  advanced_configuration = {
    custom_openssl_cipher_config_tls13 = [
      "TLS_AES_256_GCM_SHA384"
    ]
    minimum_enabled_tls_protocol = "TLS1_3"
    tls_cipher_config_mode       = "CUSTOM"
  }

}

resource "mongodbatlas_project" "project" {
  name   = "TLS13Configuration"
  org_id = var.atlas_org_id
}
