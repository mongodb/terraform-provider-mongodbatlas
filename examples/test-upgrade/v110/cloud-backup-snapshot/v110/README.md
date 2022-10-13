# MongoDB Atlas Provider -- Cloud Backup Snapshot
This example creates a project, cluster, cloud provider snapshot, and a restore job for said snapshot. The cluster is configured to use cloud backup and point in time restore.

Variables Required:
- `org_id`: ID of atlas organization
- `project_name`: Name of the project
- `cluster_name`: Name of the cluster
- `point_in_time_utc_seconds`: Point in time to restore to, a number of seconds since unix epoch.

Example `terraform.tfvars`
```
org_id                        = "627a9687f7f7f7f774de306f14"
project_name                  = "cloud_backup_snapshot_v110"
cluster_name                  = "v110-cluster"
point_in_time_utc_seconds     = 1665549600
```