# Create a Project
resource "mongodbatlas_project" "atlas-project" {
  org_id = var.atlas_org_id
  name = var.atlas_project_name
}

# Create a Database User
resource "mongodbatlas_database_user" "db-user" {
  username = "user-1"
  password = random_password.db-user-password.result
  project_id = mongodbatlas_project.atlas-project.id
  auth_database_name = "admin"
  roles {
    role_name     = "readWrite"
    database_name = "${var.atlas_project_name}-db"
  }
}

# Create a Database Password
resource "random_password" "db-user-password" {
  length = 16
  special = true
  override_special = "_%@"
  sensitive = true
}

# Create Database IP Access List 
resource "mongodbatlas_project_ip_access_list" "ip" {
  project_id = mongodbatlas_project.atlas-project.id
  ip_address = var.ip_address
}

# Create an Atlas Advanced Cluster 
resource "mongodbatlas_advanced_cluster" "atlas-cluster" {
  project_id = mongodbatlas_project.atlas-project.id
  name = "${var.atlas_project_name}-${var.environment}-cluster"
  cluster_type = "REPLICASET"
  backup_enabled = true
  mongo_db_major_version = var.mongodb_version
  replication_specs {
    region_configs {
      electable_specs {
        instance_size = var.cluster_instance_size_name
        node_count    = 3
      }
      analytics_specs {
        instance_size = var.cluster_instance_size_name
        node_count    = 1
      }
      priority      = 7
      provider_name = var.cloud_provider
      region_name   = var.atlas_region
    }
  }
}

data "mongodbatlas_advanced_cluster" "atlas-cluser" {
  project_id = mongodbatlas_project.atlas-project.id
  name       = mongodbatlas_advanced_cluster.atlas-cluster.name
  depends_on = [mongodbatlas_privatelink_endpoint_service.atlaseplink]
}

# Outputs to Display
output "connection_string" { value = mongodbatlas_advanced_cluster.atlas-cluster.connection_strings.0.standard_srv }
output "ip_access_list"    { value = mongodbatlas_project_ip_access_list.ip.ip_address }
output "project_name"      { value = mongodbatlas_project.atlas-project.name }
output "username"          { value = mongodbatlas_database_user.db-user.username } 
output "user_password"     { 
  sensitive = true
  value = mongodbatlas_database_user.db-user.password 
  }
output "plstring" {
  value = lookup(mongodbatlas_advanced_cluster.atlas-cluster.connection_strings[0].aws_private_link_srv, aws_vpc_endpoint.ptfe_service.id)
}
