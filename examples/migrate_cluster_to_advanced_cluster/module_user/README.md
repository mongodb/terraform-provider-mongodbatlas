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
* Configure the provider (can also be done with `variables.tfvars`)

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

## Usage

**Configure `variables.tfvars`:**

```tfvars
project_id             = "664619d870c247237f4b86a6"
cluster_name           = "module-cluster"
cluster_type           = "SHARDED"
instance_size          = "M10"
mongo_db_major_version = "8.0"
provider_name          = "AWS"
disk_size              = 40
tags = {
  env    = "examples"
  module = "cluster_to_advanced_cluster"
}
replication_specs = [
  {
    num_shards = 2
    zone_name  = "Zone 1"
    regions_config = [
      {
        region_name     = "US_EAST_1"
        electable_nodes = 3
        priority        = 7
        read_only_nodes = 0
      }
    ]
  }
]
```


### `v1`

```bash
cd v1
terraform init
terraform apply -var-file=variables.tfvars
```

### `v2`

```bash
cd v2
export MONGODB_ATLAS_ADVANCED_CLUSTER_V2_SCHEMA=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply -var-file=variables.tfvars # notice the same variables used as in `v1`
```

### `v3`

**Configure `variables-updated.tfvars`:**

```tfvars
# TODO: Use the new variables for latest features
```

```bash
cd v3
export MONGODB_ATLAS_ADVANCED_CLUSTER_V2_SCHEMA=true # necessary for the `moved` block to work
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply -var-file=variables-updated.tfvars # updated variables to enable latest mongodb_advanced_cluster features
```

### Cleanup with `terraform destroy`

```bash
terraform destroy
```
