# MongoDB Atlas Provider - AWS Confluent Privatelink for Atlas Streams 

This example shows how to use AWS Confluent Privatelink for Atlas Streams with a dedicated Confluent Cluster. 

You must set the following variables:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_id`: Unique 24-hexadecimal digit string that identifies your project
- `confluent_cloud_api_key`: Public API key to authenticate to Confluent Cloud
- `confluent_cloud_api_secret`: Private API key to authenticate to Confleunt Cloud
- `subnets_to_privatelink`: A map of Zone ID to Subnet ID (i.e.: {\"use1-az1\" = \"subnet-abcdef0123456789a\", ...})
- `aws_account_id`: The AWS Account ID (12 digits)
- `aws_region`: The AWS Region
