# MongoDB Atlas Provider - AWS Confluent Privatelink for Atlas Streams 

This example shows how to use AWS Confluent Privatelink for Atlas Streams for Serveless products

You must set the following variables:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_id`: Unique 24-hexadecimal digit string that identifies your project
- `confluent_cloud_api_key`: Public API key to authenticate to Confluent Cloud
- `confluent_cloud_api_secret`: Private API key to authenticate to Confleunt Cloud
- `subnets_to_privatelink`: A map of Zone ID to Subnet ID (i.e.: {\"use1-az1\" = \"subnet-abcdef0123456789a\", ...})
- `vpc_id`: The ID of the VPC in which the endpoint will be used.
- `aws_region`: The AWS Region
