# Example - MongoDB Atlas Database User Scope Use

This project aims to provide a very straight-forward example of setting up scope on database user in MongoDB Atlas. Using this, the database user access can be limited to a cluster or data lake.

![MongoDB Atlas DB User](https://github.com/nikhil-mongo/atlas-database-users/blob/master/atlas-1.png?raw=true)

You can view the MongoDB Atlas Cluster Regions from the [documentation](https://docs.atlas.mongodb.com/cloud-providers-regions/).

## Dependencies

* Terraform v0.13
* A MongoDB Atlas account - terraform-providers/mongodbatlas v0.6.5

## Usage

**1\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:


```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently does the below deployments:

- MongoDB cluster - M10
- Creates 2 database users, with the access scope to the database and data lake.

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

- Using the **terraform.tfvars** for storing the password does not shows it in the **terraform plan**. Please refer the **variables.tf** to know more and use **.tfvars** file for storing or passing the details.

```bash
+ resource "mongodbatlas_database_user" "user2" {
      + auth_database_name = "admin"
      + aws_iam_type       = "NONE"
      + id                 = (known after apply)
      + password           = (sensitive value)
...
}
```

**Output**

```bash
Outputs:
atlasclusterstring = [
  {
    "aws_private_link" = {}
    "aws_private_link_srv" = {}
    "private" = ""
    "private_srv" = ""
    "standard" = "mongodb://MongoDB_Atlas-shard-00-00.xgpi2.mongodb.net:27017,MongoDB_Atlas-shard-00-01.xgpi2.mongodb.net:27017,MongoDB_Atlas-shard-00-02.xgpi2.mongodb.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-90b49a-shard-0"
    "standard_srv" = "mongodb+srv://MongoDB_Atlas.xgpi2.mongodb.net"
  },
]
project_name = Atlas-DB-Scope
user1 = dbuser1
user2 = dbuser2
```
