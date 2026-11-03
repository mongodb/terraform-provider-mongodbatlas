# MongoDB Atlas Provider -- Stream Workspace Private Endpoint

This example shows how to register an existing AWS VPC interface endpoint for inbound private connectivity to Atlas Stream Processing in a project.

## Required Variables
- `atlas_client_id`: The MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: The MongoDB Atlas Service Account Client Secret
- `project_id`: The Atlas project ID
- `vpc_endpoint_id`: The AWS VPC interface endpoint ID
- `vpc_endpoint_dns_name`: The DNS name of the AWS VPC interface endpoint

## Outputs
- `endpoint_status`: The status of the endpoint registration, retrieved via the singular data source
- `endpoint_ids`: All workspace private endpoints registered in the project, retrieved via the plural data source
