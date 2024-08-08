# MongoDB Atlas Provider - Atlas Cloud Backup Snapshot Export Bucket in AWS

This example shows how to set up Cloud Backup Snapshot Export Bucket in Atlas through Terraform.

You must set the following variables:

- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `project_id`: Unique 24-hexadecimal digit string that identifies the project where the stream instance will be created.
- `access_key`: AWS Access Key
- `secret_key`: AWS Secret Key.
- `aws_region`: AWS region.

To learn more, see the [Export Cloud Backup Snapshot Documentation](https://www.mongodb.com/docs/atlas/backup/cloud-backup/export/).


