# MongoDB Atlas Provider — Advanced Cluster (TLS 1.2 Configuration)

This example creates a project and a replica set cluster with TLS 1.2 cipher configuration, using:
- `advanced_configuration.tls_cipher_config_mode = "CUSTOM"` to enable custom cipher lists.
- `advanced_configuration.minimum_enabled_tls_protocol = "TLS1_2"` to enforce TLS 1.2 as the minimum protocol version.

Variables required: 
- `atlas_org_id`: ID of the Atlas organization
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `provider_name`: Name of provider to use for cluster (TENANT, AWS, GCP)
- `provider_instance_size_name`: Size of the cluster (Free: M0, Dedicated: M10+.)