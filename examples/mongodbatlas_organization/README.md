
# Example - MongoDB Atlas Organization with Terraform

This project provides examples for both creating new MongoDB Atlas Organizations and importing existing ones using Terraform.

- MongoDB Atlas organization
- Private Key
- Public Key
- Organization ID
- MongoDB Atlas Project

## Dependencies

* Terraform v0.15 or greater
* A MongoDB Atlas account 
* provider.mongodbatlas: version = "~> 1.10" for creating an organization, 1.38 for importing an existing organization.
* [Cross-organization billing](https://www.mongodb.com/docs/atlas/billing/#cross-organization-billing) enabled and the requesting API Key's organization must be a paying organization. 
* Some users (see [here](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1083)) have reported issues deploying this starter example with Mac M1 CPU. you encounter this issue, try deploying instead on x86 linux if possible. See list of supported binaries [here](https://github.com/mongodb/terraform-provider-mongodbatlas/releases/tag/v1.8.1)  

## Usage - New organization
**1\. change working directry to folder organization-step-1.**

**2\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

... or utilize the `variables.tf` file and create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.


> **IMPORTANT** Hard-coding your MongoDB Atlas Service Account credentials into a Terraform configuration is not recommended. Consider the risks, especially the inadvertent submission of a configuration file containing secrets to a public repository.


**3\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```

This example currently creates the following:

- Atlas organization
- Private Key
- Public Key
- Organization ID

**4\. Execute the Terraform apply.**

Now execute the plan to provision the MongoDB Atlas resources.

``` bash
$ terraform apply
```

**Output:**

mongodbatlas_organization.test: Creating...
mongodbatlas_organization.test: Creation complete after 1s [id=b3fff2lk:NjffffyMmE2M2fffffffOTkwM2I0]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

org_id = "<ORG_ID>"
org_private_key = "<ORG_PRIVATE_KEY>"
org_public_key = "<ORG_PUBLIC_KEY>"

**5\. Retain values for org_private_key and org_public_key for next stage of example as new API key has access to create resources in new organization.**

**6\. change working directry to folder organization-step-2.**

**7\. Ensure your MongoDB Atlas credentials are set up to use new public and private key.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="<ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<ATLAS_PRIVATE_KEY>"
```

... or utilize the `variables.tf` file and create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.


> **IMPORTANT** Hard-coding your MongoDB Atlas programmatic API key pair into a Terraform configuration is not recommended. Consider the risks, especially the inadvertent submission of a configuration file containing secrets to a public repository.


**8\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```

This example currently creates the following:

- Atlas Project

**9\. Execute the Terraform apply.**

Now execute the plan to provision the MongoDB Atlas resources.

``` bash
$ terraform apply
```
mongodbatlas_project.project: Creating...
mongodbatlas_project.project: Creation complete after 4s [id=647fe6baffffffdcaee72]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

project_name = "testnew"

**9\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary charges.

``` bash
$ terraform destroy
```


**Output:**
  - project_name = "testnew" -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

mongodbatlas_project.project: Destroying... [id=647fe6ba6fc6fc0efdcaee72]
mongodbatlas_project.project: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.

cd ../organization-step-1

``` bash
$ terraform destroy
```


mongodbatlas_organization.test: Destroying... [id=b3JnX2lk:NjQ3ZfffffNWU2NzNlOTkwM2I0]
mongodbatlas_organization.test: Destruction complete after 9s

Destroy complete! Resources: 1 destroyed.


## Usage - Importing an existing organization

This is useful when you already have the MongoDB Atlas organization created but you want to start managing its settings with Infrastructure as Code.

See the example in directory [./organization-import](./organization-import).
