# Example - Okta and MongoDB Atlas DataLake Pipeline

This project provides an example of using Okta and MongoDB Atlas together.


## Dependencies

* Terraform v0.13
* A MongoDB Atlas account 
You will also need to install the Atlas Terraform provider:
```
Terraform v0.13.0
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage

**1\. Create an Atlas Organization.**

**2\. TFVARS**

Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

**3\. Review the Terraform plan.**

Execute the following command:

``` bash
$ terraform plan
```
Review the output of `terraform plan` to make sure the changes are correct.

This project will deploy the following:

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
