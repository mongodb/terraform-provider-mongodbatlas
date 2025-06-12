# Example - MongoDB Atlas Backup Compliance Policy
This example shows how to configure the `mongodbatlas_backup_compliance_policy` and the [lifecycle impact](#backup-compliance-policy-impact-on-a-cluster-with-cloud-backup-schedule) on the `mongodbatlas_advanced_cluster` and `mongodbatlas_cloud_backup_schedule`. With Backup Compliance Policy enabled, cluster backups are retained after a cluster is deleted (remember to set `retainBackups=true`) and backups can be used until retention expiration.

For more details see [Back Up, Restore, and Archive Data](https://www.mongodb.com/docs/atlas/backup-restore-cluster/)

## Backup Compliance Policy Deletion Consideration
-> **IMPORTANT NOTE:** Once a Backup Compliance Policy is enabled, no user, regardless of role, can disable the Backup Compliance Policy via Terraform, or any other method, without contacting MongoDB Support. However, if the Backup Compliance Policy is approved for disablement by MongoDB Support, the resource can be removed and the delete operation will successfully disable the Backup Compliance Policy. With Backup Compliance Policy enabled, some resources defined in Terraform will no longer be modifiable. See the full list of [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).

Therefore, the resource `mongodbatlas_backup_compliance_policy.this` defined in [main.tf](main.tf) must only be used if the intention is to use backup compliance policies for the project.

## Backup Compliance Policy Impact on a Cluster with Cloud Backup Schedule
The delete call to `mongodbatlas_cloud_backup_schedule` can fail with `BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS` when `mongodbatlas_backup_compliance_policy` is active.

The error happens because Terraform tries to delete the `mongodbatlas_cloud_backup_schedule` resource **before** deleting the cluster and the `mongodbatlas_backup_compliance_policy` resource blocks the action because the `mongodbatlas_advanced_cluster` will no longer be compliant with the policy (for more details, see the [Configure a Backup Compliance Policy](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy) docs).

The Terraform [resource dependency](https://developer.hashicorp.com/terraform/language/resources/behavior#resource-dependencies) between `mongodbatlas_advanced_cluster` and `mongodbatlas_cloud_backup_schedule` is reversed during deletion.
By following the steps below you will see how to avoid this error.


## Dependencies

- Terraform MongoDB Atlas Provider
- A MongoDB Atlas account

**Required** Variables to be set:
- `user_id`: Unique 24-hexadecimal digit string that identifies this user.
- `project_id`: Atlas Project ID.

**Optional** Variables to be set:
- `public_key`: Atlas Programmatic API public key
- `private_key`: Atlas Programmatic API private key
- `cluster_name`: Name of the cluster
- `instance_size`: Instance size of the cluster.

## Usage

**Note**: This directory contains an example of using the **Preview for MongoDB Atlas Provider 2.0.0** of `mongodbatlas_advanced_cluster`. In order to enable the Preview, you must set the environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`, otherwise the current version will be used.

You can find more info in the [resource doc page](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster%2520%2528preview%2520provider%25202.0.0%2529).


### 1. Ensure your MongoDB Atlas credentials are set up

This can be done using environment variables or the `public_key` and `private_key` variables.

```bash
export MONGODB_ATLAS_PUBLIC_KEY="<ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<ATLAS_PRIVATE_KEY>"
```
### 2. Review the Terraform plan

Execute the following command and ensure you are happy with the plan.


```bash

terraform plan
```

### 3. Execute the Terraform apply

Now execute the plan to provision the resources.

```bash
terraform apply
```

### 4. How to delete a cluster when a Backup Compliance Policy Is Enabled

To proceed with the deletion you can choose one of the methods below:  

#### 1. (Recommended) Use a `removed` block to avoid the DELETE call for `mongodbatlas_cloud_backup_schedule`
**Note**: If you are using a Terraform Module, we recommend using method 3 or follow the [module example](../module/README.md#how-to-delete-the-cluster-and-retain-their-backup-snapshots).

This method requires Terraform CLI [1.7 or later](https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources).

Add the removed block
```terraform
removed {
  from = mongodbatlas_cloud_backup_schedule.this

  lifecycle {
    destroy = false
  }
}
```
Remove the `resource "mongodbatlas_cloud_backup_schedule" "this"` definition.

Run `terraform apply`. You should see a plan similar to:

```bash
# mongodbatlas_cloud_backup_schedule.this will no longer be managed by Terraform, but will not be destroyed
 # (destroy = false is set in the configuration) [...]
 ```

 Run `terraform destroy`


#### 2. Disable `mongodbatlas_backup_compliance_policy` by contacting MongoDB Support
You will need to reach out to MongoDB Support and complete a thorough verification process. However, this process may not be ideal for most customers who wish to maintain their backup compliance policies for other active clusters.

#### 3. Use `terraform state rm` to remove `mongodbatlas_cloud_backup_schedule` from the state to avoid the DELETE call for `mongodbatlas_cloud_backup_schedule`
Note: This is identical to method 1 but requires access to `terraform state rm`.

1. Run `terraform state rm mongodbatlas_cloud_backup_schedule.this`.
2. Remove the `resource "mongodbatlas_cloud_backup_schedule" "this"` block.
3. Run `terraform destroy`
