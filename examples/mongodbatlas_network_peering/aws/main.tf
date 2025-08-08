provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

resource "mongodbatlas_project" "aws_atlas" {
  name   = "aws-atlas"
  org_id = var.atlas_org_id
}

resource "mongodbatlas_advanced_cluster" "cluster-atlas" {
  project_id     = mongodbatlas_project.aws_atlas.id
  name           = "cluster-atlas"
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AWS"
      region_name   = var.atlas_region
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
}

resource "mongodbatlas_database_user" "db-user" {
  username           = var.atlas_dbuser
  password           = var.atlas_dbpassword
  auth_database_name = "admin"
  project_id         = mongodbatlas_project.aws_atlas.id
  roles {
    role_name     = "readWrite"
    database_name = "admin"
  }
  depends_on = [mongodbatlas_project.aws_atlas]
}

resource "mongodbatlas_network_peering" "aws-atlas" {
  accepter_region_name     = var.aws_region
  project_id               = mongodbatlas_project.aws_atlas.id
  container_id             = one(values(mongodbatlas_advanced_cluster.cluster-atlas.replication_specs[0].container_id))
  provider_name            = "AWS"
  route_table_cidr_block   = aws_vpc.primary.cidr_block
  vpc_id                   = aws_vpc.primary.id
  aws_account_id           = var.aws_account_id
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}

resource "mongodbatlas_project_ip_access_list" "test" {
  project_id = mongodbatlas_project.aws_atlas.id
  cidr_block = aws_vpc.primary.cidr_block
  comment    = "cidr block for AWS VPC"
}
