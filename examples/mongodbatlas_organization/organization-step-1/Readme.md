# Example - A basic example to create a MongoDB Atlas Organization with Terraform

This project aims to provide a very straight-forward example of setting up a MongoDB Atlas Organization with Terraform. This will create the following resources in MongoDB Atlas:

- Atlas organization
- Private Key
- Public Key
- Organization ID

## Dependencies

* Terraform v0.15 or greater
* A MongoDB Atlas account 
* provider.mongodbatlas: version = "~> 1.10"
* [Cross-organization billing](https://www.mongodb.com/docs/atlas/billing/#cross-organization-billing) enabled and the requesting API Key's organization must be a paying organization. 
* Some users (see [here](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1083)) have reported issues deploying this starter example with Mac M1 CPU. you encounter this issue, try deploying instead on x86 linux if possible. See list of supported binaries [here](https://github.com/mongodb/terraform-provider-mongodbatlas/releases/tag/v1.8.1)  

## Usage

**1\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

... or utilize the `variables.tf` file and create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.


> **IMPORTANT** Hard-coding your MongoDB Atlas Service Account credentials into a Terraform configuration is not recommended. Consider the risks, especially the inadvertent submission of a configuration file containing secrets to a public repository.


**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```

This example currently creates the following:

- Atlas organization
- Private Key
- Public Key
- Organization ID

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

**Output:**
