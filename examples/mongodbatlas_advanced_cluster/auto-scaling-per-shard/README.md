# MongoDB Atlas Provider -- Sharded Cluster with Independent Shard Auto-scaling

This example creates a Sharded Cluster with 2 shards defining electable and analytics nodes. Compute auto-scaling is enabled for both `electable_specs` and `analytics_specs`, while also leveraging the [New Sharding Configuration](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema) by defining each shard with its individual `replication_specs`. This enables scaling of each shard to be independent. Please reference the [Use Auto-Scaling Per Shard](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema#use-auto-scaling-per-shard) section for more details.


## Dependencies

* Terraform MongoDB Atlas Provider v1.23.0
* A MongoDB Atlas account 

```
Terraform >= 0.13
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.23.0
```


## Usage
**1\. If you haven't already, set up your MongoDB Atlas credentials.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```
public_key           = "<MONGODB_ATLAS_PUBLIC_KEY>"
private_key          = "<MONGODB_ATLAS_PRIVATE_KEY>"
atlas_org_id         = "<MONGODB_ATLAS_ORG_ID>"
```

Alternatively, you can use [AWS Secrets Manager](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/website/docs/index.html.markdown#aws-secrets-manager).

**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- An Atlas Project
- A Sharded Cluster with independent shards with varying cluster tiers

**3\. Apply your changes.**

Now execute the plan to provision the Atlas Project and Cluster resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```

