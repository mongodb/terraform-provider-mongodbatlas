# MongoDB Atlas Federated Database Instance Examples

Query data across MongoDB Atlas clusters, AWS S3 buckets, Azure Blob Storage containers, and Google Cloud Storage buckets from a single federated database instance.

Each sibling directory is its own root module and depends only on its own variables. See each README for the credentials and values it needs.

For product documentation, see [Atlas Data Federation](https://www.mongodb.com/docs/atlas/data-federation/adf-overview/overview/).

## Sibling examples

- [`aws/`](aws/README.md) - maps an S3 bucket and an Atlas cluster as data stores, and creates the IAM role and policy that Atlas assumes.
- [`azure/`](azure/README.md) - configures the Azure cloud provider role for the instance.
- [`gcp/`](gcp/README.md) - configures the Google Cloud provider role for the instance.
- [`private-endpoint/`](private-endpoint/README.md) - an AWS VPC and private endpoint for the instance.
