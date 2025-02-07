# Module - `cluster` to `advanced_cluster`

The purpose of this example is to demonstrate the upgrade path from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster` using a Terraform module.
The example contains three module versions:

Version | Purpose | Variables | Resources
--- | --- | --- | ---
[v1](./v1) | Baseline | 5 | `mongodbatlas_cluster`
[v2](./v2) | Migrate to advanced_cluster with no change in variables or plan | 5 | `mongodbatlas_advanced_cluster`
[v3](./v3) | Use the latest features of advanced_cluster | 10 | `mongodbatlas_advanced_cluster`

<!-- TODO: Update the actual Variable counts  -->

## Dependencies
<!-- TODO: Update XX versions -->
<!-- TODO: Update the `versions.tf` inside each vX -->
* Terraform CLI >= 1.X
* Terraform MongoDB Atlas Provider v1.XX.0
* A MongoDB Atlas account
* Configure the provider (can also be done with `variables.tfvars`)

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

* Configure `variables.tfvars`
<!-- TODO: Example variables -->
```tfvars
# todo: add example variable declaration here
```

## Usage

### `v1`

```bash
# uncomment the code in main.tf marked with v1, ensure v2 and v3 is commented
export 
terraform init
terraform apply -var-file=variables.tfvars
```

### `v2`

```bash
# uncomment the code in main.tf marked with v2, ensure v1 and v3 is commented
export MONGODB_ATLAS_ADVANCED_CLUSTER_V2_SCHEMA=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply -var-file=variables.tfvars
```

### `v3`

```bash
# uncomment the code in main.tf marked with v3, ensure v1 and v2 is commented
export MONGODB_ATLAS_ADVANCED_CLUSTER_V2_SCHEMA=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply -var-file=variables.tfvars
```

### Cleanup with `terraform destroy`

```bash
terraform destroy
```
