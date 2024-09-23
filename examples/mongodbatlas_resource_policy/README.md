# MongoDB Atlas Provider -- Atlas Resource Policy
This example creates three different resource policies in an organization.

**NOTE**: Resource Policies are currently in Public Preview. You will need to
1. Enable the `Atlas Resource Policies` Beta Feature in your organization.
2. Enable the [Preview Features](../../README.md#preview-features) when running `terraform` commands.


Variables Required to be set:
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `org_id`: Organization ID where project will be created
