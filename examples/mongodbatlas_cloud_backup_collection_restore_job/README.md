# MongoDB Atlas Collection Restore Example

This example restores one collection from a Cloud Backup snapshot into an existing cluster. It also reads snapshot namespaces and per-collection restore state.

The example does not create a cluster. Point it at a cluster that already has Cloud Backup and a snapshot that contains `database_name`.

Credentials come from the environment (`MONGODB_ATLAS_PUBLIC_KEY` and `MONGODB_ATLAS_PRIVATE_KEY`, or `MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`). Do not put keys in Terraform files.

## Prerequisites

- MongoDB Atlas Service Account or API Key with Backup Manager or Project Owner on the source project, and on the target project if they differ.
- An existing dedicated cluster with Cloud Backup enabled and at least one snapshot.
- The source namespace must exist in that snapshot.

Cross-organization restore is UI-only. Programmatic credentials are org-scoped. Cross-project restore inside one org is supported: set `target_project_id` and `target_cluster_name`.

## Resources Created

This example creates the following resources:

- MongoDB Atlas collection restore job

Create waits until the job reaches `SUCCESSFUL`. Default timeout is 3 hours. Atlas keeps running if Terraform times out. There is no public cancel API. `terraform destroy` only removes the job from state.

## Local provider binary

From the provider repository root, build the binary and point Terraform at `dev.tfrc` (gitignored). See [Development Overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers).

```bash
go build -o bin/terraform-provider-mongodbatlas .
export TF_CLI_CONFIG_FILE="$PWD/dev.tfrc"
```

`dev.tfrc` must list the directory that contains `terraform-provider-mongodbatlas`, not the binary path.

## Usage

**1. Set credentials and the source snapshot.**

```bash
export MONGODB_ATLAS_PUBLIC_KEY="<ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<ATLAS_PRIVATE_KEY>"
```

Create a **terraform.tfvars** file (gitignored) with the snapshot to restore:

```hcl
project_id       = "your-project-id"
cluster_name     = "your-cluster-name"
snapshot_id      = "your-snapshot-id"
database_name    = "sample_mflix"
source_namespace = "sample_mflix.movies"
```

For a point-in-time restore, replace `snapshot_id` on the resource with `point_in_time_utc_seconds` (and optional `oplog_ts` / `oplog_inc`). Do not set both snapshot and PIT fields.

**2. Review the Terraform plan.**

```bash
terraform plan
```

**3. Apply the restore.**

```bash
terraform apply
```

Apply can take up to the create timeout. `OVERWRITE_EXISTING` drops matching namespaces on the target. Stop writers to those namespaces first.

**4. Destroy the Terraform resource.**

```bash
terraform destroy
```

Destroy does not cancel or delete the Atlas job. It only removes the resource from state.
