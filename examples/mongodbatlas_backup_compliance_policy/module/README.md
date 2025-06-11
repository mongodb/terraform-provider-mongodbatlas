# Example - MongoDB Atlas Backup Compliance Policy with a module
This example is identical to the [resource example](../resource/README.md) except that it uses a local [module](modules/cluster_with_schedule/main.tf) to manage the `mongodbatlas_advanced_cluster` and the `mongodbatlas_cloud_backup_schedule`.

It shows a workaround for deleting the `mongodbatlas_advanced_cluster` and `mongodbatlas_cloud_backup_schedule` when using Terraform Modules. See more details about the `removed` block limitation in the Hashicorp Terraform [issue 34439](https://github.com/hashicorp/terraform/issues/34439).

## Different cleanup when using the `removed` block for a module
When using a module it is likely you cannot remove the resource block from the module config directly, `module.source` can be somewhere else; therefore, you get the `Removed Resource still exists error` when trying to use the `removed` block:

```terraform
removed {
  from = module.cluster_with_schedule.mongodbatlas_cloud_backup_schedule.this

  lifecycle {
    destroy = false
  }
}
```

To workaround this limitation we can create another module instance (`cluster_without_schedule`) in our root [main.tf](main.tf) and use a `moved` block for the cluster:

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
