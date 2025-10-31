# MongoDB Atlas Provider - GCP Confluent Privatelink for Atlas Streams

This example shows how to use GCP Private Service Connect for Atlas Streams with Confluent Cloud.

You must set the following variables:

- `project_id`: Unique 24-hexadecimal digit string that identifies your project
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `gcp_region`: GCP region where your Confluent cluster is located
- `confluent_dns_domain`: DNS domain for your Confluent cluster
- `confluent_dns_subdomains`: List of DNS subdomains for your Confluent cluster
- `service_attachment_uris`: List of GCP service attachment URIs from your Confluent Cloud cluster
