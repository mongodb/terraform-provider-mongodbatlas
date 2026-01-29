# Module User: GCP Private Service Connect to Port-Mapped Architecture

The purpose of this example is to demonstrate the experience of adopting a new version of a terraform module definition which internally migrated from the legacy GCP Private Service Connect architecture to the port-mapped architecture.
Each module call represents a step on the migration path.
The example focuses on the call of the module rather than the module implementation itself (see the [module maintainer README.md](../module_maintainer/README.md) for the implementation details).

Migration Step Code | Module Version | Config Changes | Plan Changes
--- | --- | --- | ---
[Step 1](./v1) | `v1` | Baseline configuration (legacy architecture) | -
[Step 2](./v2) | `v2` | Upgrade to module version that creates both architectures | Yes (creates new port-mapped resources; legacy resources remain unchanged)
[Step 3](./v3) | `v3` | Upgrade to port-mapped-only module version | Yes (removes legacy resources)

The rest of this example is a step by step guide on how to migrate from legacy to port-mapped architecture:

- [Dependencies](#dependencies)
- [Step 1: Create the legacy architecture with `v1` of the module](#step-1-create-the-legacy-architecture-with-v1-of-the-module)
  - [Update variables](#update-variables)
  - [Run Commands](#run-commands)
- [Step 2: Add port-mapped architecture by using `v2` of the module](#step-2-add-port-mapped-architecture-by-using-v2-of-the-module)
- [Step 3: Migrate to port-mapped-only architecture using `v3` of the module](#step-3-migrate-to-port-mapped-only-architecture-using-v3-of-the-module)
- [Cleanup with `terraform destroy`](#cleanup-with-terraform-destroy)

## Dependencies

- Terraform CLI >= 1.0.
- Terraform MongoDB Atlas Provider with port-mapped architecture support.
- Google Cloud Provider >= 4.0.
- A MongoDB Atlas account.
- A Google Cloud account with appropriate permissions.
- Configure the provider (can also be done by configuring `atlas_client_id` and `atlas_client_secret` in variables).

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

## Step 1: Create the legacy architecture with `v1` of the module

### Update variables

See the example in [v1/variables.tf](v1/variables.tf) or create a `terraform.tfvars` file with your values:

```terraform
project_id     = "<ATLAS_PROJECT_ID>"
gcp_project_id = "<GCP_PROJECT_ID>"
gcp_region     = "us-central1"
```

### Run Commands

```bash
cd v1
terraform init
terraform apply
```

This creates the legacy architecture with Google Compute Addresses and Forwarding Rules. The number matches the `legacy_endpoint_count` variable (defaults to 50), which should match your Atlas project's `privateServiceConnectionsPerRegionGroup` setting. See [Set One Project Limit](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-setgrouplimit) for more information on the default value.

## Step 2: Add port-mapped architecture by using `v2` of the module

The `v2` module automatically creates both legacy and port-mapped architectures simultaneously. No configuration changes are needed - the module creates both sets of resources in parallel.

**Important:** This step creates new resources alongside your existing legacy resources. You'll need to:
1. Update your application connection strings to use the new port-mapped endpoint
2. Test the port-mapped architecture
3. Then proceed to Step 3 to remove legacy resources

```bash
cd v2
cp ../v1/terraform.tfstate . # if you are not using a remote state
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform apply
```

In the plan output, you should see:
- New `mongodbatlas_privatelink_endpoint.port_mapped` with `port_mapping_enabled = true` being created.
- New GCP resources (1 address, 1 forwarding rule) for port-mapped architecture being created.
- New `mongodbatlas_privatelink_endpoint_service.port_mapped` for port-mapped architecture being created.
- Your existing legacy resources remain unchanged (no destruction).

**Note:** After applying, update your application connection strings to use the port-mapped endpoint (connection strings will use `psc-0` identifier instead of `pl-0`). See the [migration guide](../../../../docs/guides/gcp-privatelink-port-mapping-migration.md) for details.

## Step 3: Migrate to port-mapped-only architecture using `v3` of the module

The `v3` module only supports port-mapped architecture. This step removes all legacy resources.

**Important:** Before proceeding, ensure:
1. Your applications are using the port-mapped endpoint connection strings
2. You've tested the port-mapped architecture and confirmed it works correctly
3. You have a backup of your Terraform state

```bash
cd v3
cp ../v2/terraform.tfstate . # if you are not using a remote state
terraform init -upgrade # in case your Atlas Provider version needs to be upgraded
terraform plan # Review the plan carefully
terraform apply
```

In the plan output, you should see:
- Legacy endpoint resources planned for destruction.
- Legacy GCP resources (addresses and forwarding rules matching the `legacy_endpoint_count` variable) planned for destruction.
- Only port-mapped architecture resources remain.

## Cleanup with `terraform destroy`

```bash
terraform destroy
```
