# Example - MongoDB Atlas Federated Database Instance

This project aims to provide an example of using [MongoDB Atlas Federated Database Instance](https://www.mongodb.com/docs/atlas/data-federation/overview/).


## Dependencies

* Terraform v1.10.0
* A MongoDB Atlas account 
* An AWS account

```
Terraform v1.10.0
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage

**1\. Ensure to create an Atlas project and a cluster**
**2\. Create an s3 bucket into your AWS account**
Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

**3\. Review the Terraform plan. **

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- An AWS Policy
- An AWS Role
- MongoDB Atlas Federated Database Instance

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
