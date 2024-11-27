# MongoDB Atlas Provider -- Atlas Resource Policy
This example creates three different resource policies in an organization.

**NOTE**: Resource Policies are currently in Private Preview. To use this feature, you must take the following actions:
1. Enable the `Atlas Resource Policies` Preview Feature in your organization (contact [MongoDB Support](https://www.mongodb.com/services/support)).
2. Enable the [Preview Features](../../README.md#preview-features) when running `terraform` commands.


Variables Required to be set:
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `org_id`: Organization ID where project will be created
