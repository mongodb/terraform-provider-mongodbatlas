# MongoDB Atlas Provider -- Global Cluster (Preview for MongoDB Atlas Provider v2)

This example creates a project and a Sharded Cluster with 4 independent shards with varying cluster tiers.

It uses the **Preview for MongoDB Atlas Provider v2** of `mongodbatlas_advanced_cluster`. In order to enable the Preview, you must set the enviroment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`, otherwise the current version will be used.

You can find more information in the [resource documentation page](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster%2520%2528preview%2520provider%2520v2%2529).

## Dependencies

* Terraform MongoDB Atlas Provider v1.29.0
* A MongoDB Atlas account 

```
Terraform >= 0.13
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.29.0
```


## Usage
**1\. Ensure your MongoDB Atlas credentials are set up.**

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

... or use [AWS Secrets Manager](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/website/docs/index.html.markdown#aws-secrets-manager)

**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- An Atlas Project
- A Sharded Cluster with independent shards with varying cluster tiers

**3\. Execute the Terraform apply.**

Now execute the plan to provision the Atlas Project and Cluster resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```

