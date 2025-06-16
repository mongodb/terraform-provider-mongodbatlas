---
page_title: "Guide: Delete a MongoDB Atlas Cluster with an active Backup Compliance Policy"
---

# Guide: Deleting a MongoDB Atlas Cluster with an Active Backup Compliance Policy

**Objective**: This guide explains how to delete a MongoDB Atlas Cluster using
Terraform when a Backup Compliance Policy (BCP) is enabled. It also explores the
relationship between the Terraform resources:
`mongodbatlas_backup_compliance_policy`, `mongodbatlas_cloud_backup_schedule`,
and `mongodbatlas_advanced_cluster`.

## Why Would You Need a Backup Compliance Policy?

A Backup Compliance Policy is essential if you have strict data protection
requirements. Enabling this policy ensures your backup data remains protected.
You can learn more about Backup Compliance Policy and its implications in the
[official documentation](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/).

In Terraform, the `mongodbatlas_backup_compliance_policy` resource
([documentation](../resources/backup_compliance_policy.md)) is used to configure
this feature.

## How Does a Backup Compliance Policy Impact Terraform Configuration?

If your Terraform configuration includes both a MongoDB Atlas cluster and its
associated backup schedule, it will look something like this:

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "clusterTest"
  
  ...
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_advanced_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_advanced_cluster.my_cluster.name
  
  ...
}

...
```

In this setup, `mongodbatlas_advanced_cluster` and
`mongodbatlas_cloud_backup_schedule` are defined within the same Terraform
project, with a direct dependency between the resources. If you’re using a
Terraform module, these resources might be included within that module.

If a Backup Compliance Policy is enabled for your MongoDB Atlas project
(configured via Terraform or another tool), there is an important restriction.
As stated
[here](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#prohibited-actions),
**you cannot modify the backup policy for an individual cluster below the
minimum requirements outlined in the Backup Compliance Policy**.

This means any attempt to remove the `mongodbatlas_cloud_backup_schedule` from
your Terraform configuration will trigger an error, specifically:
`BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS`. This error
occurs because the cluster's backup policy would fall **below the minimum
requirements set**. While this behavior is expected, it leads to an important
question:

**"How can I delete my cluster using Terraform while preserving my backup
snapshots?"**

## Steps to Delete a MongoDB Atlas Cluster with BCP Enabled and Retain Snapshots

To delete a MongoDB Atlas cluster in this scenario, you need to follow a
two-step process. This approach aligns with the requirements of your enabled
Backup Compliance Policy.

- **Step 1: Update Terraform to Ignore the `mongodbatlas_cloud_backup_schedule`
  Configuration**: before deleting the cluster, instruct Terraform to "ignore"
  the `mongodbatlas_cloud_backup_schedule` resource to avoid violating the
  Backup Compliance Policy.

- **Step 2: Delete the MongoDB Atlas Cluster with Terraform**: once the
  `mongodbatlas_cloud_backup_schedule` is removed from Terraform's scope,
  proceed with the cluster deletion as usual while ensuring backup snapshots
  remain intact.

To assist with implementation, we’ve provided two examples:

1. **Using Resources Directly**\
   View the
   [resource-based example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_backup_compliance_policy/resource).

2. **Using Terraform Modules**\
   Review the
   [module-based example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_backup_compliance_policy/module).

These examples outline the adjustments required for each approach to
successfully delete clusters under the constraints of a Backup Compliance
Policy.
