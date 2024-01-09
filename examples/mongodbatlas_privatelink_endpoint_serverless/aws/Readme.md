# Example - AWS and Atlas PrivateLink with Terraform

This project aims to provide a very straight-forward example of setting up PrivateLink connection between AWS and MongoDB Atlas Serverless.


## Dependencies

* Terraform v0.13
* An AWS account - provider.aws: version = "~> 3.3"
* A MongoDB Atlas account - provider.mongodbatlas: version = "~> 0.6"

## Usage

**1\. Ensure your AWS and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

``` bash
$ export AWS_SECRET_ACCESS_KEY='your secret key'
$ export AWS_ACCESS_KEY_ID='your key id'
```

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

... or the `~/.aws/credentials` file.

```
$ cat ~/.aws/credentials
[default]
aws_access_key_id = your key id
aws_secret_access_key = your secret key

```
... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

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

**Important Point**

To fetch the connection string follow the below steps:
```
output "atlasclusterstring" {
    value = data.mongodbatlas_serverless_instance.cluster_atlas.connection_strings_standard_srv
}
```
**Outputs:**
```
atlasclusterstring =  "mongodb+srv://cluster-atlas.za3fb.mongodb.net"
 
```

To fetch a private connection string, use the output of terraform as below after second apply:

```
output "plstring" {
 value = mongodbatlas_serverless_instance.cluster_atlas.connection_strings_private_endpoint_srv[0]
}
```
**Output:**
```
plstring = mongodb+srv://cluster-atlas-pe-0.za3fb.mongodb.net
```
