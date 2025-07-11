# Atlas Terraform Provider Example: Encryption at rest - AWS - Multi-Region Cluster

This example sets up encryption at rest using AWS KMS for your Atlas Project. It creates the encryption key in AWS KMS, an IAM role and policy so that Atlas can access the key, and enables encryption at rest for the Atlas Project. Finally, it creates a 2-shards Multi-Region Cluster with encryption at rest enabled.

## Dependencies

* Terraform MongoDB Atlas Provider v1.10.2
* A MongoDB Atlas account 

```
Terraform >= 0.13
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.2
```



## Usage

**1\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="<ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<ATLAS_PRIVATE_KEY>"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```hcl
public_key           = "<PUBLIC_KEY>"
private_key          = "<PRIVATE_KEY>"
atlas_org_id         = "<ATLAS_ORG_ID>"
```

... or use [AWS Secrets Manager](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/index.md#aws-secrets-manager)


**2\. Set your AWS access key & secret via environment variables:

```bash
export AWS_ACCESS_KEY_ID='<AWS_ACCESS_KEY_ID>'
export AWS_SECRET_ACCESS_KEY='<AWS_SECRET_ACCESS_KEY>'
```

**3\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.
``` bash
terraform plan
```

**4\. Execute the Terraform apply.**

Now execute the plan to provision the Federated settings resources.

``` bash
terraform apply
```

**5\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
terraform destroy
```
