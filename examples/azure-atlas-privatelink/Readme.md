# Example - Microsoft Azure and MongoDB Atlas Private Endpoint

This project aims to provide an example of using Azure and MongoDB Atlas together.


## Dependencies

* Terraform v0.13
* Microsoft Azure account
* MongoDB Atlas account

```
Terraform v0.13.0
+ provider registry.terraform.io/hashicorp/azuread v1.0.0
+ provider registry.terraform.io/hashicorp/azurerm v2.31.1
+ provider registry.terraform.io/terraform-providers/mongodbatlas v0.6.5
```

## Usage

**1\. Ensure your Azure credentials are set up.**

1. Install the Azure CLI by following the steps from the [official Azure documentation](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli).
2. Run the command `az login` and this will take you to the default browser and perform the authentication.
3. Once authenticated, it will print the user details as below:

```
⇒  az login
You have logged in. Now let us find all the subscriptions to which you have access...
The following tenants don't contain accessible subscriptions. Use 'az login --allow-no-subscriptions' to have tenant level access.
XXXXX
[
  {
    "cloudName": "AzureCloud",
    "homeTenantId": "XXXXX",
    "id": "XXXXX",
    "isDefault": true,
    "managedByTenants": [],
    "name": "Pay-As-You-Go",
    "state": "Enabled",
    "tenantId": "XXXXX",
    "user": {
      "name": "person@domain.com",
      "type": "user"
    }
  }
]
```

**2\. TFVARS**

Now create **terraform.tfvars** file with all the variable values and make sure **not to commit it**.

An existing cluster on the project can optionally be linked via the `cluster_name` variable.
If included, the azure connection string to the cluster will be output.

**3\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently does the below deployments:

- MongoDB Atlas Azure Private Endpoint
- Azure Resource Group, VNET, Subnet, Private Endpoint
- Azure-MongoDB Private Link

**4\. Execute the Terraform apply.**

Now execute the plan to provision the Azure resources.

``` bash
$ terraform apply
```

**5\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Azure and Atlas charges.

``` bash
$ terraform destroy
```
