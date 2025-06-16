---
page_title: "Guide: Delete a MongoDB Atlas Cluster with an active Backup Compliance Policy"
---

# Guide: Delete a MongoDB Atlas Cluster with an active Backup Compliance Policy

**Objective**: In this guide we explain how to delete a MongoDB Atlas Cluster
using Terraform when a Backup Compliance Policy (BCP) is enabled and how the
`mongodbatlas_backup_compliance_policy`, `mongodbatlas_cloud_backup_schedule`
and `mongodbatlas_advanced_cluster` terraform resources are related to each
other.

## When do I need a Backup Compliance Policy?

If you have strict data protection requirements, you can enable a Backup
Compliance Policy to protect your backup data. You can read more about Backup
Compliance Policy and its considerations
[at this link](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/).

In Terraform, you can use the `mongodbatlas_backup_compliance_policy` resource
([documentation](../resources/backup_compliance_policy.md)) to configure it.

## How does a Backup Compliance Policy affect my Terraform configuration?

If you have defined a
[MongoDB Atlas Cluster](../resources/advanced_cluster%20(preview%20provider%202.0.0).md)
and its [backup schedule](../resources/cloud_backup_schedule.md) in a Terraform
configuration, your configuration will look like something like this:

```
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

meaning that you have both resources defined in the same terraform project with
a direct dependency between each other (If you are using a module, those
resources are likely defined as part of the module).

At the same time, if you've enabled a Backup Compliance Policy for your MongoDB
Atlas Project (either via Terraform or via any other tool), one of the
consequences as described
[here](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#prohibited-actions)
is that you cannot modify the
[backup policy](https://www.mongodb.com/docs/atlas/backup/cloud-backup/configure-backup-policy/#std-label-configure-backup-policy)
for an individual cluster **below the minimum requirements set** in the Backup
Compliance Policy.

What this means is that any tentative to delete the
`mongodbatlas_cloud_backup_schedule` from this terraform configuration would
result into an error
(`BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS`) given the
related cluster's backup policy would end up being **below the minimum
requirements set**. While this error is expected, it poses the following
question: "How should I proceed with the deletion of my cluster using Terraform
by retaining the backup snapshot?"

## How do I delete a MongoDB Atlas Cluster with a Backup Compliance Policy enabled and retain backups?

High level, in order to delete a MongoDB Atlas Cluster in this scenario is by
following a two-step procedure. This is an expected operation to take, given
once again your MongoDB Atlas Project has an enabled Backup Compliance Policy.

- **Step 1.** Tell Terraform to "ignore" the
  `mongodbatlas_cloud_backup_schedule` configuration
- **Step 2.** Proceed with the deletion of MongoDB Atlas Cluster via Terraform

In order to explain how to do this, we've created two examples: one related to
the direct usage of
[resources](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_backup_compliance_policy/resource),
and one related to
[module](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_backup_compliance_policy/module)
usage.
