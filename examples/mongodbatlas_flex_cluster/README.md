# MongoDB Atlas Provider -- Atlas Flex Cluster
This example creates one flex cluster in a project.

**NOTE**: Flex Clusters are currently in Preview. To use this feature, you must take the following actions:
1. Enable the `Atlas USS` Preview Feature in your organization (contact [MongoDB Support](https://www.mongodb.com/services/support)).
2. Enable the [Preview Features](../../README.md#preview-features) when running `terraform` commands.


Variables Required to be set:
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `project_id`: Project ID where flex cluster will be created