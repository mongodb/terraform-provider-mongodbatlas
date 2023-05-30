# Example - MongoDB Atlas Clsuetr Outage Simulation on a multi-region cluster

This project aims to provide an example of using [MongoDB Atlas Cluster Outage Simulation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cluster-Outage-Simulation).


## Dependencies

* Terraform MongoDB Atlas Provider v1.11.0
* A MongoDB Atlas account 

```
Terraform v1.10.0
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage

**1\. Ensure to create an Atlas project**

2\. Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

**3\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- One MongoDB Atlas multi-region (US_EAST_1, US_EAST_2, US_WEST_1) cluster in the specified project.
- Cluster Outage Simulation on the created cluster.

**4\. Execute the Terraform apply.**

Now execute the plan to provision the Atlas Cluster and start outage simulation on this cluster.

``` bash
$ terraform apply
```

**6\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```
