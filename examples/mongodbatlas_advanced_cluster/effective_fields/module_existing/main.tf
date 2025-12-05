# Create Atlas Project
resource "mongodbatlas_project" "this" {
  name   = var.project_name
  org_id = var.atlas_org_id
}

/*
 This module simulates an existing implementation using lifecycle.ignore_changes blocks.

 This was the only available approach before use_effective_fields was introduced and is no longer
 recommended. This implementation is included only for understanding migration from existing modules.
 For new modules, use the module_effective_fields approach instead.
*/
resource "mongodbatlas_advanced_cluster" "this" {
  project_id        = mongodbatlas_project.this.id
  name              = var.cluster_name
  cluster_type      = var.cluster_type
  replication_specs = var.replication_specs
  tags              = var.tags

  /*
   When auto-scaling is enabled (either compute or disk auto-scaling), Atlas may adjust
   instance_size, disk_size_gb, and disk_iops regardless of which auto-scaling type is enabled.
   Therefore, all three attributes must be ignored to prevent unintended changes.

   This lifecycle block ignores changes to instance_size, disk_size_gb, and disk_iops
   for each node type (electable, analytics, read_only) in up to 2 replication_specs
   and 3 region_configs per replication_spec.

   This approach has limitations:
   - Cannot be conditional based on auto-scaling configuration
   - Module users cannot see actual provisioned values
   - Requires listing all auto-scalable attributes for maximum expected topology
  */
  lifecycle {
    ignore_changes = [
      # Replication Spec 0 - Region 0
      replication_specs[0].region_configs[0].electable_specs.instance_size,
      replication_specs[0].region_configs[0].electable_specs.disk_size_gb,
      replication_specs[0].region_configs[0].electable_specs.disk_iops,
      replication_specs[0].region_configs[0].analytics_specs.instance_size,
      replication_specs[0].region_configs[0].analytics_specs.disk_size_gb,
      replication_specs[0].region_configs[0].analytics_specs.disk_iops,
      replication_specs[0].region_configs[0].read_only_specs.instance_size,
      replication_specs[0].region_configs[0].read_only_specs.disk_size_gb,
      replication_specs[0].region_configs[0].read_only_specs.disk_iops,

      # Replication Spec 0 - Region 1
      replication_specs[0].region_configs[1].electable_specs.instance_size,
      replication_specs[0].region_configs[1].electable_specs.disk_size_gb,
      replication_specs[0].region_configs[1].electable_specs.disk_iops,
      replication_specs[0].region_configs[1].analytics_specs.instance_size,
      replication_specs[0].region_configs[1].analytics_specs.disk_size_gb,
      replication_specs[0].region_configs[1].analytics_specs.disk_iops,
      replication_specs[0].region_configs[1].read_only_specs.instance_size,
      replication_specs[0].region_configs[1].read_only_specs.disk_size_gb,
      replication_specs[0].region_configs[1].read_only_specs.disk_iops,

      # Replication Spec 0 - Region 2
      replication_specs[0].region_configs[2].electable_specs.instance_size,
      replication_specs[0].region_configs[2].electable_specs.disk_size_gb,
      replication_specs[0].region_configs[2].electable_specs.disk_iops,
      replication_specs[0].region_configs[2].analytics_specs.instance_size,
      replication_specs[0].region_configs[2].analytics_specs.disk_size_gb,
      replication_specs[0].region_configs[2].analytics_specs.disk_iops,
      replication_specs[0].region_configs[2].read_only_specs.instance_size,
      replication_specs[0].region_configs[2].read_only_specs.disk_size_gb,
      replication_specs[0].region_configs[2].read_only_specs.disk_iops,

      # Replication Spec 1 - Region 0
      replication_specs[1].region_configs[0].electable_specs.instance_size,
      replication_specs[1].region_configs[0].electable_specs.disk_size_gb,
      replication_specs[1].region_configs[0].electable_specs.disk_iops,
      replication_specs[1].region_configs[0].analytics_specs.instance_size,
      replication_specs[1].region_configs[0].analytics_specs.disk_size_gb,
      replication_specs[1].region_configs[0].analytics_specs.disk_iops,
      replication_specs[1].region_configs[0].read_only_specs.instance_size,
      replication_specs[1].region_configs[0].read_only_specs.disk_size_gb,
      replication_specs[1].region_configs[0].read_only_specs.disk_iops,

      # Replication Spec 1 - Region 1
      replication_specs[1].region_configs[1].electable_specs.instance_size,
      replication_specs[1].region_configs[1].electable_specs.disk_size_gb,
      replication_specs[1].region_configs[1].electable_specs.disk_iops,
      replication_specs[1].region_configs[1].analytics_specs.instance_size,
      replication_specs[1].region_configs[1].analytics_specs.disk_size_gb,
      replication_specs[1].region_configs[1].analytics_specs.disk_iops,
      replication_specs[1].region_configs[1].read_only_specs.instance_size,
      replication_specs[1].region_configs[1].read_only_specs.disk_size_gb,
      replication_specs[1].region_configs[1].read_only_specs.disk_iops,

      # Replication Spec 1 - Region 2
      replication_specs[1].region_configs[2].electable_specs.instance_size,
      replication_specs[1].region_configs[2].electable_specs.disk_size_gb,
      replication_specs[1].region_configs[2].electable_specs.disk_iops,
      replication_specs[1].region_configs[2].analytics_specs.instance_size,
      replication_specs[1].region_configs[2].analytics_specs.disk_size_gb,
      replication_specs[1].region_configs[2].analytics_specs.disk_iops,
      replication_specs[1].region_configs[2].read_only_specs.instance_size,
      replication_specs[1].region_configs[2].read_only_specs.disk_size_gb,
      replication_specs[1].region_configs[2].read_only_specs.disk_iops,
    ]
  }
}

/*
 Data source to read actual cluster state.
 With lifecycle.ignore_changes, the resource's replication_specs attributes return
 values from state which may differ from configuration due to ignored changes.
 Use this data source to get the real provisioned values from Atlas API.
*/
data "mongodbatlas_advanced_cluster" "this" {
  project_id = mongodbatlas_advanced_cluster.this.project_id
  name       = mongodbatlas_advanced_cluster.this.name
  depends_on = [mongodbatlas_advanced_cluster.this]
}
