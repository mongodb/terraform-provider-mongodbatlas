# Example - PrivateLink for Data Federation and Online Archive (Azure)

Set up a private connection to [Data Federation or Online Archive](https://www.mongodb.com/docs/atlas/data-federation/tutorial/config-private-endpoint/) using Azure Private Endpoint and MongoDB Atlas.

## Dependencies

- Terraform `>= 1.0`
- An Azure account with permissions to create networking resources
- A MongoDB Atlas account with a service account credential pair

## Usage

1. Authenticate to Azure:

```bash
az login
```

2. Export MongoDB Atlas credentials (or provide them in `terraform.tfvars`):

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

3. Create `terraform.tfvars` with at least:

```hcl
atlas_client_id     = "<ATLAS_CLIENT_ID>"
atlas_client_secret = "<ATLAS_CLIENT_SECRET>"
project_id          = "<ATLAS_PROJECT_ID>"
azure_location      = "East US 2"
atlas_data_federation_private_link_service_resource_id = "<AZURE_RESOURCE_ID_OF_ATLAS_PRIVATE_LINK_SERVICE>"
```

Optional variables let you customize resource group, VNet, and subnet names/CIDRs.

4. Initialize and review the plan:

```bash
terraform init
terraform plan
```

5. Apply:

```bash
terraform apply
```

6. Destroy when finished:

```bash
terraform destroy
```

## What This Example Creates

- Azure Resource Group, Virtual Network, and Subnet
- Azure `azurerm_private_endpoint` connected to the Atlas-managed Data Federation Private Link Service
- `mongodbatlas_privatelink_endpoint_service_data_federation_online_archive`
- Singular and plural data source reads for verification

## Resource Dependency Chain

1. `mongodbatlas_privatelink_endpoint` requires `project_id` and Atlas Azure region.
2. `azurerm_private_endpoint` requires Atlas private link service metadata and Azure subnet.
3. `mongodbatlas_privatelink_endpoint_service_data_federation_online_archive` requires:
   - Azure private endpoint ID (`endpoint_id`)
