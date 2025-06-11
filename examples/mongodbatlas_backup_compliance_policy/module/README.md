# Example - MongoDB Atlas Backup Compliance Policy with a module
This example is identical to the [resource example](../resource/README.md) except that it uses a local [module](modules/cluster_with_schedule/main.tf) to manage the `mongodbatlas_advanced_cluster` and the `mongodbatlas_cloud_backup_schedule` via the `add_schedule` variable.
The [cleanup step below](#different-cleanup-when-using-the-removed-block-for-a-module) shows how the `moved` and `removed` block can remove the `mongodbatlas_cloud_backup_schedule` from state to avoid the [BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS](../resource/README.md#4-cleanup-extra-steps-when-a-backup-compliance-policy-is-enabled) error when running `terraform destroy`.

## Different cleanup when using the `removed` block for a module

**Note**: If you can use the `terraform state rm` command, follow the simpler steps in the [resource example readme](../resource/README.md#3-use-terraform-state-rm-to-remove-mongodbatlas_cloud_backup_schedule-from-the-state-to-avoid-the-delete-call-for-mongodbatlas_cloud_backup_schedule).

When using a module, modifying the module terraform files can be inaccessible:
- The `module.source` can be pointing at an online repository, for example [terraform-mongodbatlas-atlas-basic](https://github.com/terraform-mongodbatlas-modules/terraform-mongodbatlas-atlas-basic)
- If you can modify the module terraform files, follow the steps in the [resource example readme](../resource/README.md#1-recommended-use-a-removed-block-to-avoid-the-delete-call-for-mongodbatlas_cloud_backup_schedule)

If you try to use the `removed` block without deleting the `from` resouce you get the error: `Removed Resource still exists error`:

```terraform
removed {
  from = module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this

  lifecycle {
    destroy = false
  }
}
```

To workaround this limitation (see more details about the `removed` block limitation in the Hashicorp Terraform [issue 34439](https://github.com/hashicorp/terraform/issues/34439)) we can create another module instance (`cluster_without_schedule`) in our root [main.tf](main.tf) and use a `moved` block for the cluster:

```terraform
module "cluster_without_schedule" {
    source = "./modules/cluster_with_schedule"

    project_id = var.project_id
    instance_size = var.instance_size
    cluster_name = var.cluster_name
    add_schedule = false # changed flag
}
# Keep the cluster
moved {
  from = module.cluster_with_schedule.mongodbatlas_advanced_cluster.this
  to  = module.cluster_without_schedule.mongodbatlas_advanced_cluster.this
}
```

And comment/remove the old module (`cluster_with_schedule`) instance:
```terraform
# Removed or Commented out
# module "cluster_with_schedule" {
#   source = "./modules/cluster_with_schedule"

#   project_id    = var.project_id
#   instance_size = var.instance_size
#   cluster_name  = var.cluster_name
#   add_schedule  = true
# }
```

Then when we run:
```sh
terraform init # initialize the cluster_without_schedule instance
terraform apply
```

We should see:

```bash
[...]
module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this will no longer be managed by Terraform, but will not be destroyed
[...]
module.cluster_with_schedule.mongodbatlas_advanced_cluster.this has moved to module.cluster_without_schedule.mongodbatlas_advanced_cluster.this
[...]
Plan: 0 to add, 0 to change, 0 to destroy.
```

Reply `yes` to confirm the move and state removal of `mongodbatlas_cloud_backup_schedule`.

Then run `terraform destroy` to destroy the:
- Cluster defined in `module.cluster_without_schedule.mongodbatlas_advanced_cluster.this`.
- Root compliance policy defined in `mongodbatlas_backup_compliance_policy`.
