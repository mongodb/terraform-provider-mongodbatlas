# Atlas Terraform Provider Example: Encryption at rest - AWS - Cluster

This example sets up encryption at rest using AWS KMS for your Atlas Project. It creates the encryption key in AWS KMS, an IAM role and policy so that Atlas can access the key, and enables encryption at rest for the Atlas Project. Finally, it creates a Cluster with encryption at rest enabled.

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
public_key           = "examplepksy"
private_key          = "22b722a9-34f4-3b1b-aada-298329a5c128"
atlas_org_id         = "63f4d4a47baeac59406dc131"
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

## (Optional) Enabling encryption at rest for an existing cluster

1. Import the cluster using the Project ID and cluster name (e.g. `5beae24579358e0ae95492af-MyCluster`):

        $ terraform import mongodbatlas_advanced_cluster.cluster ProjectId-ClusterName

2. Add any non-default values to the cluster resource *mongodbatlas_advanced_cluster.cluster* in *main.tf*. And add the following attribute: `encryption_at_rest_provider = "AWS"`
3. Run terraform apply to enable encryption at rest for the cluster: `terraform apply`
4. (Optional) To remove the cluster from TF state, in case you want to disable project-level encryption and delete the role and key without deleting the imported cluster:
    1. First disable encryption on the cluster by changing the attribute `encryption_at_rest_provider = "NONE"` for the cluster resource *mongodbatlas_advanced_cluster.cluster* in *main.tf*. If you skip this and the next step, you won't be able to disable encryption on the project-level
    2. Run terraform apply to disable encryption for the cluster: `terraform apply`
    3. Finally, remove the cluster from TF state:

            terraform state rm mongodbatlas_advanced_cluster.cluster

    4. You should now be able to run terraform destroy without deleting the cluster: `terraform destroy`
