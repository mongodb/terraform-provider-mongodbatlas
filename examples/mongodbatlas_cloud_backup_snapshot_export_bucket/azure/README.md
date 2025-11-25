# MongoDB Atlas Provider - Atlas Cloud Backup Snapshot Export Bucket in Azure

This example shows how to set up Cloud Backup Snapshot Export Bucket in Atlas through Terraform.

You must set the following variables:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID.
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret.
- `project_id`: Unique 24-hexadecimal digit string that identifies the project where the stream instance will be created.
- `azure_tenant_id`: The Tenant ID which should be used.
- `subscription_id`: Azure Subscription ID.
- `client_id`: Azure Client ID.
- `client_secret`: Azure Client Secret.
- `tenant_id`: Azure Tenant ID.
- `azure_atlas_app_id`: The client ID of the application for which to create a service principal.
- `azure_resource_group_location`: The Azure Region where the Resource Group should exist.
- `storage_account_name`: Specifies the name of the storage account.

To learn more, see the [Export Cloud Backup Snapshot Documentation](https://www.mongodb.com/docs/atlas/backup/cloud-backup/export/).


