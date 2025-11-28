# MongoDB Atlas Provider — Advanced Cluster (TLS 1.3 Configuration)

This example creates a project and a sharded cluster with TLS 1.3 cipher configuration, using:
- `advanced_configuration.tls_cipher_config_mode = "CUSTOM"` to enable custom cipher lists.
- `advanced_configuration.minimum_enabled_tls_protocol = "TLS1_3"` to enforce TLS 1.3 as the minimum protocol version.

Variables required: 
- `atlas_org_id`: ID of the Atlas organization
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `provider_name`: Name of provider to use for cluster (TENANT, AWS, GCP)
- `provider_instance_size_name`: Size of the cluster (Free: M0, Dedicated: M10+.)