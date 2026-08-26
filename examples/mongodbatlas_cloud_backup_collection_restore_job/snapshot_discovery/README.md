# Snapshot namespace discovery

List Cloud Backup snapshots for a cluster, resolve a `snapshot_id`, and browse databases and collections in that snapshot.

Run this example from `snapshot_discovery/` before [`snapshot_restore/`](../snapshot_restore/README.md) when you need to pick a snapshot or namespace names. Point-in-time restores have no snapshot ID; use [`pit_restore/`](../pit_restore/README.md) instead.

Apply lists all databases and collections in the selected snapshot via the snapshot data sources.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Defaults

Omit optional variables to use these values from `variables.tf`:

- `snapshot_id`: unset — uses the latest completed snapshot for the cluster
- `available_snapshots_limit`: `5` — newest snapshots included in `available_snapshots`

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

# snapshot_id = "6789abcd1234ef5678901234"  # default: latest completed snapshot when unset
```

**3. Apply and read outputs.**

```bash
terraform apply
```

- `available_snapshots` — newest completed snapshots for the cluster (up to `available_snapshots_limit`)
- `snapshot_id` — snapshot used for database and collection listing
- `snapshot_details` — `created_at`, databases, and collections for that snapshot

Copy `snapshot_id`, database names, and collection namespaces from `snapshot_details` into [`snapshot_restore/terraform.tfvars`](../snapshot_restore/README.md).
