# Module User - `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`

The purpose of this example is to demonstrate the experience of adopting a new version of a terraform module definition which internally migrated from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`.
Each module call represent a step on the migration path.
The example focus on the call of the module rather than the module implementation itself (see the [module maintainer README.md](../module_maintainer/README.md) for the implementation details).

Migration Step Code | `mongodbatlas` version | Config Changes | Plan Changes
--- | --- | --- | ---
[Step 1](./v1) | `<= 1.27.0` | Baseline configuration | -
[Step 2](./v2) | `>= 1.27.0` | Only change to `v2` module version | No changes. Only moved.
[Step 3](./v3) | `>= 1.27.0` | Usage of new variables to support multi-cloud and independent shard scaling | Yes (new features)


The rest of this example is a step by step guide on how to migrate from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`:

- [Dependencies](#dependencies)
- [Step 1: Create the `mongodbatlas_cluster` with `v1` of the module](#step-1-create-the-mongodbatlas_cluster-with-v1-of-the-module)
  - [Update `v1_v2.tfvars`](#update-v1_v2tfvars)
  - [Run Commands](#run-commands)
- [Step 2: Upgrade to `mongodbatlas_advanced_cluster` by using `v2` of the module](#step-2-upgrade-to-mongodbatlas_advanced_cluster-by-using-v2-of-the-module)
- [Step 3: Use the latest `mongodbatlas_advanced_cluster` features by using `v3` of the module](#step-3-use-the-latest-mongodbatlas_advanced_cluster-features-by-using-v3-of-the-module)
  - [Update `v3_no_plan_changes`](#update-v3_no_plan_changes)
  - [Run `terraform plan` to ensure there are no plan changes](#run-terraform-plan-to-ensure-there-are-no-plan-changes)
  - [Update `v3.tfvars`](#update-v3tfvars)
  - [Run `terraform apply` to upgrade the cluster to an Asymmetric Sharded Cluster](#run-terraform-apply-to-upgrade-the-cluster-to-an-asymmetric-sharded-cluster)
- [Cleanup with `terraform destroy`](#cleanup-with-terraform-destroy)

## Dependencies
- Terraform CLI >= 1.8-
- Terraform MongoDB Atlas Provider `>=v1.27.0`-
- A MongoDB Atlas account.
- Configure the provider (can also be done by configuring `public_key` and `private_key` in a `provider.tfvars`).

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

## Step 1: Create the `mongodbatlas_cluster` with `v1` of the module

### Update `v1_v2.tfvars`

See the example in [v1_v2.tfvars](v1_v2.tfvars).

### Run Commands
```bash
cd v1
terraform init
terraform apply -var-file=../v1_v2.tfvars
```

## Step 2: Upgrade to `mongodbatlas_advanced_cluster` by using `v2` of the module

```bash
cd v2
cp ../v1/terraform.tfstate . # if you are not using a remote state
export MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply -var-file=../v1_v2.tfvars # notice the same variables used as in `v1`
```
In the plan output, you should see a line simlar to:
```text
# module.cluster.mongodbatlas_cluster.this has moved to module.cluster.mongodbatlas_advanced_cluster.this
```

## Step 3: Use the latest `mongodbatlas_advanced_cluster` features by using `v3` of the module

### Update `v3_no_plan_changes`

See the example in [`v3_no_plan_changes`](v3_no_plan_changes.tfvars).
The example changes:
1. Use the new `replication_specs_new` variable.
2. Remove old `replication_specs`, `provider_name`, `instance_size`, `disk_size` variables.

### Run `terraform plan` to ensure there are no plan changes

```bash
cd v3
cp ../v2/terraform.tfstate . # if you are not using a remote state
export MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform plan -var-file=../v3_no_plan_changes.tfvars # updated variables to enable latest mongodb_advanced_cluster features
```

### Update `v3.tfvars`

See the example in [v3.tfvars](v3.tfvars)
The example changes:
1. Increase `disk_size_gb` from `40` -> `50`
2. New `read_only_specs`
3. Increased `instance_size`:
   1. In shard 1 from `M10` -> `M30`
   1. In shard 2 from `M10` -> `M50`

### Run `terraform apply` to upgrade the cluster to an Asymmetric Sharded Cluster

```bash
cd v3
cp ../v2/terraform.tfstate . # if you are not using a remote state
export MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply -var-file=../v3.tfvars # updated variables to enable latest mongodb_advanced_cluster features
```

## Cleanup with `terraform destroy`

```bash
terraform destroy
```
