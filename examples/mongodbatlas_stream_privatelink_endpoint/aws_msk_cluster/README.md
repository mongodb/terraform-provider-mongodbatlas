# MongoDB Atlas Provider - AWS MSK Privatelink for Atlas Streams

This example shows how to use AWS Privatelink for Atlas Streams with an AWS MSK cluster.

You must set the following variables:

- `public_key`: Public API key to authenticate to Atlas
- `private_key`: Private API key to authenticate to Atlas
- `project_id`: Unique 24-hexadecimal digit string that identifies your project
- `aws_account_id`: The AWS Account ID (12 digits)
- `msk_cluster_name`: The MSK cluster's desired name
- `aws_secret_arn`: AWS Secrets Manager secret ARN. Must meet the criteria outlined in https://docs.aws.amazon.com/msk/latest/developerguide/msk-password-tutorial.html
