# MongoDB Atlas Provider -- Cluster with pinned FCV

Example shows how to pin the FCV of a cluster making use of `pinned_fcv` block. This enables direct control to pin cluster’s FCV before performing an upgrade on the `mongo_db_major_version`. Users can then downgrade to the previous MongoDB version with minimal risk if desired, as the FCV is maintained.

The unpin operation can be performed by removing the `pinned_fcv` block. **Note**: Once FCV is unpinned it will not be possible to downgrade the `mongo_db_major_version`. If FCV is unpinned past the expiration date the `pinned_fcv` attribute must be removed.

The following [knowledge hub article](https://kb.corp.mongodb.com/article/000021785/) and [FCV documentation](https://www.mongodb.com/docs/atlas/tutorial/major-version-change/#manage-feature-compatibility--fcv--during-upgrades) can be referenced for more details.

## Dependencies

* Terraform MongoDB Atlas Provider v2.0.0 or later
* A MongoDB Atlas account 

```
Terraform >= 0.13
+ provider registry.terraform.io/terraform-providers/mongodbatlas v2.0.0
```


## Usage
**1\. Ensure your MongoDB Atlas credentials are set up.**

Following `variables.tf` file create **terraform.tfvars** file with all the variable values, as demonstrated below:
```
public_key           = "<MONGODB_ATLAS_PUBLIC_KEY>"
private_key          = "<MONGODB_ATLAS_PRIVATE_KEY>"
atlas_project_id     = "<MONGODB_ATLAS_PROJECT_ID>"
fcv_expiration_date  = "<FCV pin expiration date, e.g. 2024-11-22T10:50:00Z>"
```

**2\. Review the Terraform plan.**

Execute the following command.

``` bash
$ terraform plan
```
This project currently supports the following deployments:

- A Cluster with pinned FCV configured.

**3\. Execute the Terraform apply.**

Execute the following plan to provision the Atlas Project and Cluster resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

Once you finished your testing, ensure you destroy the resources to avoid unnecessary Atlas charges.

``` bash
$ terraform destroy
```
