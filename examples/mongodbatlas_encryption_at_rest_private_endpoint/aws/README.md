# MongoDB Atlas Provider - Encryption At Rest using Customer Key Management via Private Network Interfaces (AWS)
This example shows how to configure encryption at rest using AWS with customer managed keys ensuring all communication with AWS Key Management Service (KMS) happens exclusively over AWS PrivateLink.

## Dependencies

* Terraform MongoDB Atlas Provider v1.28.0 minimum
* A MongoDB Atlas account 
* Terraform AWS provider
* An AWS account

## Usage

**1\. Provide the appropriate values for the input variables.**

- `atlas_public_key`: The public API key for MongoDB Atlas
- `atlas_private_key`: The private API key for MongoDB Atlas
- `atlas_project_id`: Atlas Project ID
- `aws_kms_key_id`: ARN that identifies the Amazon Web Services (AWS) Customer Master Key (CMK) to use to encrypt and decrypt
- `atlas_aws_region`: Region in which the Encryption At Rest private endpoint is located

**2\. Review the Terraform plan.**

Execute the following command and ensure you are happy with the plan.

``` bash
$ terraform plan
```
This project will execute the following changes to acheive successful encryption at rest over AWS PrivateLink for customer managed keys:

- Configure encryption at rest in an existing project using a custom AWS KMS Key. For successful private networking configuration, the `requires_private_networking` attribute in `mongodbatlas_encryption_at_rest.aws_kms_config` is set to `true`.
- Create a private endpoint for the existing project under a certain AWS region using `mongodbatlas_encryption_at_rest_private_endpoint`. 

**3\. Execute the Terraform apply.**

Now execute the plan to provision the resources.

``` bash
$ terraform apply
```

**4\. Destroy the resources.**

When you have finished your testing, ensure you destroy the resources to avoid unnecessary |service| charges.

``` bash
$ terraform destroy
```

