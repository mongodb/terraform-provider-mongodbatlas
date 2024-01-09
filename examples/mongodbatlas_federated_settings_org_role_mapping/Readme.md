# Example - Okta and MongoDB Atlas Federated Settings Configuration

This project aims to provide an example of using Okta and MongoDB Atlas together.


## Dependencies

* Terraform v0.13
* Okta account 
* A MongoDB Atlas account 

```
Terraform v0.13.0
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.4.0
```

## Usage

**1\. Ensure your Okta/Mongodb Atlas Federal settings configuration is set up to have a working set of organizations, verified domains, and identity providers.**

**2\. TFVARS**

Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

**3\. Review the Terraform plan. **

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently does the below deployments:

- MongoDB Atlas Federated Settings Organizational Role Mapping
- MongoDB Atlas Federated Settings Organizational Identity Provider
- MongoDB Atlas Federated Settings Organizational configuration

**4\. Execute the Terraform import for 2 resources that do not support create.**
``` bash
$ terraform import mongodbatlas_federated_settings_identity_provider.identity_provider  6287a67f7f7f7f7f441c6c-0oad7f7f7f7fk1297
  terraform import mongodbatlas_federated_settings_org_config.org_connections_import  6287a67f7f7f7f7f441c6c-627a96837f7f7f7f7e306f14

```

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
