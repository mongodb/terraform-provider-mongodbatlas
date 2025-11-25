# Example - MongoDB Atlas Cloud Backup Schedule for setting up policies for multiple clusters

This project aims to provide an example of using [Cloud Backup Schedule in Atlas](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/modify-one-schedule/).


## Dependencies

* Terraform MongoDB Atlas Provider v1.10.0
* A MongoDB Atlas account 

```
Terraform v1.10.0
+ provider registry.terraform.io/terraform-providers/mongodbatlas v1.10.0
```

## Usage

**1\. Ensure your MongoDB Atlas credentials are set up.**

This can be done using environment variables:

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

... or follow as in the `variables.tf` file and create **terraform.tfvars** file with all the variable values, ex:
```
client_id     = "<ATLAS_CLIENT_ID>"
client_secret = "<ATLAS_CLIENT_SECRET>"
```
**2\. Update required variables.**
Now create/update **terraform.tfvars** file with all the variable values and make sure **not to commit it**. For this example, you just need to provide `org_id` and a `project_name`.

**3\. Review the Terraform plan.**

Execute the below command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project currently supports the below deployments:

- MongoDB Atlas Project
- MongoDB Atlas Clusters (2 AWS M10 clusters in different regions) 
- MongoDB Cloud Backup Schedule(s) with various policies which is set up for each created cluster.

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
