# Example - MongoDB Atlas Backup Compliance Policy with a module
This example module is identical to the [resource example](../resource/README.md) except that it is designed to serve as a reference for platform teams who have created their own module and make it available to internal teams for leveraging the MongoDB Atlas Terraform provider. Typically, these users do not have the ability to execute `terraform state` commands or modify the Terraform state manually.

As in the resource example, the attention is focused on how to avoid the [BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS](../resource/README.md#4-cleanup-extra-steps-when-a-backup-compliance-policy-is-enabled) when running `terraform destroy` on a `mongodbatlas_advanced_cluster` with `mongodbatlas_cloud_backup_schedule` and enabled backup compliance policy.

To do that, we'll use:
- an `add_schedule` variable provided by the module maintainer that controls the presence of the `mongodbatlas_cloud_backup_schedule` in the configuration
- the `moved` and `removed` blocks used by the module user and added at the root level.

-> **IMPORTANT NOTE:** Read the [Backup Compliance Policy Deletion Consideration](../resource/README.md#backup-compliance-policy-deletion-consideration) before running this example.

## How to delete the cluster and retain their backup snapshots

Let's begin by assuming you've created a module instance with the following configuration:

```terraform
module "cluster_with_schedule" {
  source = "./modules/cluster_with_schedule"

  project_id    = var.project_id
  instance_size = var.instance_size
  cluster_name  = var.cluster_name
  add_schedule  = true
}
```
Note: The `add_schedule` field is set to true, indicating that a `mongodbatlas_cloud_backup_schedule` resource has been defined, as reflected in the module's source code.

To proceed with the deletion, we'll update the configuration as follows:

```terraform
module "cluster_with_schedule" {
  source = "./modules/cluster_with_schedule"

  project_id    = var.project_id
  instance_size = var.instance_size
  cluster_name  = var.cluster_name
  add_schedule  = false # changed
}

moved {
  from = module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this[0] # must be deleted with the `add_schedule` variable set to false
  to   = mongodbatlas_cloud_backup_schedule.to_be_deleted                              # any resource name that doesn't exist works!
}

removed {
  from = mongodbatlas_cloud_backup_schedule.to_be_deleted # any resource name that doesn't exist works!

  lifecycle {
    destroy = false
  }
}
```

Then when you run `terraform apply`, you should see:

```bash
[...]
mongodbatlas_cloud_backup_schedule.to_be_deleted will no longer be managed by Terraform, but will not be destroyed
 (destroy = false is set in the configuration)
  (moved from module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this[0])
[...]
Plan: 0 to add, 0 to change, 0 to destroy.
```

Reply `yes` to confirm the state removal of `mongodbatlas_cloud_backup_schedule`.

Then run `terraform destroy` to destroy the cluster defined in `module.cluster_with_schedule.mongodbatlas_advanced_cluster.this`.
See [Backup Compliance Policy Deletion Consideration](../resource/README.md#backup-compliance-policy-deletion-consideration) for details on `mongodbatlas_backup_compliance_policy` deletion.

## FAQ
I get a `Removed Resource still exists error` when running `terraform apply`, how do I fix it?

This error happens because the configuration still has the `mongodbatlas_cloud_backup_schedule` active.
Remember to add the `moved` block and set `add_schedule=false` on the `cluster_with_schedule` module.
