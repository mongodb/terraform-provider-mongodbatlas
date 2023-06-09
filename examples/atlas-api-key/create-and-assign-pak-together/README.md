# MongoDB Atlas Provider -- Create and Assign PAK together
This example creates a project API Key and access list, showing how to attach a CIDR block.

Variables Required to be set:
- `project_id`: ID of the Atlas project
- `org_id`: ID of Atlas organization
- `public_key`: Atlas public key
- `private_key`: Atlas  private key

In this example, we will set up a project API key and attach an access list to it.


**Note:** in this example parameter role_names is deprecated and will be removed in v1.12.0 release from codebase. Use `project_assignment`  parameter instead. 