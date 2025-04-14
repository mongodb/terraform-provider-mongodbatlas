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
export MONGODB_ATLAS_PUBLIC_KEY="<YOUR_ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<YOUR_ATLAS_PRIVATE_KEY>"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```hcl
public_key           = "examplepksy"
private_key          = "22b722a9-34f4-3b1b-aada-298329a5c128"
atlas_org_id         = "63f4d4a47baeac59406dc131"
```

... or use [AWS Secrets Manager](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/index.md#aws-secrets-manager)


**2\. Set your AWS access key & secret via environment variables:

```bash
export AWS_ACCESS_KEY_ID="exampleawsy"
export AWS_SECRET_ACCESS_KEY="ye3G3+example+TOez2BBK+P22+raXXn2l2hywcuq1"

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
