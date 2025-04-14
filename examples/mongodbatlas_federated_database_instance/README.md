# Example - MongoDB Atlas Federated Database Instance with an Amazon S3 and MongoDB Cluster as storage databases

This project aims to provide an example of using [MongoDB Atlas Federated Database Instance](https://www.mongodb.com/docs/atlas/data-federation/overview/).


## Dependencies

* Terraform MongoDB Atlas Provider v1.10.0
* A MongoDB Atlas account 
* An AWS account

```
Terraform v1.10.0
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage

**1\. Ensure your AWS and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="<ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<ATLAS_PRIVATE_KEY>"
```

```bash
export AWS_SECRET_ACCESS_KEY="<AWS_SECRET_ACCESS_KEY>"
export AWS_ACCESS_KEY_ID="<AWS_ACCESS_KEY_ID>"
```

... or the `~/.aws/credentials` file.

```
$ cat ~/.aws/credentials
[default]
aws_access_key_id = <AWS_ACCESS_KEY_ID>
aws_secret_access_key = <AWS_SECRET_ACCESS_KEY>
```
... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```
access_key   = "<AWS_ACCESS_KEY_ID>"
secret_key   = "<AWS_SECRET_ACCESS_KEY>"
public_key   = "<ATLAS_PUBLIC_KEY>"
private_key  = "<ATLAS_PRIVATE_KEY>"
```
**2\. Create an S3 bucket into your AWS account**
Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

**3\. Review the Terraform plan. **

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- An AWS Policy
- An AWS Role
- MongoDB Atlas Federated Database Instance with an S3 bucket and MongoDB Atlas Cluster as storage databases

**5\. Execute the Terraform apply.**

Now execute the plan to provision the Federated settings resources.

``` bash
$ terraform apply
```

**6\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```
