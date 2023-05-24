# Example - Privatelink for Data Federation and Online Archive

Setup private connection to a [Data Federation or Online Archive](https://www.mongodb.com/docs/atlas/data-federation/tutorial/config-private-endpoint/) utilizing [Amazon Virtual Private Cloud (aws vpc)](https://docs.aws.amazon.com/vpc/latest/userguide/what-is-amazon-vpc.html).


## Dependencies

* Terraform v0.13
* An AWS account - provider.aws: version = "~> 4"
* A MongoDB Atlas account - provider.mongodbatlas: version = "~> 1.8"

## Usage

**1\. Ensure your AWS and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

``` bash
$ export AWS_SECRET_ACCESS_KEY='your secret key'
$ export AWS_ACCESS_KEY_ID='your key id'
```

... or the `~/.aws/credentials` file.

```
$ cat ~/.aws/credentials
[default]
aws_access_key_id = your key id
aws_secret_access_key = your secret key

```
... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```
access_key   = "<AWS_ACCESS_KEY_ID>"
secret_key   = "<AWS_SECRET_ACCESS_KEY>"
public_key   = "<MONGODB_ATLAS_PUBLIC_KEY>"
private_key  = "<MONGODB_ATLAS_PRIVATE_KEY>"
project_id   = "<MONGODB_ATLAS_PROJECT_ID>"
```

**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently does the below deployments:

- MongoDB cluster - M10
- AWS Custom VPC, Internet Gateway, Route Tables, Subnets with Public and Private access
- PrivateLink Connection at MongoDB Atlas
- Create VPC Endpoint in AWS

**3\. Configure the security group as required.**

The security group in this configuration allows All Traffic access in Inbound and Outbound Rules.

**4\. Execute the Terraform apply.**

Now execute the plan to provision the AWS and Atlas resources.

``` bash
$ terraform apply
```

**5\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary charges.

``` bash
$ terraform destroy
```

**What's the resource dependency chain?**
1. `mongodbatlas_project` must exist for any of the following
2. `aws_vpc_endpoint` is dependent on the `mongodbatlas_privatelink_endpoint`, and its dependencies.
3. `mongodbatlas_privatelink_endpoint_service` is dependent on `aws_vpc_endpoint` and its dependencies.
4. `mongodbatlas_privatelink_endpoint_service_data_federation_online_archive` is dependent on the `mongodbatlas_project` and `aws_vpc_endpoint`

