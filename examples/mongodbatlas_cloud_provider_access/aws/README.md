# MongoDB Atlas Provider -- Cloud Provider Access Role with AWS
This example shows how to perform authorization for a cloud provider AWS role.

## Dependencies

* Terraform MongoDB Atlas Provider v1.10.0
* A MongoDB Atlas account 
* An AWS account


```
Terraform v1.5.2
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage

**1\. Ensure your AWS and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="<YOUR_ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<YOUR_ATLAS_PRIVATE_KEY>"
```

``` bash
export AWS_ACCESS_KEY_ID='<YOUR_AWS_KEY_ID>'
export AWS_SECRET_ACCESS_KEY='<YOUR_AWS_SECRET_ACCESS_KEY>'
```

... or the `~/.aws/credentials` file.

```
$ cat ~/.aws/credentials
[default]
aws_access_key_id = <YOUR_AWS_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_AWS_SECRET_ACCESS_KEY>
```
... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```
access_key   = "<YOUR_AWS_ACCESS_KEY_ID>"
secret_key   = "<YOUR_AWS_SECRET_ACCESS_KEY>"
public_key   = "<YOUR_ATLAS_PUBLIC_KEY>"
private_key  = "<YOUR_ATLAS_PRIVATE_KEY>"
```

**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- An AWS Policy
- An AWS Role
- Confiture Atlas to use your AWS Role

Please note: the policy is intentionally restricted to a _Deny All_. You need to update it accordingly based on the permissions you will need for this role.

**3\. Execute the Terraform apply.**

Now execute the plan to provision the resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```

