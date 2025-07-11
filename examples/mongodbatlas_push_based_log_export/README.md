# MongoDB Atlas Provider -- Push-Based Log Export 
This example shows how to configure push-based log export for an Atlas project.

## Dependencies

* Terraform MongoDB Atlas Provider v1.16.0 minimum
* Terraform AWS provider
* A MongoDB Atlas account 
* An AWS account


```
Terraform v1.5.2
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.16.0
```

## Usage

**1\. Ensure your AWS and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="<ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<ATLAS_PRIVATE_KEY>"
```

``` bash
export AWS_ACCESS_KEY_ID='<AWS_ACCESS_KEY_ID>'
export AWS_SECRET_ACCESS_KEY='<AWS_SECRET_ACCESS_KEY>'
```

... or the `~/.aws/credentials` file.

```
$ cat ~/.aws/credentials
[default]
aws_access_key_id = <AWS_ACCESS_KEY_ID>
aws_secret_access_key = <AWS_SECRET_ACCESS_KEY>
```
... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values. For example:
```
access_key           = "<AWS_ACCESS_KEY_ID>"
secret_key           = "<AWS_SECRET_ACCESS_KEY>"
public_key           = "<ATLAS_PUBLIC_KEY>"
private_key          = "<ATLAS_PRIVATE_KEY>"
```

**2\. Review the Terraform plan.**

Execute the following command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the following deployments:

- An AWS IAM Policy
- An AWS IAM Role
- An AWS S3 bucket
- An IAM role policy for the S3 bucket
- Configure Atlas to use your AWS Role
- An Atlas project in the configured Atlas organization
- Configure push-based log export to the S3 bucket for Atlas project

**3\. Execute the Terraform apply.**

Now execute the plan to provision the resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

When you have finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```

