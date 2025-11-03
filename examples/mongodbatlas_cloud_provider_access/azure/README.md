# MongoDB Atlas Provider -- Cloud Provider Access Role with AZURE
This example shows how to perform authorization for a cloud provider Azure Service Principal. 

## Dependencies

* Terraform MongoDB Atlas Provider v1.11.0
* A MongoDB Atlas account 
* An AZURE account


```
Terraform v1.5.2
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.11.0
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

**2\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```terraform
client_id     = "<ATLAS_CLIENT_ID>"
client_secret = "<ATLAS_CLIENT_SECRET>"
```

**3\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- An Azure Service Principal
- Confiture Atlas to use your Azure Service Principal

**5\. Execute the Terraform apply.**

Now execute the plan to provision the resources.

``` bash
$ terraform apply
```

**6\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```

