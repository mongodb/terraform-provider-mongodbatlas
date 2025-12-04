# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

# Set use_effective_fields = true on the resource
resource "mongodbatlas_advanced_cluster" "this" {
  project_id           = mongodbatlas_project.this.id
  name                 = var.cluster_name
  cluster_type         = var.cluster_type
  use_effective_fields = true
  replication_specs    = var.replication_specs
  tags                 = var.tags
}

/*
 Phase 1 (current, backward compatible):
 - Omit use_effective_fields flag (defaults to false) on the data source
 - *_specs (electable_specs, analytics_specs, read_only_specs) return actual provisioned values
 - Recommended for migrating from lifecycle.ignore_changes while maintaining compatibility with module_existing behavior

 Phase 2 (breaking change, prepares for v3.x):
 - Set use_effective_fields = true on the data source
 - *_specs return configured values, effective_*_specs return actual provisioned values
 - Module users must switch from *_specs to effective_*_specs for actual values
 - Recommended for new modules or when preparing for provider v3.x
*/
data "mongodbatlas_advanced_cluster" "this" {
  project_id = mongodbatlas_advanced_cluster.this.project_id
  name       = mongodbatlas_advanced_cluster.this.name
  # use_effective_fields = true  # Uncomment for Phase 2
  depends_on = [mongodbatlas_advanced_cluster.this]
}
