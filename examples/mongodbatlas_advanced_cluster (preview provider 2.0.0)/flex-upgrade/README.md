# MongoDB Atlas Provider -- Flex cluster (Preview for MongoDB Atlas Provider 2.0.0)

This example creates a project and a Flex cluster using `mongodbatlas_advanced_cluster` resource. It is intended to show how to create a Flex cluster, upgrade an M0 cluster to Flex and upgrade a Flex cluster to a Dedicated cluster.

It uses the **Preview for MongoDB Atlas Provider 2.0.0** of `mongodbatlas_advanced_cluster`. In order to enable the Preview, you must set the enviroment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`, otherwise the current version will be used.

You can find more information in the [resource documentation page](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster%2520%2528preview%2520provider%25202.0.0%2529).

Variables Required:
- `atlas_org_id`: ID of the Atlas organization
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `provider_name`: Name of provider to use for cluster (TENANT, AWS, GCP)
- `backing_provider_name`: If provider_name is tenant, the backing provider (AWS, GCP)
- `provider_instance_size_name`: Size of the cluster (Shared: M0. Dedicated: M10+.)

For this example, first we'll start out on the Free tier, then upgrade to a flex cluster and finally to a Dedicated tier cluster.

Utilize the following to execute a working example, replacing the org id, public and private key with your values:

Apply with the following `terraform.tfvars` to first create a free tier cluster:
```
atlas_org_id                = <YOUR_ORG_ID>
public_key                  = <REDACTED>
private_key                 = <REDACTED>
provider_name               = "TENANT"
backing_provider_name       = "AWS"
provider_instance_size_name = "M0"
node_count 					= null
```

Apply with the following `terraform.tfvars` to upgrade the free tier cluster you just created to flex tier:
```
atlas_org_id                = <YOUR_ORG_ID>
public_key                  = <REDACTED>
private_key                 = <REDACTED>
provider_name               = "FLEX"
backing_provider_name       = "AWS"
provider_instance_size_name = null
node_count 					= null
```
Apply with the following `terraform.tfvars` to upgrade the flex tier cluster you just created to dedicated tier:
```
atlas_org_id                = <YOUR_ORG_ID>
public_key                  = <REDACTED>
private_key                 = <REDACTED>
provider_name               = "AWS"
backing_provider_name       = null
provider_instance_size_name = "M10"
node_count 					= 3
```
