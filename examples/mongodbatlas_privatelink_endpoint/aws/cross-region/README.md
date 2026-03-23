# Example - Cluster Private Connection String via AWS Cross-Region PrivateLink

Setup private connection to a [MongoDB Atlas Cluster](https://www.mongodb.com/basics/clusters/mongodb-cluster-setup) utilizing [Amazon Virtual Private Cloud (aws vpc)](https://docs.aws.amazon.com/vpc/latest/userguide/what-is-amazon-vpc.html) with cross-region connectivity.

Unlike the [geosharded example](../cluster-geosharded/) which creates a separate endpoint service per region, this example uses a single endpoint service with `supported_remote_regions` to accept connections from multiple AWS regions.

## Dependencies

* Terraform >= 1.0
* An AWS account - provider.aws: version = "~> 5.0"
* A MongoDB Atlas account - provider.mongodbatlas

## Usage

**1\. Ensure your AWS and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

``` bash
export AWS_SECRET_ACCESS_KEY='<AWS_SECRET_ACCESS_KEY>'
export AWS_ACCESS_KEY_ID='<AWS_ACCESS_KEY_ID>'
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```
access_key      = "<AWS_ACCESS_KEY_ID>"
secret_key      = "<AWS_SECRET_ACCESS_KEY>"
client_id       = "<ATLAS_CLIENT_ID>"
client_secret   = "<ATLAS_CLIENT_SECRET>"
project_id      = "<ATLAS_PROJECT_ID>"
cluster_name    = "aws-cross-region-private-connection"
```

**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently does the below deployments:

- MongoDB cluster - M10
- AWS Custom VPC, Internet Gateway, Route Tables, Subnets in us-east-1 (primary)
- AWS Custom VPC, Internet Gateway, Route Tables, Subnets in us-west-2 (remote)
- PrivateLink Connection at MongoDB Atlas with cross-region support
- VPC Endpoints in both regions connecting to the single Atlas endpoint service

**3\. Execute the Terraform apply.**

Now execute the plan to provision the AWS and Atlas resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary charges.

``` bash
$ terraform destroy
```

**What's the resource dependency chain?**
1. `mongodbatlas_privatelink_endpoint` creates a single endpoint service with `supported_remote_regions` for cross-region connectivity.
2. `aws_vpc_endpoint` in each region (east and west) connects to the same Atlas endpoint service.
3. `mongodbatlas_privatelink_endpoint_service` links each VPC endpoint to Atlas.
4. `mongodbatlas_advanced_cluster` depends on both endpoint services for private connection strings.

**Cross-Region vs Multi-Region**

The [geosharded example](../cluster-geosharded/) creates separate `mongodbatlas_privatelink_endpoint` resources per region. With `supported_remote_regions`, you create one endpoint service and connect from multiple regions, simplifying the setup.
