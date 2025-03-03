# Basic Migration from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`

This example demonstrates how to migrate a `mongodbatlas_cluster` resource to `mongodbatlas_advanced_cluster` (see alternatives, and more details in the [cluster to advanced cluster migration guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide)).
In this example we use specific files, but the same approach can be applied to any configuration file with `mongodbatlas_cluster` resource(s).
The main steps are:

1. [Create the `mongodbatlas_cluster`](#create-the-mongodbatlas_cluster) (skip if you already have a configuration with one or more `mongodbatlas_cluster` resources).
2. [Use the Atlas CLI Plugin Terraform to create the `mongodbatlas_advanced_cluster` configuration](#use-the-atlas-cli-plugin-terraform-to-create-the-mongodbatlas_advanced_cluster-resource).
3. [Manually update the Terraform configuration](#manual-updates-to-the-terraform-configuration).
4. [Perform the Move](#perform-the-move).
   - [Troubleshooting](#troubleshooting).

## Create the `mongodbatlas_cluster`

**Note**: This step is only to demonstrate the migration. If you want to manage a cluster with Terraform, we recommend you use a [mongodbatlas_advanced_cluster](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster%2520%2528preview%2520provider%2520v2%2529) resource instead.

This step can be skipped if you already have a configuration file with a `mongodbatlas_cluster` created.

1. Uncomment the code in [outputs.tf](outputs.tf) (lines marked with `# BEFORE`) and [cluster.tf](cluster.tf).
2. Comment the code in [outputs.tf](outputs.tf) (lines marked with `# AFTER`) and [advanced_cluster.tf](advanced_cluster.tf).
3. Create a `vars.auto.tfvars` file, for example:
```terraform
project_id = "{PROJECT_ID}" # replace with your project ID, should be similar to 664619d870c247237f4b86a6
cluster_name = "cluster-mig-resource"
instance_size = "M10"
mongo_db_major_version = "8.0"
```
4. Run `terraform init`
5. Run `terraform apply`

## Use the Atlas CLI Plugin Terraform to create the `mongodbatlas_advanced_cluster` resource

The [CLI Plugin](https://github.com/mongodb-labs/atlas-cli-plugin-terraform) helps you convert your existing configuration into one using the `mongodbatlas_advanced_cluster` resource. This step is not required but recommended; alternatively you can write the equivalent configuration on your own.
1. Ensure you have the Atlas CLI [installed](https://www.mongodb.com/docs/atlas/cli/current/install-atlas-cli/) by running `atlas --version`.
   - You should see a multi line output starting with `atlascli version: 1.38.0` (or a later version).
2. Install the Terraform CLI plugin: `atlas plugin install mongodb-labs/atlas-cli-plugin-terraform`.
3. Run `atlas terraform` to ensure the plugin was installed correctly.
   - You should see a multi line output starting with `Utilities for Terraform's MongoDB Atlas Provider`.
4. Run `atlas tf clusterToAdvancedCluster --file {CLUSTER_IN}.tf --output {CLUSTER_OUT}.tf`.
   1. For example `cluster.tf` for `{CLUSTER_IN}.tf`.
   2. For example `advanced_cluster.tf` for `{CLUSTER_OUT}.tf`.
   3. If your config is not supported, see [limtations](https://github.com/mongodb-labs/atlas-cli-plugin-terraform?tab=readme-ov-file#limitations) of the plugin, try to update your `mongodbatlas_cluster`, or see alternatives in the [Migration Guide: Cluster to Advanced Cluster](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide).

## Manual updates to the Terraform configuration

1. Ensure all references are updated (see example of updates in [outputs.tf](outputs.tf))
2. Comment out the `mongodbatlas_cluster` in `{CLUSTER_IN}.tf`
3. Add the moved block for each resource migrated in `{CLUSTER_OUT}.tf`
```terraform
moved {
  from = mongodbatlas_cluster.this # change `this` to your specific resource identifier
  to   = mongodbatlas_advanced_cluster.this # change `this` to your specific resource identifier
}
```
- The [moved block](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring#moved-block-syntax) can be kept to record where you have historically moved or renamed an object.

## Perform the Move

1. Ensure you are using V2 schema: `export MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`.
2. Run `terraform validate` to ensure there are no missing reference updates. You might see errors like:
   - `Error: Reference to undeclared resource`: You forgot to update the resource type to `mongodbatlas_advanced_cluster`
   ```text
    │   on outputs.tf line 7, in output "container_id":
    │    7:     value = mongodbatlas_cluster.this.replication_specs[0].container_id
   ```
   - `Error: Unsupported attribute`:  The attribute is no longer supported in `mongodbatlas_advanced_cluster` and needs to be removed.
   ```text
    │   on outputs.tf line 7, in output "provider_name":
    │    7:   value = mongodbatlas_advanced_cluster.this.provider_name
    │ 
    │ This object has no argument, nested block, or exported attribute named "provider_name".
   ```
3. Run `terraform apply` and accept the move.
   - You should expect to see
   ```text
   Terraform will perform the following actions:
    # mongodbatlas_cluster.this has moved to mongodbatlas_advanced_cluster.this
        resource "mongodbatlas_advanced_cluster" "this" {
            name                                 = "cluster-mig-resource"
            # (24 unchanged attributes hidden)
        }

    Plan: 0 to add, 0 to change, 0 to destroy.
    Do you want to perform these actions?
        Terraform will perform the actions described above.
        Only 'yes' will be accepted to approve.

        Enter a value:
   ```
   - Type `yes` and hit enter

### Troubleshooting

- You might see: `Changes to Outputs`, consider where the output is used and take action accordingly.
```text
container_id               = "67a09ae45cc3a60e55b6f42d" -> "67ac794392f9196661de88e1"
```
- If there are any Plan Changes, try updating the `mongodbatlas_advanced_cluster` in `{CLUSTER_OUT}.tf` manually.
- For example the below plan, would require you to explicitly set `backup_enabled = false` in the `mongodbatlas_advanced_cluster.this` resource.
```text
Terraform will perform the following actions:

  # mongodbatlas_advanced_cluster.this will be updated in-place
  # (moved from mongodbatlas_cluster.this)
  ~ resource "mongodbatlas_advanced_cluster" "this" {
      ~ backup_enabled                       = false -> true
      ~ connection_strings                   = {
          + private          = (known after apply)
          + private_endpoint = (known after apply)
          + private_srv      = (known after apply)
          ~ standard         = "mongodb://cluster-mig-resource-shard-00-00.jciib.mongodb-dev.net:27017,cluster-mig-resource-shard-00-01.jciib.mongodb-dev.net:27017,cluster-mig-resource-shard-00-02.jciib.mongodb-dev.net:27017,cluster-mig-resource-shard-00-03.jciib.mongodb-dev.net:27017,cluster-mig-resource-shard-00-04.jciib.mongodb-dev.net:27017,cluster-mig-resource-shard-00-05.jciib.mongodb-dev.net:27017,cluster-mig-resource-shard-00-06.jciib.mongodb-dev.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-46mqxd-shard-0" -> (known after apply)
          ~ standard_srv     = "mongodb+srv://cluster-mig-resource.jciib.mongodb-dev.net" -> (known after apply)
        } -> (known after apply)
        name                                 = "cluster-mig-resource"
        # (22 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```
- You might see errors like:
- `Error: Invalid index`: This is due to type changes and can usually be resolved by removing the `[0]` or `.0` reference.
```text
│   on outputs.tf line 3, in output "connection_string_standard":
│    3:     value = mongodbatlas_advanced_cluster.this.connection_strings[0].standard
│     ├────────────────
│     │ mongodbatlas_advanced_cluster.this.connection_strings is object with 5 attributes
│
│ The given key does not identify an element in this collection value. An object only supports looking up attributes by name, not by numeric index.
```
