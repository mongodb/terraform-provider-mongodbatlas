# MongoDB Atlas Provider -- Atlas IP Access List
This example creates a project API access list showing how to attach multiple IP addresses and CIDR Blocks.

Variables Required to be set:
- `project_id`: ID of the Atlas project
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `ip_address`: IP addresses you want to permit access to
- `cidr_block`: CIDR block you want to permit access to
- `comment`: If provider_name is tenant, the backing provider (AWS, GCP)


For this example, we will setup two access ranges to show multiple IP support and multiple CIDR block.


