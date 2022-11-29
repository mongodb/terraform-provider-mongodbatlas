# Example - A basic example configuring MongoDB Atlas Third Party Integrations and Terraform

This project aims to provide a very straight-forward example of setting up Terraform with MongoDB Atlas. This will create the following resources in MongoDB Atlas:

- Atlas Project
- Microst Teams Third Party Integration
- Prometheus Third Party Integration


You can refer to the MongoDB Atlas documentation to know about the parameters that support Third Party Integrations.

[Prometheus](https://www.mongodb.com/docs/atlas/tutorial/prometheus-integration/#std-label-httpsd-prometheus-config)

[Microsoft Teams](https://www.mongodb.com/docs/atlas/tutorial/integrate-msft-teams/)

## Dependencies

* Terraform v0.13 or greater
* A MongoDB Atlas account 
* provider.mongodbatlas: version = "~> 0.9.1"

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
- Microst Teams Third Party Integration
- Prometheus Third Party Integration

**3\. Execute the Terraform apply.**

Now execute the plan to provision the MongoDB Atlas resources. 
Note: you can find the Prometheus Config details under Outputs section. 

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary charges.

``` bash
$ terraform destroy
```

