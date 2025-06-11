# Example - MongoDB Atlas Backup Compliance Policy with a module
This example is identical to the [resource example](../resource/README.md) except that it uses a local [module](modules/cluster_with_schedule/main.tf) to manage the `mongodbatlas_advanced_cluster` and the `mongodbatlas_cloud_backup_schedule` via the `add_schedule` variable.
The [cleanup step below](#different-cleanup-when-using-the-removed-block-for-a-module) shows how the `moved` and `removed` block can remove the `mongodbatlas_cloud_backup_schedule` from your Terraform state to avoid the [BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS](../resource/README.md#4-cleanup-extra-steps-when-a-backup-compliance-policy-is-enabled) error when running `terraform destroy`.

## Different cleanup when using the `removed` block for a module

**Note**: If you can use the `terraform state rm` command or edit the module TF files directly, follow the simpler steps in the [resource example readme](../resource/README.md#4-cleanup-extra-steps-when-a-backup-compliance-policy-is-enabled).

If you try to use the `removed` block without deleting the `from` resource you get the error: `Removed Resource still exists error`:

```terraform
removed {
  from = module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this

  lifecycle {
    destroy = false
  }
}
```

To workaround this limitation (see more details about the `removed` block limitation in the Hashicorp Terraform [issue 34439](https://github.com/hashicorp/terraform/issues/34439)) you can use a `moved` block and change the `add_schedule` flag in the root [main.tf](main.tf):

```terraform
module "cluster_without_schedule" {
    source = "./modules/cluster_with_schedule"

    project_id = var.project_id
    instance_size = var.instance_size
    cluster_name = var.cluster_name
    add_schedule = false # changed flag
}
# Rename the resource to avoid the `Removed Resource still exists error`
moved {
  from = module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this[0]
  to  = mongodbatlas_cloud_backup_schedule.deleted # any resource name that doesn't exist works!
}

removed {
  from = mongodbatlas_cloud_backup_schedule.deleted

  lifecycle {
    destroy = false
  }
}
```
Then when you run:
```sh
terraform init # initialize the cluster_without_schedule instance
terraform apply
```

You should see:

```bash
[...]
module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this[0] will no longer be managed by Terraform, but will not be destroyed
[...]
Plan: 0 to add, 0 to change, 0 to destroy.
```

Reply `yes` to confirm the state removal of `mongodbatlas_cloud_backup_schedule`.

Then run `terraform destroy` to destroy the:
- Cluster defined in `module.cluster_with_schedule.mongodbatlas_advanced_cluster.this`.
- Root compliance policy defined in `mongodbatlas_backup_compliance_policy`.
