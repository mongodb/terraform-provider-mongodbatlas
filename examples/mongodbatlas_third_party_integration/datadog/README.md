# Example - DataDog third party integration with MongoDB Atlas Cluster

## Dependencies

* Terraform v0.13
* A [DataDog](https://www.datadoghq.com/) API Key.
    * As of March 2023, worth noting that this is *not* an Application Key. API Keys can be created and accessed via Organization Settings from DataDog UI (not from MongoDB Atlas).
* A MongoDB Atlas account - provider.mongodbatlas: version = "~> 1.8"

## Usage

**1\. Ensure your AWS and MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```
public_key      = "<MONGODB_ATLAS_PUBLIC_KEY>"
private_key     = "<MONGODB_ATLAS_PRIVATE_KEY>"
project_id      = "<MONGODB_ATLAS_PROJECT_ID>"
datadog_api_key = "<DATADOG_API_KEY>"
```

**2\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently does the below deployments:

- MongoDB cluster - M10
- Third Party Integration

**4\. Execute the Terraform apply.**

Now execute the plan to provision the AWS and Atlas resources.

``` bash
$ terraform apply
```

**5\. Destroy the resources.**

Once you are finished your testing, ensure you destroy the resources to avoid unnecessary charges.

``` bash
$ terraform destroy
```
