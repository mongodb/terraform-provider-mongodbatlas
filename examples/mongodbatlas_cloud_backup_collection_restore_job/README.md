# MongoDB Atlas Collection Restore Examples

Restore selected databases and collections from Cloud Backup into an existing cluster without wiping everything else.

Restore examples default to `write_strategy = "CREATE_NEW"` and `index_strategy = "ALL"` in `variables.tf` when those variables are omitted. Optional `database_renames`, `collection_renames`, `database_suffix`, and `collection_suffix` control target namespace naming. With `OVERWRITE_EXISTING`, Atlas drops and replaces matching namespaces on the target; stop writers to those namespaces before apply to avoid losing writes made during the restore.

The example does not create a cluster. Point it at a cluster that already has Cloud Backup enabled.

For product limits, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_backup_collection_restore_job#limitations) and [Restore from Selected Databases and Collections](https://www.mongodb.com/docs/atlas/backup/cloud-backup/restore-from-db-coll/).

## Sibling examples

Read-only:

- [`snapshot_discovery/`](snapshot_discovery/README.md) — list snapshots, resolve `snapshot_id`, browse databases and collections in a snapshot. Run this before a snapshot restore when you need to pick inputs.
- [`restore_job_inspection/`](restore_job_inspection/README.md) — read back collection restore jobs on a cluster without creating a new job.

Create:

- [`snapshot_restore/`](snapshot_restore/README.md) — restore from a scheduled or on-demand snapshot.
- [`pit_restore/`](pit_restore/README.md) — restore from continuous backup at a point in time. No snapshot ID.

Typical snapshot flow: `snapshot_discovery` → copy outputs into `snapshot_restore` tfvars → apply.
