locals {
  mongodb_uri = mongodbatlas_cluster.this.connection_strings[0].standard
}

data "mongodbatlas_federated_settings" "this" {
  org_id = var.org_id
}
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.org_id
  tags   = local.tags
}

resource "mongodbatlas_project_ip_access_list" "mongo-access" {
  project_id = mongodbatlas_project.this.id
  cidr_block = "0.0.0.0/0"
}

resource "mongodbatlas_cluster" "this" {
  project_id             = mongodbatlas_project.this.id
  name                   = var.project_name
  mongo_db_major_version = "7.0"
  cluster_type           = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = var.region
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  cloud_backup                 = false
  auto_scaling_disk_gb_enabled = false
  provider_name                = "AWS"
  disk_size_gb                 = 10
  provider_instance_size_name  = "M10"
}

resource "mongodbatlas_federated_settings_identity_provider" "oidc" {
  federation_settings_id = data.mongodbatlas_federated_settings.this.id
  audience               = var.token_audience
  authorization_type     = "USER"
  description            = "oidc-for-azure"
  # e.g. "https://sts.windows.net/91405384-d71e-47f5-92dd-759e272cdc1c/"
  issuer_uri = "https://sts.windows.net/${azurerm_user_assigned_identity.this.tenant_id}/"
  idp_type   = "WORKLOAD"
  name       = "OIDC-for-azure"
  protocol   = "OIDC"
  # groups_claim = null
  user_claim = "sub"
}

resource "mongodbatlas_federated_settings_org_config" "this" {
  federation_settings_id            = data.mongodbatlas_federated_settings.this.id
  org_id                            = var.org_id
  domain_restriction_enabled        = false
  domain_allow_list                 = []
  data_access_identity_provider_ids = [mongodbatlas_federated_settings_identity_provider.oidc.idp_id]
}

resource "mongodbatlas_database_user" "oidc" {
  project_id         = mongodbatlas_project.this.id
  username           = "${mongodbatlas_federated_settings_identity_provider.oidc.idp_id}/${azurerm_user_assigned_identity.this.principal_id}"
  oidc_auth_type     = "USER"
  auth_database_name = "$external" # required when using OIDC USER authentication

  roles {
    role_name     = "atlasAdmin"
    database_name = "admin"
  }
}
