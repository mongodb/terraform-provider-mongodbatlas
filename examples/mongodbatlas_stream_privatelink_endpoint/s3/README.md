# MongoDB Atlas Provider - AWS S3 Privatelink for Atlas Streams

This example shows how to use AWS Privatelink for Atlas Streams with an AWS S3 bucket.

You must set the following variables:

- `project_id`: Unique 24-hexadecimal digit string that identifies your project
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `region`: Region where the S3 bucket is located
- `service_endpoint_id`: Service endpoint ID (should follow the format `com.amazonaws.<region>.s3`)
- `s3_bucket_name`: Name of the S3 bucket for stream data
