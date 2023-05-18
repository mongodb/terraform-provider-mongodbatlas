# Example - Okta and MongoDB Atlas DataLake Pipeline

This project aims to provide an example of using Okta and MongoDB Atlas together.


## Dependencies

* Terraform v0.13
* Okta account 
* A MongoDB Atlas account 

```
Terraform v0.13.0
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage

**1\. Create an Atlas Organization.**

**2\. TFVARS**

Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

**3\. Review the Terraform plan. **

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently does the below deployments:

- MongoDB Atlas Project
- MongoDB Atlas Cluster
- MongoDB Atlas DataLake Pipeline

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
