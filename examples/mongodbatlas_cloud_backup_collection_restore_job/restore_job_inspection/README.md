# Collection restore job inspection

Read back collection restore jobs on a cluster without creating a new job.

Use this example after a restore apply, or against any cluster that already has collection restore jobs. To start a restore, use [`snapshot_restore/`](../snapshot_restore/README.md) or [`pit_restore/`](../pit_restore/README.md).

## Prerequisites

- Backup Manager or Project Owner on the source project.
- At least one existing collection restore job on the cluster when `job_id` is omitted, or set `job_id` explicitly.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Defaults

Omit optional variables to use these values from `variables.tf`:

- `job_id`: unset — uses the last `job_id` from `mongodbatlas_cloud_backup_collection_restore_jobs`
- `source_namespace`: unset — uses the last `source_namespace` from `mongodbatlas_cloud_backup_collection_restore_job_collections`

The singular collection data source returns no state while the job is `INITIALIZING`.

## Usage

**1. Set credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**2. Create `terraform.tfvars`.**

```hcl
project_id   = "your-prod-project-id"
cluster_name = "your-prod-cluster-name"

# job_id            = "6789abcd1234ef5678901234"  # default: last job when unset
# source_namespace  = "sample_mflix.movies"       # default: last namespace when unset
```

**3. Apply and read outputs.**

```bash
terraform apply
```

- `job_id` — resolved job ID
- `job_state` — job state from the singular job data source
- `collection_states` — per-collection state from the plural collections data source
- `collection_state` — state for one collection from the singular collection data source
- `restore_jobs` — trimmed list of jobs (`job_id`, `state`, `created_at`)
