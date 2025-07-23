# Example - MongoDB Atlas Federated Database Instance with Microsoft Azure Blob Storage and MongoDB Cluster as storage databases

This project aims to provide an example of using [MongoDB Atlas Federated Database Instance](https://www.mongodb.com/docs/atlas/data-federation/overview/).

## Dependencies

* Terraform MongoDB Atlas Provider v1.39.0
* A MongoDB Atlas account
* An Azure account

```
Terraform v1.39.0
+ provider registry.terraform.io/mongodb/mongodbatlas v1.39.0
```

## Usage

**1\. Set up your Azure credentials.**

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
export MONGODB_ATLAS_PUBLIC_KEY="<ATLAS_PUBLIC_KEY>"
export MONGODB_ATLAS_PRIVATE_KEY="<ATLAS_PRIVATE_KEY>"
```

... or create a **terraform.tfvars** file with all variable values:

```terraform
public_key                 = "<MONGODB_ATLAS_PUBLIC_KEY>"
private_key                = "<MONGODB_ATLAS_PRIVATE_KEY>"
project_id                 = "<ATLAS_PROJECT_ID>"
azure_atlas_app_id         = "<AZURE_ATLAS_APP_ID>"
azure_service_principal_id = "<AZURE_SERVICE_PRINCIPAL_ID>"
azure_tenant_id            = "<AZURE_TENANT_ID>"
```

**3\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the following deployments:

- MongoDB Atlas Cloud Provider Access Setup for Azure
- MongoDB Atlas Cloud Provider Access Authorization
- MongoDB Atlas Federated Database Instance with Azure cloud provider configuration

**5\. Run the Terraform apply command to apply the plan.**

Now run the plan to provision the resources.

``` bash
$ terraform apply
```

**6\. Destroy the resources.**

Once you are finished with our testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```
