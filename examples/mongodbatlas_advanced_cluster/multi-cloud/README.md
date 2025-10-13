# MongoDB Atlas Provider -- Multi-Cloud Advanced Cluster

This example creates a project and a Multi Cloud Advanced Cluster with 2 shards.

## Dependencies

* Terraform MongoDB Atlas Provider v2.0.0 or later
* A MongoDB Atlas account 

```
Terraform >= 0.13
+ provider registry.terraform.io/terraform-providers/mongodbatlas v2.0.0
```


## Usage
**1\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="<ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<ATLAS_PRIVATE_KEY>"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```
public_key           = "<MONGODB_ATLAS_PUBLIC_KEY>"
private_key          = "<MONGODB_ATLAS_PRIVATE_KEY>"
atlas_org_id         = "<MONGODB_ATLAS_ORG_ID>"
```

... or use [AWS Secrets Manager](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/index.md#aws-secrets-manager)

**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- An Atlas Project
- A Multi-Cloud Cluster

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
