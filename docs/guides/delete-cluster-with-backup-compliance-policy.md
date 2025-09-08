---
page_title: "Guide: Delete a MongoDB Atlas Cluster with an active Backup Compliance Policy"
subcategory: "Older Guides"
---

# Guide: Delete a MongoDB Atlas Cluster with an Active Backup Compliance Policy

**Objective**: Learn how to delete a MongoDB Atlas Cluster using Terraform when
a Backup Compliance Policy (BCP) is enabled and how the following Terraform
resources are related to each other: `mongodbatlas_backup_compliance_policy`,
`mongodbatlas_cloud_backup_schedule`, and `mongodbatlas_advanced_cluster`.

## Why Do You Need a Backup Compliance Policy?

You must use a Backup Compliance policy if you have strict data protection
requirements. Enabling this policy ensures your backup data remains protected.
To learn more, see
[Backup Compliance Policy](https://mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/)
in the MongoDB Atlas Documentation.

To configure a Backup Compliance Policy in Terraform, use the
[mongodbatlas_backup_compliance_policy](../resources/backup_compliance_policy.md)
resource.

## How Does a Backup Compliance Policy Impact Terraform Configuration?

If a Backup Compliance Policy is enabled for your MongoDB Atlas project
(configured via Terraform or another tool), the following
[actions are prohibited](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#prohibited-actions),
leading to some implications for Terraform.

Let's assume your configuration is similar to the following:

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = "<PROJECT-ID>"
  name           = "clusterTest"
  
  ...
}

resource "mongodbatlas_cloud_backup_schedule" "this" {
  project_id   = mongodbatlas_advanced_cluster.this.project_id
  cluster_name = mongodbatlas_advanced_cluster.this.name
  
  ...
}

...
```

In this setup, `mongodbatlas_advanced_cluster` and
`mongodbatlas_cloud_backup_schedule` are defined within the same Terraform
project, with a direct dependency between the resources. If you’re using a
Terraform module, these resources might be included within that module.

When you attempt to run a `terraform destroy` on a configuration similar to the
above, as expected Terraform will delete resources in the inverse order of
dependency, starting from `mongodbatlas_cloud_backup_schedule`. Attempting to
delete the backup schedule results in the following error:
`BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS`. In fact, If
this deletion was successful, it would leave the individual cluster below the
minimum requirements outlined in the Backup Compliance Policy. While this
behavior is expected, it leads to the important question that we are going to
cover in the following section: **How can you delete your cluster using
Terraform when a Backup Compliance Policy is enabled?**.

## Steps to Delete a MongoDB Atlas Cluster with BCP Enabled and Retain Snapshots

To delete a MongoDB Atlas cluster in this scenario, follow a two-step process.
This approach aligns with the requirements of your enabled Backup Compliance
Policy.

- **Step 1: Update Terraform to remove `mongodbatlas_cloud_backup_schedule` from
  the state**. Before deleting the cluster, instruct Terraform to "ignore" the
  `mongodbatlas_cloud_backup_schedule` resource to avoid violating the Backup
  Compliance Policy.

- **Step 2: Delete the MongoDB Atlas Cluster with Terraform**. Once you remove
  the `mongodbatlas_cloud_backup_schedule` from Terraform's state, proceed with
  deleting the cluster with `terraform destroy`.

Use the following examples to assist with deleting a cluster. These examples
outline the adjustments required for each approach to successfully delete
clusters under the constraints of a Backup Compliance Policy.

1. **Using Resources Directly**\
   [resource-based example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_backup_compliance_policy/resource).

2. **Using Terraform Modules**\
   [module-based example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_backup_compliance_policy/module).
