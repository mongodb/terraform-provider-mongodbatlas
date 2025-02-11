# Module User - `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`

The purpose of this example is to demonstrate the User experience when upgrading `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`.

Module Version | `mongodbatlas` version | Config Changes | Plan Changes
--- | --- | --- | ---
[v1](./v1) | <= PLACEHOLDER_TPF_RELEASE | Baseline configuration | -
[v2](./v2) | >= PLACEHOLDER_TPF_RELEASE | Only change to `v2` module version | No changes. Only moved.
[v3](./v3) | >= PLACEHOLDER_TPF_RELEASE | Usage of new variables to support multi-cloud and independent shard scaling | Yes (new features)


## Dependencies
<!-- TODO: Update XX versions -->
<!-- TODO: Update the `versions.tf` inside each vX -->
* Terraform CLI >= 1.X
* Terraform MongoDB Atlas Provider v1.XX.0
* A MongoDB Atlas account
* Configure the provider (can also be done by configuring `public_key` and `private_key` in a `provider.tfvars`)

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

## Usage
* These steps show example of calling the modules in [../module_maintainer](../module_maintainer/).
* To follow the steps and use the same cluster you need to use a remote state or copy the `terraform.tfstate` file when switching from `vx` to `vy` (e.g., `v1` to `v2`)

### Update `v1_v2.tfvars`
See the example in [v1_v2.tfvars](v1_v2.tfvars)

### `v1`

```bash
cd v1
terraform init
terraform apply -var-file=../v1_v2.tfvars
```

### `v2`

```bash
cd v2
export MONGODB_ATLAS_ADVANCED_CLUSTER_V2_SCHEMA=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply -var-file=../v1_v2.tfvars # notice the same variables used as in `v1`
```

### Update `v3.tfvars`
See the example in [v3.tfvars](v3.tfvars)

### `v3`

```bash
cd v3
export MONGODB_ATLAS_ADVANCED_CLUSTER_V2_SCHEMA=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply -var-file=../v3.tfvars # updated variables to enable latest mongodb_advanced_cluster features
```

### Cleanup with `terraform destroy`

```bash
terraform destroy
```
