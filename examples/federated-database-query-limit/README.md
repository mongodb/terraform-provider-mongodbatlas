# Example - MongoDB Atlas Federated Database Query Limit with Atlas clusters

This project aims to provide an example of using [MongoDB Atlas Federated Database Query Limit](https://www.mongodb.com/docs/atlas/data-federation/overview/).


## Dependencies

* Terraform MongoDB Atlas Provider v1.10.0
* A MongoDB Atlas account 

```
Terraform v1.10.0
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage

**1\. Ensure to create an Atlas project**

Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

**2\. Review the Terraform plan**


Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- Two MongoDB Atlas clusters in the specified project
- MongoDB Atlas Federated Database Instance based on Atlas clusters
- MongoDB Atlas Federated Database Query Limit

**3\. Execute the Terraform apply.**

Now execute the plan to provision the Federated settings resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```