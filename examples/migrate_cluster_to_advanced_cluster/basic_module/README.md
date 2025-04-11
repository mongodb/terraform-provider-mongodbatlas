# Basic module for mongodbatlas_advanced_cluster (Preview for MongoDB Atlas Provider 2.0.0)

The purpose of this example is to demonstrate how a module can be used to manage the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster%2520%2528preview%2520provider%25202.0.0%2529) resource.

## Dependencies
- Terraform CLI >= 1.0
- Terraform MongoDB Atlas Provider `>=v1.29.0`
- A MongoDB Atlas account.
- Configure the provider (can also be done by configuring `public_key` and `private_key` in a `provider.tfvars`).

```bash
export MONGODB_ATLAS_PUBLIC_KEY="xxxx"
export MONGODB_ATLAS_PRIVATE_KEY="xxxx"
```

## Usage

### Update `example.tfvars`

See the example in [example.tfvars](example.tfvars).

### Create the Cluster with `terraform apply`
```bash
terraform init
terraform apply -var-file=example.tfvars
```

### Delete the Cluster with `terraform destroy`
```bash
terraform destroy
```
