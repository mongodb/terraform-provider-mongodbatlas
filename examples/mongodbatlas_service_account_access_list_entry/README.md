# MongoDB Atlas Provider -- Service Account Access List Entry

This example shows how to create Service Account Access List entries in MongoDB Atlas.

## Variables Required to be set:

- `org_id`: MongoDB Atlas Organization ID
- `atlas_client_id`: MongoDB Atlas Service Account Client ID (for provider authentication)
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret (for provider authentication)

## Prerequisites
- Service Account with Organization Owner permissions

## Outputs

- `access_list_entry_cidr_block`: The CIDR block from the access list entry
- `all_access_list_entries`: All access list entries for the Service Account
