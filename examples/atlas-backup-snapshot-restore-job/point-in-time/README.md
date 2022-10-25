# MongoDB Atlas Provider -- Cloud Backup Snapshot
This example creates a project, cluster, cloud provider snapshot, and a restore job for said snapshot. The cluster is configured to use cloud backup and point in time restore.

Variables Required:
- `org_id`: ID of atlas organization
- `project_name`: Name of the project
- `cluster_name`: Name of the cluster
- `point_in_time_utc_seconds`: Point in time to restore to, a number of seconds since unix epoch.

In order to utilize the backup restore job via point in time, fist you need a backup with which to restore.
This example has been configured to only create the backup restore if `point_in_time_utc_seconds` is a non-zero number.
As such, utilize the following example `terraform.tfvars` and pseudo-code to execute a workign example:

Example `terraform.tfvars`
```
org_id                        = "627a9687f7f7f7f774de306f14"
project_name                  = "cloud_backup_snapshot_v110"
cluster_name                  = "v110-cluster"
point_in_time_utc_seconds     = 0
```

- Run `terraform apply`
- Update `point_in_time_utc_seconds` to the [current epoch time](https://www.epoch101.com/)
- Run `terraform apply`

You'll now have a project, cluster, backup snapshot, and restore job pointing to specific point in time which to restore.