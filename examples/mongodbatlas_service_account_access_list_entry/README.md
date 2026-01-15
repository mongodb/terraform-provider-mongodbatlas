# MongoDB Atlas Provider -- Service Account Access List Entry

This example shows how to create Service Account Access List entries in MongoDB Atlas.

## Example Overview

The included [`main.tf`](./main.tf) shows how to:

1. **Create a Service Account** at the organization level with `mongodbatlas_service_account`.
2. **Add IP Access List entries** to the Service Account using `mongodbatlas_service_account_access_list_entry` with either CIDR blocks or IP addresses.
3. **Read access list entries** using data sources.

## Variables Required to be set:

- `org_id`: MongoDB Atlas Organization ID
- `atlas_client_id`: MongoDB Atlas Service Account Client ID (for provider authentication)
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret (for provider authentication)

## Important Notes

~> **IMPORTANT WARNING:** Managing Service Accounts with Terraform **exposes sensitive organizational secrets** in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

When you remove an entry from the access list, existing connections from the removed address(es) may remain open for a variable amount of time.

## Outputs

- `access_list_entry_cidr_block`: The CIDR block from the access list entry
- `all_access_list_entries`: All access list entries for the Service Account
