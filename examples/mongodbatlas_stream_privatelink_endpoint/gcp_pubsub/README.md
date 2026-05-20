# MongoDB Atlas Provider - GCP Pub/Sub Privatelink for Atlas Streams

This example shows how to use GCP Private Service Connect for Atlas Streams with GCP Pub/Sub.

A GCP cluster must be provisioned in the same region before creating a GCP Pub/Sub private endpoint.

You must set the following variables:

- `project_id`: Unique 24-hexadecimal digit string that identifies your project
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `cluster_name`: Name of the GCP cluster to provision in the same region
- `gcp_region`: GCP region where your Pub/Sub resources are located
