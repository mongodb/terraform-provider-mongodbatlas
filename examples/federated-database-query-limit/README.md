# Example - MongoDB Atlas Federated Database Query Limit with Atlas clusters

This project aims to provide an example of using [MongoDB Atlas Federated Database Query Limit](https://www.mongodb.com/docs/atlas/data-federation/overview/).


## Dependencies

* Terraform MongoDB Atlas Provider v1.10.0
* A MongoDB Atlas account 

```
Terraform v1.4.6
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage
**1\. Create a .tfvars file**

This example requires an Atlas Project to already exist. Once a project is created create the terraform.tfvars file and enter the values for all the required variables, including the project, and make sure **not to commit it**.

**2\. Review the Terraform plan**

Execute the below command and ensure you are happy with the plan. The `terraform plan` command lets you to preview the actions Terraform would take to modify your infrastructure, or save a speculative plan which you can apply later.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- Two MongoDB Atlas clusters in the specified project
- MongoDB Atlas Federated Database Instance based on Atlas clusters
- MongoDB Atlas Federated Database Query Limit

**3\. Execute the Terraform apply.**

Now execute the plan to provision the Federated settings resources. The `terraform apply` command performs a plan just like `terraform plan` does, but then actually carries out the planned changes to each resource using the relevant infrastructure provider's API. It asks for confirmation from the user before making any changes, unless it was explicitly told to skip approval.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges. Calling the `terraform destroy` command will instruct Terraform to terminate / destroy all the resources managed. This will enable you to completely tear down and remove all resources defined in the Terraform State that have previously been deployed.

``` bash
$ terraform destroy
```
