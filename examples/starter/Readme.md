# Example - A basic example to start with the MongoDB Atlas and Terraform

This project aims to provide a very straight-forward example of setting up Terraform with MongoDB Atlas. This will create the following resources in MongoDB Atlas:

- Atlas Project
- MongoDB Cluster - M10
- Database User
- IP Access List

You can refer to the MongoDB Atlas documentation to know about the region names used in MongoDB Atlas respective to the Cloud Provider's region name.
[Amazon Web Services (AWS)](https://docs.atlas.mongodb.com/reference/amazon-aws/#amazon-aws)
[Google Cloud Platform (GCP)](https://docs.atlas.mongodb.com/reference/google-gcp/#google-gcp)
[Microsoft Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/#microsoft-azure)

## Dependencies

* Terraform v0.13 or greater
* A MongoDB Atlas account 
* provider.mongodbatlas: version = "~> 0.9.1"
* Some users (see [here](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1083)) have reported issues deploying this starter example with Mac M1 CPU. If issue persists, suggest deploying instead on x86 linux if possible. See list of supported binaries [here](https://github.com/mongodb/terraform-provider-mongodbatlas/releases/tag/v1.8.1)  

## Usage

**1\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.


> **IMPORTANT** Hard-coding your MongoDB Atlas programmatic API key pair into a Terraform configuration is not recommended. Consider the risks, especially the inadvertent submission of a configuration file containing secrets to a public repository.


**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```

This project currently creates the below deployments:

- Atlas Project
- MongoDB cluster - M10
- Database User
- IP Access list

**3\. Execute the Terraform apply.**

Now execute the plan to provision the MongoDB Atlas resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary charges.

``` bash
$ terraform destroy
```

**Important Point**

You can fetch the connection string as per the use case by following the MongoDB Atlas documentation on [Connect to your cluster](https://docs.atlas.mongodb.com/tutorial/connect-to-your-cluster/index.html).

Or to fetch the connection string using terraform follow the below steps:

```hcl
output "atlasclusterstring" {
    value = mongodbatlas_cluster.cluster.connection_strings
}
```
**Outputs:**
```hcl
atlasclusterstring = [
  {
    "aws_private_link" = {
      "vpce-0ebb76559e8affc96" = "mongodb://pl-0-us-east-1.za3fb.mongodb.net:1024,pl-0-us-east-1.za3fb.mongodb.net:1025,pl-0-us-east-1.za3fb.mongodb.net:1026/?ssl=true&authSource=admin&replicaSet=atlas-d177ke-shard-0"
    }
    "aws_private_link_srv" = {
      "vpce-0ebb76559e8affc96" = "mongodb+srv://mongodb-atlas-pl-0.za3fb.mongodb.net"
    }
    "private" = ""
    "private_srv" = ""
    "standard" = "mongodb://mongodb-atlas-shard-00-00.za3fb.mongodb.net:27017,mongodb-atlas-shard-00-01.za3fb.mongodb.net:27017,mongodb-atlas-shard-00-02.za3fb.mongodb.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-d177ke-shard-0"
    "standard_srv" = "mongodb+srv://mongodb-atlas.za3fb.mongodb.net"
  },
]
```

To fetch a particular connection string, use the **lookup()** function of terraform as below:

```
output "plstring" {
    value = lookup(mongodbatlas_cluster.cluster.connection_strings[0].aws_private_link_srv, aws_vpc_endpoint.ptfe_service.id)
}
```
**Output:**
```
plstring = mongodb+srv://cluster-atlas-pl-0.za3fb.mongodb.net
```
