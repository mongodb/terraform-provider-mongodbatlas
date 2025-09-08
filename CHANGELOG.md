## (Unreleased)

## 1.41.0 (September 08, 2025)

ENHANCEMENTS:

* data-source/mongodbatlas_cloud_provider_access_setup: Adds support for GCP as a Cloud Provider. ([#3637](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3637))
* data-source/mongodbatlas_encryption_at_rest: Supports role_id in google_cloud_kms_config ([#3636](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3636))
* resource/mongodbatlas_cloud_provider_access_authorization: Changes to `project_id` or `role_id` will now result in the destruction and recreation of the authorization resource ([#3646](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3646))
* resource/mongodbatlas_cloud_provider_access_authorization: Supports GCP cloud provider ([#3639](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3639))
* resource/mongodbatlas_cloud_provider_access_setup: Adds long running operation support for GCP ([#3644](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3644))
* resource/mongodbatlas_cloud_provider_access_setup: Adds support for GCP as a Cloud Provider. ([#3637](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3637))
* resource/mongodbatlas_encryption_at_rest: Supports role_id in google_cloud_kms_config ([#3636](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3636))

BUG FIXES:

* resource/mongodbatlas_advanced_cluster: Fixes `Value Conversion Error` when replication_specs are unknown ([#3652](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3652))

## 1.40.0 (August 21, 2025)

ENHANCEMENTS:

* data-source/mongodbatlas_stream_privatelink_endpoint: Support S3 PrivateLink Endpoints for Atlas Stream Processing ([#3554](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3554))
* resource/mongodbatlas_stream_privatelink_endpoint: Support S3 PrivateLink Endpoints for Atlas Stream Processing ([#3554](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3554))

## 1.39.0 (July 24, 2025)

NOTES:

* New environment variables: We added support for the `MONGODB_ATLAS_PUBLIC_API_KEY` and `MONGODB_ATLAS_PRIVATE_API_KEY` environment variables which are widely used across the MongoDB ecosystem. ([#3505](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3505))

ENHANCEMENTS:

* data-source/mongodbatlas_federated_database_instance: Adds `azure` attribute to support reading federated databases with Azure cloud provider configuration ([#3484](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3484))
* data-source/mongodbatlas_federated_database_instances: Adds `azure` attribute to support reading federated databases with Azure cloud provider configuration ([#3484](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3484))
* resource/mongodbatlas_federated_database_instance: Adds `azure` attribute to allow the creation of federated databases with Azure cloud provider configuration ([#3484](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3484))

BUG FIXES:

* resource/mongodbatlas_organization: Sets org_id on import ([#3513](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3513))

## 1.38.0 (July 10, 2025)

NOTES:

* data-source/mongodbatlas_stream_connections: Deprecates the `id` attribute as it is a random assigned value which should not be used ([#3476](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3476))
* data-source/mongodbatlas_stream_instances: Deprecates the `id` attribute as it is a random assigned value which should not be used ([#3476](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3476))

FEATURES:

* **New Data Source:** `data-source/mongodbatlas_api_key_project_assignment` ([#3461](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3461))
* **New Data Source:** `data-source/mongodbatlas_api_key_project_assignments` ([#3461](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3461))
* **New Resource:** `resource/mongodbatlas_api_key_project_assignment` ([#3461](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3461))

ENHANCEMENTS:

* data-source/mongodbatlas_third_party_integration Adds `send_user_provided_resource_tags` attribute to support sending $querystats to DataDog ([#3454](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3454))
* data-source/mongodbatlas_third_party_integrations Adds `send_user_provided_resource_tags` attribute to support sending $querystats to DataDog ([#3454](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3454))
* resource/mongodbatlas_organization: Adds import support ([#3475](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3475))
* resource/mongodbatlas_third_party_integration Adds `send_user_provided_resource_tags` attribute to support sending $querystats to DataDog ([#3454](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3454))

BUG FIXES:

* data-source/mongodbatlas_cloud_backup_snapshot_export_buckets: Fix pagination when `items_per_page` or `page_num` are not set ([#3459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3459))
* data-source/mongodbatlas_cloud_backup_snapshot_export_jobs: Fix pagination when `items_per_page` or `page_num` are not set ([#3459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3459))
* data-source/mongodbatlas_cloud_backup_snapshot_restore_jobs: Fix pagination when `items_per_page` or `page_num` are not set ([#3459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3459))
* data-source/mongodbatlas_cloud_backup_snapshots: Fix pagination when `items_per_page` or `page_num` are not set ([#3459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3459))
* data-source/mongodbatlas_federated_settings_org_configs: Fix pagination when `items_per_page` or `page_num` are not set ([#3459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3459))
* data-source/mongodbatlas_organizations: Fix pagination when `items_per_page` or `page_num` are not set ([#3459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3459))

## 1.37.0 (June 26, 2025)

ENHANCEMENTS:

* data-source/mongodbatlas_stream_connection Adds `cluster_project_id` to allow connections to clusters in other projects within an organization ([#3424](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3424))
* resource/mongodbatlas_stream_connection Adds `cluster_project_id` to allow connections to clusters in other projects within an organization ([#3424](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3424))

## 1.36.0 (June 11, 2025)

FEATURES:

* **New Data Source:** `mongodbatlas_stream_account_details` ([#3364](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3364))

## 1.35.1 (May 29, 2025)

BUG FIXES:

* provider: Fixes Realm Client authentication ([#3362](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3362))

## 1.35.0 (May 28, 2025)

ENHANCEMENTS:

* resource/mongodbatlas_advanced_cluster (preview provider 2.0.0): Adds `delete_on_create_timeout` a flag that indicates whether to delete the cluster if the cluster creation times out, default is false ([#3333](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3333))
* resource/mongodbatlas_advanced_cluster: Adds `delete_on_create_timeout` a flag that indicates whether to delete the cluster if the cluster creation times out, default is false ([#3333](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3333))
* resource/mongodbatlas_search_deployment: Adds `delete_on_create_timeout` a flag that indicates whether to delete the search deployment if the search deployment creation times out, default is false ([#3344](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3344))

BUG FIXES:

* resource/mongodbatlas_private_endpoint_regional_mode: Increases update wait time so cluster connection strings are updated ([#3320](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3320))

## 1.34.0 (April 30, 2025)

ENHANCEMENTS:

* data-source/mongodbatlas_database_user: Adds `description` field ([#3280](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3280))
* data-source/mongodbatlas_database_users: Adds `description` field ([#3280](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3280))
* data-source/mongodbatlas_maintenance_window: Adds `protected_hours` and `time_zone_id` ([#3195](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3195))
* resource/mongodbatlas_database_user: Adds `description` field ([#3280](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3280))
* resource/mongodbatlas_maintenance_window: Adds `protected_hours` and `time_zone_id` ([#3195](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3195))

BUG FIXES:

* resource/mongodbatlas_auditing: Fixes JSON comparison in `audit_filter` field ([#3302](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3302))

## 1.33.0 (April 16, 2025)

NOTES:

* data-source/mongodbatlas_resource_policies: Enables usage without preview environment flag ([#3276](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3276))
* data-source/mongodbatlas_resource_policy: Enables usage without preview environment flag ([#3276](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3276))
* resource/mongodbatlas_resource_policy: Enables usage without preview environment flag ([#3276](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3276))

ENHANCEMENTS:

* data-source/mongodbatlas_organization: Adds `security_contact` attribute ([#3263](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3263))
* data-source/mongodbatlas_organizations: Adds `security_contact` attribute ([#3263](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3263))
* data-source/mongodbatlas_third_party_integration: Adds support for `send_collection_latency_metrics` and `send_database_metrics` for Datadog integrations ([#3259](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3259))
* data-source/mongodbatlas_third_party_integrations: Adds support for `send_collection_latency_metrics` and `send_database_metrics` for Datadog integrations ([#3259](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3259))
* resource/mongodbatlas_organization: Adds `security_contact` attribute ([#3263](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3263))
* resource/mongodbatlas_third_party_integration: Adds support for `send_collection_latency_metrics` and `send_database_metrics` for Datadog integrations ([#3259](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3259))

## 1.32.0 (April 09, 2025)

ENHANCEMENTS:

* data-source/mongodbatlas_encryption_at_rest: Adds `enabled_for_search_nodes` attribute ([#3142](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3142))
* data-source/mongodbatlas_resource_policies: Adds support for the new `description` field ([#3214](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3214))
* data-source/mongodbatlas_resource_policy: Adds support for the new `description` field ([#3214](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3214))
* data-source/mongodbatlas_search_deployment: Adds `encryption_at_rest_provider` computed attribute ([#3152](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3152))
* resource/mongodbatlas_encryption_at_rest: Adds `enabled_for_search_nodes` attribute ([#3142](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3142))
* resource/mongodbatlas_resource_policy: Adds support for the new `description` field ([#3214](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3214))
* resource/mongodbatlas_search_deployment: Adds `encryption_at_rest_provider` computed attribute ([#3152](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3152))
* resource/mongodbatlas_search_deployment: Adds `skip_wait_on_update` to avoid waiting for completion of update operations ([#3237](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3237))
* resource/mongodbatlas_stream_processor: Adds update support ([#3180](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3180))

## 1.31.0 (March 25, 2025)

ENHANCEMENTS:

* data-source/mongodbatlas_organization: Adds support for `skip_default_alerts_settings` setting ([#2933](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2933))
* data-source/mongodbatlas_organizations: Adds support for `skip_default_alerts_settings` setting ([#2933](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2933))
* resource/mongodbatlas_organization: Adds support for `skip_default_alerts_settings` setting ([#2933](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2933))

## 1.30.0 (March 20, 2025)

ENHANCEMENTS:

* data-source/mongodbatlas_stream_connection: Adds `Https` connection ([#3150](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3150))
* data-source/mongodbatlas_stream_connections: Adds `Https` connection ([#3150](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3150))
* resource/mongodbatlas_cloud_backup_snapshot: Adds `timeouts` attribute for create operation ([#3171](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3171))
* resource/mongodbatlas_cloud_backup_snapshot: Adjusts creation default timeout from 20 minutes to 1 hour ([#3171](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3171))
* resource/mongodbatlas_stream_connection: Adds `Https` connection ([#3150](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3150))
* resource/mongodbatlas_stream_privatelink_endpoint: Adds `error_message`, `interface_endpoint_name`, and `provider_account_id` attributes ([#3161](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3161))
* resource/mongodbatlas_stream_privatelink_endpoint: Adds support for AWS MSK clusters ([#3179](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3179))

BUG FIXES:

* data-source/mongodbatlas_global_cluster_config: Adds support for reading clusters with independent shard scaling ([#3177](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3177))
* resource/mongodbatlas_advanced_cluster (preview provider 2.0.0): Avoids error when removing `read_only_specs` in `region_configs` that does not define `electable_specs` ([#3162](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3162))
* resource/mongodbatlas_global_cluster_config: Adds support for reading clusters with independent shard scaling ([#3177](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3177))

## 1.29.0 (March 12, 2025)

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: Available as Preview of MongoDB Atlas Provider 2.0.0 ([#3147](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3147))
* data-source/mongodbatlas_advanced_clusters: Available as Preview of MongoDB Atlas Provider 2.0.0 ([#3147](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3147))
* data-source/mongodbatlas_stream_connection: Adds `AWSLambda` connection ([#3085](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3085))
* resource/mongodbatlas_advanced_cluster: Available as Preview of MongoDB Atlas Provider 2.0.0 ([#3147](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3147))
* resource/mongodbatlas_stream_connection: Adds `AWSLambda` connection ([#3085](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3085))

BUG FIXES:

* data-source/mongodbatlas_organizations: Avoids nil pointer error when individual getOrganizationSettings API call fails ([#3118](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3118))
* resource/mongodbatlas_backup_compliance_policy: Changes `on_demand_policy_item` attribute from `required` to `optional` ([#3119](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3119))

## 1.28.0 (February 25, 2025)

ENHANCEMENTS:

* resource/mongodbatlas_encryption_at_rest: Adds support for `aws_kms_config.requirePrivateNetworking` ([#2951](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2951))

## 1.27.0 (February 20, 2025)

NOTES:

* data-source/mongodbatlas_serverless_instance: Deprecates data source ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))
* data-source/mongodbatlas_serverless_instances: Deprecates data source ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))
* data-source/mongodbatlas_shared_tier_restore_job: Deprecates data source ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))
* data-source/mongodbatlas_shared_tier_restore_jobs: Deprecates data source ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))
* data-source/mongodbatlas_shared_tier_snapshot: Deprecates data source ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))
* data-source/mongodbatlas_shared_tier_snapshot: Deprecates data source ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))
* resource/mongodbatlas_advanced_cluster: Deprecates `M2` and `M5` instance size for the attribute `instance_size` inside of `analytics_specs`, `electable_specs` and `read_only_specs` ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))
* resource/mongodbatlas_cluster: Deprecates `M2` and `M5` instance size values for the attribute `provider_instance_size_name` ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))
* resource/mongodbatlas_serverless_instance: Deprecates resource ([#3012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3012))

FEATURES:

* **New Data Source:** `mongodbatlas_flex_restore_job` ([#3041](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3041))
* **New Data Source:** `mongodbatlas_flex_restore_jobs` ([#3041](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3041))
* **New Data Source:** `mongodbatlas_flex_snapshot` ([#3036](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3036))
* **New Data Source:** `mongodbatlas_flex_snapshots` ([#3036](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3036))

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: This data source can now read Flex clusters ([#3001](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3001))
* data-source/mongodbatlas_advanced_clusters: This data source can now read Flex clusters ([#3001](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3001))
* data-source/mongodbatlas_flex_cluster: Reaches GA (General Availability) ([#3003](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3003))
* data-source/mongodbatlas_flex_cluster: Reaches GA (General Availability) ([#3003](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3003))
* resource/mongodbatlas_advanced_cluster: This resource can now create, read, update, and delete Flex clusters ([#3001](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3001))
* resource/mongodbatlas_flex_cluster: Reaches GA (General Availability) ([#3003](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3003))
* resource/mongodbatlas_global_cluster_config: Supports update operation ([#3060](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3060))

BUG FIXES:

* resource/mongodbatlas_alert_configuration: Removes UseStateForUnknown plan modifier for interval_min ([#3051](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3051))
* resource/mongodbatlas_database_user: Avoids error in read if resource no longer exists ([#3069](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3069))
* resource/mongodbatlas_maintenance_window: Avoids error in creation when `hour_of_day` is set to zero or not defined ([#3086](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3086))

## 1.26.1 (February 07, 2025)

BUG FIXES:

* resource/mongodbatlas_advanced_cluster: Adds `PENDING` status for update and delete operations ([#3034](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3034))
* resource/mongodbatlas_cluster: Adds `PENDING` status for update and delete operations ([#3034](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3034))

## 1.26.0 (January 23, 2025)

FEATURES:

* **New Data Source:** `mongodbatlas_stream_privatelink_endpoint` ([#2897](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2897))
* **New Data Source:** `mongodbatlas_stream_privatelink_endpoints` ([#2897](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2897))
* **New Resource:** `mongodbatlas_stream_privatelink_endpoint` ([#2890](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2890))

ENHANCEMENTS:

* resource/mongodbatlas_backup_compliance_policy: Adds support for disabling Backup Compliance Policy on resource delete ([#2953](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2953))
* resource/mongodbatlas_stream_connection: Supports Privatelink networking access type for Kafka Stream Connections ([#2940](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2940))

BUG FIXES:

* resource/mongodbatlas_search_index: Don't send empty `analyzers` attribute to Atlas ([#2994](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2994))

## 1.25.0 (January 07, 2025)

NOTES:

* resource/mongodbatlas_cloud_backup_snapshot_export_bucket: Deprecates `tenant_id` argument as the `mongodbatlas_cloud_provider_access_authorization.azure.tenant_id` is used instead ([#2932](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2932))
* resource/mongodbatlas_cloud_backup_snapshot_export_job: Changes `custom_data` changed attribute from required -> optional ([#2929](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2929))

ENHANCEMENTS:

* data-source/mongodbatlas_project_ip_addresses: Adds support for `future_inbound` and `future_outbound` fields ([#2934](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2934))
* data-source/mongodbatlas_stream_connection: Adds `networking` attribute ([#2474](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2474))
* data-source/mongodbatlas_stream_connections: Adds `networking` attribute ([#2474](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2474))
* resource/mongodbatlas_stream_connection: Adds `networking` attribute ([#2474](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2474))

BUG FIXES:

* data-source/mongodbatlas_team: Fixes pagination logic when retrieved users of a team ([#2919](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2919))
* resource/mongodbatlas_database_user: Avoids import error for database_user when both username and auth database contain hyphens ([#2928](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2928))
* resource/mongodbatlas_team: Fixes pagination logic when retrieved users of a team ([#2919](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2919))

## 1.24.0 (December 20, 2024)

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: Adds `advanced_configuration.0.tls_cipher_config_mode` and `advanced_configuration.0.custom_openssl_cipher_config_tls12` attribute ([#2872](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2872))
* data-source/mongodbatlas_advanced_clusters: Adds `advanced_configuration.0.tls_cipher_config_mode` and `advanced_configuration.0.custom_openssl_cipher_config_tls12` attribute ([#2872](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2872))
* data-source/mongodbatlas_cluster: Adds `advanced_configuration.0.tls_cipher_config_mode` and `advanced_configuration.0.custom_openssl_cipher_config_tls12` attribute ([#2872](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2872))
* data-source/mongodbatlas_cluster: Adds `advanced_configuration.0.tls_cipher_config_mode` and `advanced_configuration.0.custom_openssl_cipher_config_tls12` attribute ([#2872](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2872))
* resource/mongodbatlas_advanced_cluster: Adds `advanced_configuration.0.tls_cipher_config_mode` and `advanced_configuration.0.custom_openssl_cipher_config_tls12` attribute ([#2872](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2872))
* resource/mongodbatlas_cluster: Adds `advanced_configuration.0.tls_cipher_config_mode` and `advanced_configuration.0.custom_openssl_cipher_config_tls12` attribute ([#2872](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2872))

## 1.23.0 (December 17, 2024)

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: Adds `advanced_configuration.0.default_max_time_ms` attribute ([#2825](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2825))
* data-source/mongodbatlas_advanced_cluster: Adds `pinned_fcv` attribute ([#2789](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2789))
* data-source/mongodbatlas_advanced_clusters: Adds `advanced_configuration.0.default_max_time_ms` attribute ([#2825](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2825))
* data-source/mongodbatlas_advanced_clusters: Adds `pinned_fcv` attribute ([#2789](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2789))
* data-source/mongodbatlas_cluster: Adds `pinned_fcv` attribute ([#2817](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2817))
* data-source/mongodbatlas_clusters: Adds `pinned_fcv` attribute ([#2817](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2817))
* resource/mongodbatlas_advanced_cluster: Adds `advanced_configuration.0.default_max_time_ms` attribute ([#2825](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2825))
* resource/mongodbatlas_advanced_cluster: Adds `pinned_fcv` attribute ([#2789](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2789))
* resource/mongodbatlas_advanced_cluster: Adjusts create operation to support cluster tier auto scaling per shard. ([#2836](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2836))
* resource/mongodbatlas_advanced_cluster: Adjusts update operation to support cluster tier auto scaling per shard. ([#2814](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2814))
* resource/mongodbatlas_cluster: Adds `pinned_fcv` attribute ([#2817](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2817))

BUG FIXES:

* data-source/mongodbatlas_advanced_cluster: `mongo_db_major_version` attribute is populated with binary version when FCV pin is active ([#2789](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2789))
* data-source/mongodbatlas_advanced_clusters: `mongo_db_major_version` attribute is populated with binary version when FCV pin is active ([#2789](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2789))
* data-source/mongodbatlas_cluster: `mongo_db_major_version` attribute is populated with binary version when FCV pin is active ([#2817](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2817))
* data-source/mongodbatlas_clusters: `mongo_db_major_version` attribute is populated with binary version when FCV pin is active ([#2817](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2817))
* resource/mongodbatlas_advanced_cluster: `mongo_db_major_version` attribute is populated with binary version when FCV pin is active ([#2789](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2789))
* resource/mongodbatlas_cluster: `mongo_db_major_version` attribute is populated with binary version when FCV pin is active ([#2817](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2817))
* resource/mongodbatlas_search_index: Fixes resource create and update when `wait_for_index_build_completion` attribute is used ([#2887](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2887))

## 1.22.0 (November 26, 2024)

NOTES:

* data-source/mongodbatlas_organization: Adds new `gen_ai_features_enabled` attribute ([#2724](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2724))
* data-source/mongodbatlas_organizations: Adds new `gen_ai_features_enabled` attribute ([#2724](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2724))
* data-source/mongodbatlas_privatelink_endpoint_service_serverless: Deprecates data source ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* data-source/mongodbatlas_privatelink_endpoints_service_serverless: Deprecates data source ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* data-source/mongodbatlas_serverless_instance: Deprecates `auto_indexing` attribute ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* data-source/mongodbatlas_serverless_instance: Deprecates `continuous_backup_enabled` attribute ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* data-source/mongodbatlas_serverless_instances: Deprecates `auto_indexing` attribute ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* data-source/mongodbatlas_serverless_instances: Deprecates `continuous_backup_enabled` attribute ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* resource/mongodbatlas_organization: Adds new `gen_ai_features_enabled` attribute ([#2724](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2724))
* resource/mongodbatlas_privatelink_endpoint_serverless: Deprecates resource ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* resource/mongodbatlas_privatelink_endpoint_service_serverless: Deprecates resource ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* resource/mongodbatlas_serverless_instance: Deprecates `auto_indexing` attribute ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))
* resource/mongodbatlas_serverless_instance: Deprecates `continuous_backup_enabled` attribute ([#2742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2742))

FEATURES:

* **New Data Source:** `mongodbatlas_flex_cluster` ([#2738](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2738))
* **New Data Source:** `mongodbatlas_flex_clusters` ([#2767](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2767))
* **New Resource:** `mongodbatlas_flex_cluster` ([#2716](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2716))

ENHANCEMENTS:

* data-source/mongodbatlas_cloud_backup_snapshot_restore_job: Adds `failed` attribute ([#2781](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2781))
* data-source/mongodbatlas_cloud_backup_snapshot_restore_jobs: Adds `failed` attribute ([#2781](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2781))
* resource/mongodbatlas_cloud_backup_snapshot_restore_job: Adds `failed` attribute ([#2781](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2781))
* resource/mongodbatlas_network_peering: Improve error message when networking peering reaches a failed status ([#2766](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2766))
* resource/mongodbatlas_privatelink_endpoint: Improves error message when privatelink endpoint returns error after POST ([#2803](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2803))
* resource/mongodbatlas_privatelink_endpoint_service: Decreases delay time when creating or deleting a resource ([#2819](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2819))
* resource/mongodbatlas_privatelink_endpoint_service: Improves error message when privatelink endpoint service returns error after POST ([#2803](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2803))

## 1.21.4 (October 29, 2024)

NOTES:

* data-source/mongodbatlas_resource_policies: Deprecates `resource_policies` attribute ([#2740](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2740))

ENHANCEMENTS:

* data-source/mongodbatlas_resource_policies: Adds `results` attribute ([#2740](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2740))

BUG FIXES:

* resource/mongodbatlas_project_api_key: Validates `project_id` are unique across `project_assignment` blocks and fixes update issues with error `API_KEY_ALREADY_IN_GROUP` ([#2737](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2737))

## 1.21.3 (October 25, 2024)

NOTES:

* data-source/mongodbatlas_project: Deprecates `is_slow_operation_thresholding_enabled`. Attribute will be supported in a separate data source as it requires different set of permissions ([#2731](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2731))
* data-source/mongodbatlas_projects: Deprecates `is_slow_operation_thresholding_enabled`. Attribute will be supported in a separate data source as it requires different set of permissions ([#2731](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2731))
* resource/mongodbatlas_project: Deprecates `is_slow_operation_thresholding_enabled`. Attribute will be supported in a separate resource as it requires different set of permissions ([#2731](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2731))

BUG FIXES:

* data-source/mongodbatlas_project: Avoids error when user doesn't have project owner permission ([#2731](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2731))
* data-source/mongodbatlas_projects: Avoids error when user doesn't have project owner permission ([#2731](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2731))
* resource/mongodbatlas_project: Avoids error when user doesn't have project owner permission ([#2731](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2731))

## 1.21.2 (October 22, 2024)

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: Adds new `config_server_management_mode` and `config_server_type` fields ([#2670](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2670))
* data-source/mongodbatlas_advanced_clusters: Adds new `config_server_management_mode` and `config_server_type` fields ([#2670](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2670))
* data-source/mongodbatlas_project: Adds `is_slow_operation_thresholding_enabled` attribute ([#2698](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2698))
* data-source/mongodbatlas_projects: Adds `is_slow_operation_thresholding_enabled` attribute ([#2698](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2698))
* resource/mongodbatlas_advanced_cluster: Adds new `config_server_management_mode` and `config_server_type` fields ([#2670](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2670))
* resource/mongodbatlas_project: Adds `is_slow_operation_thresholding_enabled` attribute ([#2698](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2698))

BUG FIXES:

* resource/mongodbatlas_event_trigger: Always includes `disabled` in the PUT payload ([#2690](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2690))
* resource/mongodbatlas_organization: Avoids inconsistent result returned by provider when `USER_NOT_FOUND` ([#2684](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2684))
* resource/mongodbatlas_search_deployment: Fixes inconsistent result for a multi-region cluster that always uses a single spec. ([#2685](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2685))

## 1.21.1 (October 09, 2024)

BUG FIXES:

* resource/mongodbatlas_team: Fixes update logic of `usernames` attribute ensuring team is never emptied ([#2669](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2669))

## 1.21.0 (October 07, 2024)

NOTES:

* data-source/mongodbatlas_global_cluster_config: Deprecates `custom_zone_mapping` in favor of `custom_zone_mapping_zone_id` ([#2637](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2637))
* resource/mongodbatlas_global_cluster_config: Deprecates `custom_zone_mapping` in favor of `custom_zone_mapping_zone_id` ([#2637](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2637))

FEATURES:

* **New Data Source:** `data-source/mongodbatlas_resource_policies` ([#2598](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2598))
* **New Data Source:** `data-source/mongodbatlas_resource_policy` ([#2598](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2598))
* **New Resource:** `resource/mongodbatlas_resource_policy` ([#2585](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2585))

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: Supports `redact_client_log_data` attribute ([#2600](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2600))
* data-source/mongodbatlas_advanced_clusters: Supports `redact_client_log_data` attribute ([#2600](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2600))
* data-source/mongodbatlas_cluster: Supports `redact_client_log_data` attribute ([#2601](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2601))
* data-source/mongodbatlas_clusters: Supports `redact_client_log_data` attribute ([#2601](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2601))
* data-source/mongodbatlas_global_cluster_config: Adds `custom_zone_mapping_zone_id` attribute ([#2637](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2637))
* resource/mongodbatlas_advanced_cluster: Supports `redact_client_log_data` attribute ([#2600](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2600))
* resource/mongodbatlas_cluster: Supports `redact_client_log_data` attribute ([#2601](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2601))
* resource/mongodbatlas_global_cluster_config: Adds `custom_zone_mapping_zone_id` attribute ([#2637](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2637))

BUG FIXES:

* resource/mongodbatlas_advanced_cluster: Enforces `priority` descending order in `region_configs` avoiding potential non-empty plans after apply ([#2640](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2640))

## 1.20.0 (September 20, 2024)

BREAKING CHANGES:

* data-source/mongodbatlas_cloud_backup_snapshot_export_job: Removes `err_msg` attribute ([#2617](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2617))
* data-source/mongodbatlas_cloud_backup_snapshot_export_jobs: Removes `err_msg` attribute ([#2617](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2617))
* data-source/mongodbatlas_federated_database_instance: Removes `storage_stores.#.cluster_id` attribute ([#2617](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2617))
* data-source/mongodbatlas_federated_database_instances: Removes `storage_stores.#.cluster_id` attribute ([#2617](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2617))
* resource/mongodbatlas_cloud_backup_snapshot_export_job: Removes `err_msg` attribute ([#2617](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2617))
* resource/mongodbatlas_federated_database_instance: Removes `storage_stores.#.cluster_id` attribute ([#2617](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2617))

NOTES:

* data-source/mongodbatlas_data_lake_pipeline: Data Lake is deprecated. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation ([#2599](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2599))
* data-source/mongodbatlas_data_lake_pipeline_run: Data Lake is deprecated. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation ([#2599](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2599))
* data-source/mongodbatlas_data_lake_pipeline_runs: Data Lake is deprecated. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation ([#2599](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2599))
* data-source/mongodbatlas_data_lake_pipelines: Data Lake is deprecated. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation ([#2599](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2599))
* resource/mongodbatlas_alert_configuration: Updates `notification.#.integration_id` to be Optional & Computed ([#2603](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2603))
* resource/mongodbatlas_data_lake_pipeline: Data Lake is deprecated. To learn more, see https://dochub.mongodb.org/core/data-lake-deprecation ([#2599](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2599))

FEATURES:

* **New Data Source:** `data-source/mongodbatlas_mongodb_employee_access_grant` ([#2591](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2591))
* **New Resource:** `resource/mongodbatlas_mongodb_employee_access_grant` ([#2591](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2591))

BUG FIXES:

* resource/mongodbatlas_advanced_cluster: Set `advanced_configuration.change_stream_options_pre_and_post_images_expire_after_seconds` only for compatible MongoDB versions ([#2592](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2592))
* resource/mongodbatlas_advanced_cluster: Supports using decimal in advanced_configuration `oplog_min_retention_hours` ([#2604](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2604))
* resource/mongodbatlas_cluster: Set `advanced_configuration.change_stream_options_pre_and_post_images_expire_after_seconds` only for compatible MongoDB versions ([#2592](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2592))
* resource/mongodbatlas_cluster: Supports using decimal in advanced_configuration `oplog_min_retention_hours` ([#2604](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2604))
* resource/mongodbatlas_stream_processor: Error during create should only show one error message and required actions ([#2590](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2590))

## 1.19.0 (September 11, 2024)

NOTES:

* data-source/mongodbatlas_project: Deprecates the `ip_addresses` attribute. Use the new `mongodbatlas_project_ip_addresses` data source to obtain this information instead. ([#2541](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2541))
* data-source/mongodbatlas_projects: Deprecates the `ip_addresses` attribute. Use the new `mongodbatlas_project_ip_addresses` data source to obtain this information instead. ([#2541](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2541))
* resource/mongodbatlas_project: Deprecates the `ip_addresses` attribute. Use the new `mongodbatlas_project_ip_addresses` data source to obtain this information instead. ([#2541](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2541))

FEATURES:

* **New Data Source:** `data-source/mongodbatlas_encryption_at_rest` ([#2538](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2538))
* **New Data Source:** `data-source/mongodbatlas_encryption_at_rest_private_endpoint` ([#2527](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2527))
* **New Data Source:** `data-source/mongodbatlas_encryption_at_rest_private_endpoints` ([#2536](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2536))
* **New Data Source:** `data-source/mongodbatlas_project_ip_addresses` ([#2533](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2533))
* **New Data Source:** `data-source/mongodbatlas_stream_processor` ([#2497](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2497))
* **New Data Source:** `data-source/mongodbatlas_stream_processors` ([#2566](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2566))
* **New Resource:** `mongodbatlas_stream_processor` ([#2501](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2501))
* **New Resource:** `resource/mongodbatlas_encryption_at_rest_private_endpoint` ([#2512](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2512))

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: Supports change_stream_options_pre_and_post_images_expire_after_seconds attribute ([#2528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2528))
* data-source/mongodbatlas_advanced_cluster: Supports change_stream_options_pre_and_post_images_expire_after_seconds attribute ([#2528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2528))
* data-source/mongodbatlas_advanced_cluster: supports replica_set_scaling_strategy attribute ([#2539](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2539))
* data-source/mongodbatlas_advanced_clusters: supports replica_set_scaling_strategy attribute ([#2539](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2539))
* data-source/mongodbatlas_cluster: Supports change_stream_options_pre_and_post_images_expire_after_seconds attribute ([#2528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2528))
* data-source/mongodbatlas_clusters: Supports change_stream_options_pre_and_post_images_expire_after_seconds attribute ([#2528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2528))
* resource/mongodbatlas_advanced_cluster: Supports change_stream_options_pre_and_post_images_expire_after_seconds attribute ([#2528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2528))
* resource/mongodbatlas_advanced_cluster: supports replica_set_scaling_strategy attribute ([#2539](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2539))
* resource/mongodbatlas_cluster: Supports change_stream_options_pre_and_post_images_expire_after_seconds attribute ([#2528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2528))
* resource/mongodbatlas_encryption_at_rest: Adds `aws_kms_config.0.valid`, `azure_key_vault_config.0.valid` and `google_cloud_kms_config.0.valid` attribute ([#2538](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2538))
* resource/mongodbatlas_encryption_at_rest: Adds new `azure_key_vault_config.#.require_private_networking` field to enable connection to Azure Key Vault over private networking ([#2509](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2509))

BUG FIXES:

* data-source/mongodbatlas_advanced_clusters: Sets correct `zone_id` when `use_replication_spec_per_shard` is false ([#2568](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2568))

## 1.18.1 (August 26, 2024)

NOTES:

* resource/mongodbatlas_advanced_cluster: Documentation adjustment in resource and migration guide to clarify potential `Internal Server Error` when applying updates with new sharding configuration ([#2525](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2525))

## 1.18.0 (August 14, 2024)

BREAKING CHANGES:

* data-source/mongodbatlas_cloud_backup_snapshot_export_bucket: Changes `id` attribute from optional to computed only ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* data-source/mongodbatlas_cloud_backup_snapshot_export_job: Changes `id` attribute from optional to computed only ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* data-source/mongodbatlas_cloud_backup_snapshot_restore_job: Removes `created_at` attribute ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* data-source/mongodbatlas_cloud_backup_snapshot_restore_job: Removes `job_id` attribute and defines `snapshot_restore_job_id` attribute as required ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* data-source/mongodbatlas_cloud_backup_snapshot_restore_jobs: Removes `created_at` attribute ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* data-source/mongodbatlas_federated_settings_identity_providers: Removes `page_num` and `items_per_page` attributes ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* data-source/mongodbatlas_privatelink_endpoint_service: Removes `endpoints.*.service_attachment_name` attribute ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* data-source/mongodbatlas_third_party_integration: Removes `scheme` attribute ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* data-source/mongodbatlas_third_party_integrations: Removes `scheme` attribute ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* resource/mongodbatlas_cloud_backup_snapshot_restore_job: Removes `created_at` attribute ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* resource/mongodbatlas_privatelink_endpoint_service: Removes `endpoints.*.service_attachment_name` attribute ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))
* resource/mongodbatlas_third_party_integration: Removes `scheme` attribute ([#2499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2499))

NOTES:

* data-source/mongodbatlas_advanced_cluster: Deprecates `replication_specs.#.id`, `replication_specs.#.num_shards`, `disk_size_gb`, `advanced_configuration.0.default_read_concern`, and  `advanced_configuration.0.fail_index_key_too_long`. To learn more, see the [1.18.0 Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide). ([#2420](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2420))
* data-source/mongodbatlas_advanced_clusters: Deprecates `replication_specs.#.id`, `replication_specs.#.num_shards`, `disk_size_gb`, `advanced_configuration.0.default_read_concern`, and  `advanced_configuration.0.fail_index_key_too_long`. To learn more, see the [1.18.0 Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide). ([#2420](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2420))
* resource/mongodbatlas_advanced_cluster: Deprecates `replication_specs.#.id`, `replication_specs.#.num_shards`, `disk_size_gb`, `advanced_configuration.0.default_read_concern`, and  `advanced_configuration.0.fail_index_key_too_long`. To learn more, see the [1.18.0 Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide). ([#2420](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2420))
* resource/mongodbatlas_advanced_cluster: Using this new version impacts the possibility of editing the definition of multi shard clusters in the Atlas UI. This impact is limited to the first weeks of September. ([#2478](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2478))

FEATURES:

* **New Guide:** [Migration Guide: Advanced Cluster New Sharding Schema](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema). This enables Independent Shard Scaling. ([#2505](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2505))
* **New Guide:** [Migration Guide: Cluster to Advanced Cluster](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide) ([#2505](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2505))

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: Adds `use_replication_spec_per_shard`, `replication_specs.*.zone_id`, and `replication_specs.*.region_configs.*.(electable_specs|analytics_specs|read_only_specs).disk_size_gb` attributes. To learn more, see the [1.18.0 Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide) and data source documentation. ([#2478](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2478))
* data-source/mongodbatlas_advanced_clusters: Adds `use_replication_spec_per_shard`, `replication_specs.*.zone_id`, and `replication_specs.*.region_configs.*.(electable_specs|analytics_specs|read_only_specs).disk_size_gb` attributes. To learn more, see the [1.18.0 Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide) and data source documentation. ([#2478](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2478))
* data-source/mongodbatlas_cloud_backup_schedule: Adds new `use_zone_id_for_copy_settings` and `copy_settings.#.zone_id` attributes and deprecates `copy_settings.#.replication_spec_id`. These new attributes enable you to reference cluster zones using independent shard scaling, which no longer supports `replication_spec.*.id` ([#2464](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2464))
* data-source/mongodbatlas_cloud_backup_snapshot_export_bucket: Adds Azure support ([#2486](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2486))
* data-source/mongodbatlas_cloud_backup_snapshot_export_buckets: Adds Azure support ([#2486](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2486))
* resource/mongodbatlas_advanced_cluster: Adds `replication_specs.*.zone_id` and `replication_specs.*.region_configs.*.(electable_specs|analytics_specs|read_only_specs).disk_size_gb` attributes. To learn more, see the [1.18.0 Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.18.0-upgrade-guide) and resource documentation. ([#2478](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2478))
* resource/mongodbatlas_advanced_cluster: Supports defining cluster shards with independent `replication_specs` objects. This feature enables defining independent scaled shards. To learn more, see the [Advanced Cluster New Sharding Schema Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema). ([#2478](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2478))
* resource/mongodbatlas_cloud_backup_schedule: Adds `copy_settings.#.zone_id` and deprecates `copy_settings.#.replication_spec_id` for referencing zones of a cluster. This enables referencing zones of clusters using independent shard scaling which no longer support `replication_spec.*.id`. ([#2459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2459))
* resource/mongodbatlas_cloud_backup_snapshot_export_bucket: Adds Azure support ([#2486](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2486))

## 1.17.6 (August 07, 2024)

BUG FIXES:

* resource/mongodbatlas_backup_compliance_policy: Fixes an issue where the update operation modified attributes that were not supposed to be modified" ([#2480](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2480))

## 1.17.5 (July 30, 2024)

NOTES:

* data-source/mongodbatlas_cloud_backup_snapshot_export_job: Deprecates the `err_msg` attribute ([#2436](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2436))
* data-source/mongodbatlas_cloud_backup_snapshot_export_jobs: Deprecates the `err_msg` attribute ([#2436](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2436))
* resource/mongodbatlas_cloud_backup_snapshot_export_job: Deprecates the `err_msg` attribute ([#2436](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2436))

BUG FIXES:

* resource/mongodbatlas_alert_configuration: Fixes an issue where the `terraform apply` command crashes if you attempt to edit an existing `mongodbatlas_alert_configuration` by adding a value to `threshold_config` ([#2463](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2463))
* resource/mongodbatlas_organization: Fixes a bug in organization resource creation where the provider crashed ([#2462](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2462))

## 1.17.4 (July 19, 2024)

ENHANCEMENTS:

* data-source/mongodbatlas_search_index: Adds attribute `stored_source` ([#2388](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2388))
* data-source/mongodbatlas_search_indexes: Adds attribute `stored_source` ([#2388](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2388))
* resource/mongodbatlas_search_index: Adds attribute `stored_source` ([#2388](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2388))

BUG FIXES:

* resource/mongodbatlas_advanced_cluster: Fixes `disk_iops` attribute for Azure cloud provider ([#2396](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2396))
* resource/mongodbatlas_cloud_backup_schedule: Updates `copy_settings` on changes (even when empty) ([#2387](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2387))
* resource/mongodbatlas_search_index: Returns error if the `analyzers` attribute contains unknown fields ([#2394](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2394))

## 1.17.3 (June 27, 2024)

## 1.17.2 (June 20, 2024)

ENHANCEMENTS:

* data-source/mongodbatlas_advanced_cluster: Adds attribute `global_cluster_self_managed_sharding` ([#2348](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2348))
* data-source/mongodbatlas_advanced_clusters: Adds attribute `global_cluster_self_managed_sharding` ([#2348](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2348))
* resource/mongodbatlas_advanced_cluster: Adds attribute `global_cluster_self_managed_sharding` ([#2348](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2348))

BUG FIXES:

* resource/mongodbatlas_project_ip_access_list: Fixes resource removal in Read() if resource doesn't exist after creation ([#2349](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2349))

## 1.17.1 (June 14, 2024)

NOTES:

* resource/mongodbatlas_federated_settings_identity_provider: OIDC Workforce and Workload are now in GA ([#2344](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2344))

## 1.17.0 (June 10, 2024)

BREAKING CHANGES:

* data-source/mongodbatlas_federated_settings_identity_provider: Replaces `audience_claim` field with `audience` ([#2310](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2310))
* data-source/mongodbatlas_federated_settings_identity_providers: Replaces `audience_claim` field with `audience` ([#2310](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2310))
* data-source/mongodbatlas_privatelink_endpoints_service_serverless: Removes `page_num` and `items_per_page` arguments ([#2336](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2336))
* resource/mongodbatlas_federated_settings_identity_provider: Replaces `audience_claim` field with `audience` ([#2310](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2310))

FEATURES:

* **New Data Source:** `mongodbatlas_control_plane_ip_addresses` ([#2331](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2331))

ENHANCEMENTS:

* data-source/mongodbatlas_federated_settings_identity_provider: Adds OIDC Workload support ([#2318](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2318))
* data-source/mongodbatlas_federated_settings_identity_provider: Adds `description` and `authorization_type` fields ([#2310](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2310))
* data-source/mongodbatlas_federated_settings_identity_providers: Adds OIDC Workload support ([#2318](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2318))
* data-source/mongodbatlas_federated_settings_identity_providers: Adds `description` and `authorization_type` fields ([#2310](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2310))
* data-source/mongodbatlas_federated_settings_identity_providers: Adds filtering support for Protocol and IdP type ([#2318](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2318))
* data-source/mongodbatlas_federated_settings_org_config: Adds `data_access_identity_provider_ids` ([#2322](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2322))
* data-source/mongodbatlas_federated_settings_org_configs: Adds `data_access_identity_provider_ids` ([#2322](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2322))
* resource/mongodbatlas_database_user: Supports Workload OIDC `mongodbatlas_database_user` ([#2323](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2323))
* resource/mongodbatlas_federated_settings_identity_provider: Adds OIDC Workload support ([#2318](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2318))
* resource/mongodbatlas_federated_settings_identity_provider: Adds `description` and `authorization_type` fields ([#2310](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2310))
* resource/mongodbatlas_federated_settings_identity_provider: Adds create and delete operations for Workforce OIDC IdP ([#2310](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2310))
* resource/mongodbatlas_federated_settings_org_config: Adds `data_access_identity_provider_ids` ([#2322](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2322))
* resource/mongodbatlas_federated_settings_org_config: Adds `user_conflicts` ([#2322](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2322))
* resource/mongodbatlas_federated_settings_org_config: Supports detaching and updating the `identity_provider_id` ([#2322](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2322))

## 1.16.2 (May 31, 2024)

BUG FIXES:

* resource/mongodbatlas_network_peering: Correctly handles GCP updates of mongodbatlas_network_peering ([#2306](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2306))
* resource/mongodbatlas_network_peering: Fixes computed values of `altas_gcp_project_id` and `atlas_vpc_name` to provide GCP related values ([#2315](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2315))

## 1.16.1 (May 28, 2024)

ENHANCEMENTS:

* data-source/mongodbatlas_cloud_backup_snapshot_export_bucket: Marks `id` as computed not required ([#2241](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2241))
* data-source/mongodbatlas_cloud_backup_snapshot_export_job: Marks `id` as computed and therefore, not required ([#2234](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2234))
* data-source/mongodbatlas_cloud_backup_snapshot_restore_job: Uses `snapshot_restore_job_id` instead of encodedID in `job_id` ([#2257](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2257))
* resource/mongodbatlas_federated_settings_org_rolemapping: Adds `role_mapping_id` computed attribute ([#2258](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2258))

BUG FIXES:

* data-source/mongodbatlas_federated_database_instance: Populates value of `data_process_region` when returned by the API ([#2223](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2223))
* data-source/mongodbatlas_federated_database_instances: Populates value of `data_process_region` when returned by the API ([#2223](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2223))
* resource/mongodbatlas_cloud_backup_schedule: Fixes behavior when resource is deleted outside of Terraform ([#2268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2268))
* resource/mongodbatlas_cloud_backup_snapshot_export_bucket Adds missing `project_id` during Read ([#2232](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2232))
* resource/mongodbatlas_cloud_backup_snapshot_export_bucket: Calls DeleteExportBucket before checking for a status update so that the delete operation doesn't hang ([#2269](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2269))
* resource/mongodbatlas_encryption_at_rest: Fixes behavior when resource is deleted outside of Terraform ([#2268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2268))
* resource/mongodbatlas_global_cluster_config: Blocks updates on global_cluster_config resource to avoid leaving the cluster in an inconsistent state ([#2282](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2282))
* resource/mongodbatlas_ldap_configuration: Disables LDAP when the resource is destroyed, instead of deleting userToDNMapping document ([#2221](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2221))
* resource/mongodbatlas_network_peering: Sets all attributes of Azure network peering as ForceNew, forcing recreation of the resource when updating ([#2299](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2299))
* resource/mongodbatlas_project: Fixes inconsistent result after apply when region_usage_restrictions are not set in configuration but returned from server ([#2291](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2291))
* resource/mongodbatlas_push_based_log_export: Fixes behavior when resource is deleted outside of Terraform ([#2268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2268))
* resource/mongodbatlas_search_deployment: Fixes behavior when resource is deleted outside of Terraform ([#2268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2268))
* resource/mongodbatlas_stream_connection: Fixes behavior when resource is deleted outside of Terraform ([#2268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2268))
* resource/mongodbatlas_stream_instance: Fixes behavior when resource is deleted outside of Terraform ([#2268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2268))

## 1.16.0 (April 29, 2024)

BREAKING CHANGES:

* data-source/mongodbatlas_database_user: Removes `password` attribute ([#2190](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2190))
* data-source/mongodbatlas_database_users: Removes `password` attribute ([#2190](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2190))
* data-source/mongodbatlas_federated_settings_identity_provider: The only allowed format for `identity_provider_id` is a 24-hexadecimal digit string ([#2185](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2185))
* data-source/mongodbatlas_organizations: Removes `include_deleted_orgs` attribute ([#2190](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2190))
* resource/mongodbatlas_federated_settings_identity_provider: Import can only use a 24-hexadecimal digit string that identifies the IdP, `idp_id`, instead of `okta_idp_id` ([#2185](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2185))
* resource/mongodbatlas_project_api_key: Removes `project_id` attribute ([#2190](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2190))

NOTES:

* data-source/mongodbatlas_federated_settings_identity_providers: Deprecates `page_num` and `items_per_page` attributes. They are not being used and will not be relevant once all results are fetched internally. ([#2207](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2207))
* data-source/mongodbatlas_teams: Deprecates data source in favour of `mongodbatlas_team` which has the same implementation. This aligns the name of the resource with the implementation which fetches a single team. ([#2208](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2208))
* data-source/mongodbatlas_third_party_integration: Deprecates `scheme` attribute as it is not being used ([#2216](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2216))
* data-source/mongodbatlas_third_party_integrations: Deprecates `scheme` attribute as it is not being used ([#2216](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2216))
* provider: New changelog format has been incorporated following [Terraform Changelog Specification](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning#changelog-specification) ([#2124](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2124))
* resource/mongodbatlas_teams: Deprecates resource in favour of `mongodbatlas_team` which has the same implementation. This aligns the name of the resource with the implementation which manages a single team. ([#2208](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2208))
* resource/mongodbatlas_third_party_integration: Deprecates `scheme` attribute as it is not being used ([#2216](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2216))

FEATURES:

* **New Data Source:** `data-source/mongodbatlas_push_based_log_export` ([#2169](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2169))
* **New Resource:** `resource/mongodbatlas_push_based_log_export` ([#2169](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2169))

ENHANCEMENTS:

* data-source/mongodbatlas_alert_configuration: Adds `integration_id` attribute ([#2212](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2212))
* data-source/mongodbatlas_alert_configurations: Adds `integration_id` attribute ([#2212](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2212))
* data-source/mongodbatlas_backup_compliance_policy: Adds `policy_item_yearly` attribute ([#2109](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2109))
* data-source/mongodbatlas_cloud_backup_schedule: Adds `policy_item_yearly` attribute ([#2109](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2109))
* data-source/mongodbatlas_project: Adds `tags` attribute ([#2135](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2135))
* data-source/mongodbatlas_projects: Adds `tags` attribute ([#2135](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2135))
* data-source/mongodbatlas_serverless_instance: Adds `auto_indexing` attribute ([#2100](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2100))
* data-source/mongodbatlas_stream_connection: Reaches GA (General Availability) ([#2209](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2209))
* data-source/mongodbatlas_stream_connections: Reaches GA (General Availability) ([#2209](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2209))
* data-source/mongodbatlas_stream_instance: Reaches GA (General Availability) ([#2209](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2209))
* data-source/mongodbatlas_stream_instances: Reaches GA (General Availability) ([#2209](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2209))
* data-source/mongodbatlas_third_party_integration: New `id` value which can be used for referencing a third party integration ([#2217](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2217))
* data-source/mongodbatlas_third_party_integrations: New `id` value which can be used for referencing a third party integration ([#2217](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2217))
* resource/mongodbatlas_alert_configuration: Adds `integration_id` attribute ([#2212](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2212))
* resource/mongodbatlas_backup_compliance_policy: Adds `policy_item_yearly` attribute ([#2109](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2109))
* resource/mongodbatlas_cloud_backup_schedule: Adds `policy_item_yearly` attribute ([#2109](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2109))
* resource/mongodbatlas_privatelink_endpoint_service_serverless: Adds support for updating `comment` attribute in-place. ([#2133](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2133))
* resource/mongodbatlas_project: Adds `tags` attribute ([#2135](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2135))
* resource/mongodbatlas_serverless_instance: Adds `auto_indexing` attribute ([#2100](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2100))
* resource/mongodbatlas_stream_connection: Reaches GA (General Availability) ([#2209](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2209))
* resource/mongodbatlas_stream_instance: Reaches GA (General Availability) ([#2209](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2209))
* resource/mongodbatlas_third_party_integration: New `id` value which can be used for referencing a third party integration ([#2217](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2217))

BUG FIXES:

* data-source/mongodbatlas_advanced_cluster: Converts `replication_specs` from TypeSet to TypeList. This fixes an issue where some items were not returned in the results. ([#2145](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2145))
* data-source/mongodbatlas_advanced_clusters: Converts `replication_specs` from TypeSet to TypeList. This fixes an issue where some items were not returned in the results. ([#2145](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2145))
* data-source/mongodbatlas_network_peering: Ensures `accepter_region_name` is set when it is has the same value as the container resource ([#2105](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2105))
* resource/mongodbatlas_cluster: Fixes nil pointer dereference if `advanced_configuration` update fails ([#2139](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2139))
* resource/mongodbatlas_maintenance_window: Fixes `day_of_week` param as **required** when calling the API ([#2163](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2163))
* resource/mongodbatlas_privatelink_endpoint_serverless: Removes setting default comment during create. ([#2133](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2133))
* resource/mongodbatlas_project: Reads `region_usage_restrictions` attribute value from get request ([#2104](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2104))

## [v1.15.3](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.15.3) (2024-03-27)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.15.2...v1.15.3)

**Bug Fixes**

- fix: Fixes `network_container` resource update [\#2055](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2055) ([lantoli](https://github.com/lantoli))
- fix: Uses `overwriteBackupPolicies` in `mongodbatlas_backup_compliance_policy` to avoid overwriting non complying backup policies in updates [\#2054](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2054) ([maastha](https://github.com/maastha))

**Internal Improvements**

- chore: Allows user to specify to use an existing tag for release [\#2053](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2053) ([maastha](https://github.com/maastha))
- chore: Fixes Slack notification button to GH action run text [\#2093](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2093) ([oarbusi](https://github.com/oarbusi))
- doc: Fixes import command in `mongodbatlas_third_party_integration` doc [\#2083](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2083) ([zach-carr](https://github.com/zach-carr))
- chore: Reuses project in tests - `mongodbatlas_auditing` [\#2082](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2082) ([lantoli](https://github.com/lantoli))
- test: Converting a test case to a migration test [\#2081](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2081) ([EspenAlbert](https://github.com/EspenAlbert))
- doc: Specifies that upgrades From Replica Sets to Multi-Sharded Instances of `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` might lead to error [\#2080](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2080) ([oarbusi](https://github.com/oarbusi))
- chore: Reuses project in tests - `mongodbatlas_project` [\#2078](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2078) ([lantoli](https://github.com/lantoli))
- doc: Adds an example using `ignore_changes` when `autoscaling` is enabled [\#2077](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2077) ([zach-carr](https://github.com/zach-carr))
- chore: Uses mocks for unit tests in Atlas Go SDK [\#2075](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2075) ([lantoli](https://github.com/lantoli))
- chore: Updates Atlas Go SDK [\#2074](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2074) ([github-actions[bot]](https://github.com/apps/github-actions))
- doc: Improve Readme Requirements to point to section in HashiCorp Registry Docs  [\#2073](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2073) ([Zuhairahmed](https://github.com/Zuhairahmed))
- doc: Updates bug report to include Terraform version support guidance [\#2072](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2072) ([maastha](https://github.com/maastha))
- chore: Fixes send notification when test suite fails [\#2071](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2071) ([oarbusi](https://github.com/oarbusi))
- chore: Fixes federated test [\#2070](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2070) ([lantoli](https://github.com/lantoli))
- chore: Follow up to use global mig project in tests [\#2068](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2068) ([lantoli](https://github.com/lantoli))
- chore: Reuses project in tests - `mongodbatlas_project_ip_access_list` [\#2067](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2067) ([lantoli](https://github.com/lantoli))
- chore: Adds mig tests and refactor - `mongodbatlas_search_index` [\#2065](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2065) ([lantoli](https://github.com/lantoli))
- chore: Corrects order of checks in `data_source_federated_settings_identity_providers_test` [\#2064](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2064) ([oarbusi](https://github.com/oarbusi))
- chore: Removes old service from mockery [\#2063](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2063) ([EspenAlbert](https://github.com/EspenAlbert))
- chore: Enables Github action linter and removes set terminal in release action [\#2062](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2062) ([lantoli](https://github.com/lantoli))
- chore: Allows `MONGODB_ATLAS_PROJECT_ID` for local executions [\#2060](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2060) ([lantoli](https://github.com/lantoli))
- doc: Clarifies private endpoint resource docs [\#2059](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2059) ([zach-carr](https://github.com/zach-carr))
- chore: Automates changing Terraform supported versions in provider documentation [\#2058](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2058) ([maastha](https://github.com/maastha))
- test: Enables simulation of cloud-dev using hoverfly in alert configuration acceptance tests [\#2057](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2057) ([AgustinBettati](https://github.com/AgustinBettati))
- refactor: Uses mocks on `admin.APIClient` instead of custom `ClusterService` [\#2056](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2056) ([EspenAlbert](https://github.com/EspenAlbert))
- chore: Send Slack message for Terraform Compatibility Matrix is executed [\#2052](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2052) ([oarbusi](https://github.com/oarbusi))
- chore: Follow-up to Atlas SDK upgrade [\#2051](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2051) ([lantoli](https://github.com/lantoli))
- chore: Fixes test names for `mongodbatlas_network_container` resource [\#2046](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2046) ([lantoli](https://github.com/lantoli))
- refactor: Uses mocks on admin.APIClient for Project+Teams+ClustersAPIs instead of custom service [\#2045](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2045) ([EspenAlbert](https://github.com/EspenAlbert))
- chore: Updates Atlas Go SDK [\#2044](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2044) ([github-actions[bot]](https://github.com/apps/github-actions))
- test: Mocks the EncryptionAtRestUsingCustomerKeyManagementApi instead of using custom service [\#2043](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2043) ([EspenAlbert](https://github.com/EspenAlbert))
- chore: Reuses project in tests for `mongodbatlas_advanced_cluster` resource [\#2042](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2042) ([lantoli](https://github.com/lantoli))
- doc: Updates online archive resource docs [\#2041](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2041) ([zach-carr](https://github.com/zach-carr))
- chore: Reuses project in tests for `mongodbatlas_cluster_outage_simulation` resource [\#2040](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2040) ([lantoli](https://github.com/lantoli))
- chore: Reuses project in tests for `mongodbatlas_network_container` resource [\#2039](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2039) ([lantoli](https://github.com/lantoli))
- chore: Reuses project in tests for `mongodbatlas_x509_authentication_database_user` resource [\#2038](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2038) ([lantoli](https://github.com/lantoli))
- chore: Reuses project in tests for `mongodbatlas_project_api_key` resource [\#2037](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2037) ([lantoli](https://github.com/lantoli))
- chore: Reuses project in tests for `mongodbatlas_cluster` resource [\#2036](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2036) ([lantoli](https://github.com/lantoli))
- chore: Bump tj-actions/verify-changed-files from 843c0b95f87cd81a2efe729380c6d1f11fb3ea12 to 1e517a7f5663673148ceb7c09c1900e5af48e7a1 [\#2092](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2092) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/hcl/v2 from 2.20.0 to 2.20.1 [\#2091](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2091) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/zclconf/go-cty from 1.14.3 to 1.14.4 [\#2089](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2089) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.51.3 to 1.51.8 [\#2088](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2088) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/mongodb-forks/digest from 1.0.5 to 1.1.0 [\#2087](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2087) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.50.32 to 1.51.3 [\#2049](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2049) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/stretchr/testify from 1.8.4 to 1.9.0 [\#2048](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2048) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump actions/checkout from 4.1.1 to 4.1.2 [\#2047](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2047) ([dependabot[bot]](https://github.com/apps/dependabot))

**Closed Issues**

- \[Bug\]: Removing `mongodbatlas_cloud_backup_schedule` resource is resulting in error `Continuous Cloud Backup cannot be on without an hourly policy item` [\#2029](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/2029)

## [v1.15.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.15.2) (2024-03-15)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.15.1...v1.15.2)

**Enhancements**

- feat: Adds support for Sample stream type to `mongodbatlas_stream_connection` [\#2026](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2026) ([maastha](https://github.com/maastha))
- feat: Adds support for using DEV/QA for mongodbgov [\#2009](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2009) ([maastha](https://github.com/maastha))
- refactor: Renames `MONGODB_ATLAS_ENABLE_BETA` to `MONGODB_ATLAS_ENABLE_PREVIEW` for features in preview [\#2004](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2004) ([EspenAlbert](https://github.com/EspenAlbert))
- feat: Adds StreamConfig attribute to `mongodbatlas_stream_instance` resource and datasources [\#1989](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1989) ([oarbusi](https://github.com/oarbusi))
- feat: Adds support for `region` & `customer_endpoint_dns_name` in privatelink [\#1982](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1982) ([EspenAlbert](https://github.com/EspenAlbert))
- feat: Updates `mongodbatlas_stream_connection` resource & data sources to support `dbRoleToExecute` attribute [\#1980](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1980) ([maastha](https://github.com/maastha))

**Bug Fixes**

- fix: Avoids sending database user password in update request if the value has not changed [\#2005](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2005) ([AgustinBettati](https://github.com/AgustinBettati))
- fix: Removes escape logic on IP address in mongodbatlas\_access\_list\_api\_key [\#1998](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1998) ([maastha](https://github.com/maastha))
- fix: Enables creation of database event trigger with watch against database and not collection [\#1968](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1968) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Updates Atlas Go SDK and unifies resource prefix names in tests [\#1966](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1966) ([github-actions[bot]](https://github.com/apps/github-actions))
- fix: Removes `SchemaConfigModeAttr` from resources [\#1961](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1961) ([maastha](https://github.com/maastha))
- fix: Fixes timeout and removes deletion logic on update failure for `mongodbatlas_search_index` resource [\#1950](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1950) ([maastha](https://github.com/maastha))
- chore: Fixes `search_deployment` template to fix doc structure [\#1943](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1943) ([oarbusi](https://github.com/oarbusi))
- fix: Converts `snapshot_id` field from required to optional [\#1924](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1924) ([marcosuma](https://github.com/marcosuma))

**Closed Issues**

- \[Bug\]: Many provider crashes when running from GitLab CI pipeline [\#1944](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1944)
- Terraform lifecycle `ignore_changes` tags [\#2006](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/2006)
- \[Bug\]: `mongodbatlas_access_list_api_key` CIDRs/subnets \(not single IPs\) fresh create after upgrade from 1.14 to 1.15 [\#1984](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1984)
- \[Feature\]: `mongodbatlas_event_trigger` does not support Database operations [\#1967](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1967)
- \[Bug\] `mongodbatlas_custom_db_role` created with the wrong permissions [\#1963](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1963)
- \[Bug\]:  `mongodbatlas_privatelink_endpoint_service` for GCP - Provider produced inconsistent final plan [\#1957](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1957)
- \[Bug\]: Impossible to create a database scope \(database watch against\) in `mongodbatlas_event_trigger` [\#1956](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1956)
- \[Feature\]: Add resources to create function [\#1954](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1954)
- What will happen to the password field when DB user is imported? [\#1952](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1952)
- \[Feature\]: Add resources to automatically create `mongodbatlas_event_trigger` resource [\#1949](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1949)
- \[Bug\]: missing data source for `app_id` & `service_id` [\#1942](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1942)
- \[Feature\]: unable to setup log forwarding to S3 [\#1933](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1933)
- \[Bug\]: Unable to create billing alert configuration. [\#1927](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1927)
- \[Bug\]: Can't setup privatelink\_endpoint\_service\_data\_federation\_online\_archive region or VPC Endpoint DNS Name for AWS [\#1878](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1878)
- \[Bug\]: Undocumented behaviour privatelink\_endpoint / circle dependency [\#1872](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1872)
- \[Bug\]: Provider produces inconsistent result after importing encryption\_at\_rest [\#1805](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1805)
- \[Bug\]: backup\_compliance\_policy resource missing required attribute [\#1800](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1800)
- `mongodbatlas_cloud_backup_snapshot_export_bucket` resource stuck on `Still distroying...` [\#1569](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1569)
- Assigning a custom role to a new user fails \(400\) UNSUPPORTED\_ROLE [\#1522](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1522)

**Internal Improvements**

- chore: Reverts actionlint to enable release gh action [\#2034](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2034) ([lantoli](https://github.com/lantoli))
- chore: Specifies shell in release action [\#2033](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2033) ([lantoli](https://github.com/lantoli))
- chore: Fixes user terminal in release Github action [\#2032](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2032) ([lantoli](https://github.com/lantoli))
- refactor: Modifies mocking of search deployment unit test directly to SDK removing intermediate service [\#2028](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2028) ([AgustinBettati](https://github.com/AgustinBettati))
- doc: Fixes doc header for datasource `mongodbatlas_cloud_provider_access_setup` [\#2027](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2027) ([jvincent-mongodb](https://github.com/jvincent-mongodb))
- chore: Removes usage of vars.MONGODB\_ATLAS\_ENABLE\_PREVIEW [\#2024](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2024) ([EspenAlbert](https://github.com/EspenAlbert))
- chore: Run Go tests using specific packages [\#2014](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2014) ([lantoli](https://github.com/lantoli))
- chore: Creates project for test execution [\#2010](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2010) ([lantoli](https://github.com/lantoli))
- chore: Adjusts go.mod files after removal of integration tests [\#2008](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2008) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Reuses projects in executions in some resources and rename mig tests [\#2007](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2007) ([lantoli](https://github.com/lantoli))
- chore: Simplifies makefile [\#2003](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2003) ([lantoli](https://github.com/lantoli))
- doc: Adds export example to mongodbatlas\_cloud\_backup\_snapshot\_export\_job resource doc [\#2002](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2002) ([zach-carr](https://github.com/zach-carr))
- chore: Enables tests in CI for `mongodbatlas_privatelink_endpoint_service_data_federation_online_archive` resource [\#2001](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2001) ([lantoli](https://github.com/lantoli))
- chore: Deletes outdated integration tests [\#1999](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1999) ([lantoli](https://github.com/lantoli))
- chore: Unifies pass of org variables in tests [\#1997](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1997) ([lantoli](https://github.com/lantoli))
- chore: Stops using Network org in tests [\#1996](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1996) ([lantoli](https://github.com/lantoli))
- chore: Changes skip project logic for cleanup action [\#1995](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1995) ([lantoli](https://github.com/lantoli))
- chore: Adds Github Actions linter [\#1988](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1988) ([lantoli](https://github.com/lantoli))
- doc: Adds additional PR title guidelines [\#1986](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1986) ([zach-carr](https://github.com/zach-carr))
- doc: Fixes mongodbatlas\_roles\_org\_id data source documentation [\#1985](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1985) ([maastha](https://github.com/maastha))
- doc: Correct reference to cluster attribute `auto_scaling_disk_gb_enabled` [\#1983](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1983) ([SimonPerryman](https://github.com/SimonPerryman))
- chore: Updates Atlas Go SDK, default MongoDB version updated from 6.0 to 7.0 [\#1981](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1981) ([github-actions[bot]](https://github.com/apps/github-actions))
- chore: Updates Go to 1.22 and Terraform to 1.7.4 [\#1979](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1979) ([lantoli](https://github.com/lantoli))
- doc: Updates third\_party\_integration documentation. [\#1973](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1973) ([marcosuma](https://github.com/marcosuma))
- refactor: Remove redundant parameter in checkExists for trigger tests [\#1972](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1972) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Upgrades `ldap_configuration` and `ldap_verify` resources to auto-generated SDK [\#1971](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1971) ([lantoli](https://github.com/lantoli))
- chore: Fixes some acceptance tests [\#1970](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1970) ([lantoli](https://github.com/lantoli))
- doc: HashiCorp Terraform Version Compatibility Matrix [\#1969](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1969) ([Zuhairahmed](https://github.com/Zuhairahmed))
- fix: Uses `google_cloud_kms_config` correctly in `mongodbatlas_encryption_at_rest` creation [\#1962](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1962) ([DasLampe](https://github.com/DasLampe))
- chore: Signs created tag during release process [\#1960](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1960) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Avoids skipping stream resources in qa tests [\#1959](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1959) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Restores atlas streams guide in examples section [\#1958](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1958) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Drop dependency on pointy [\#1953](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1953) ([gssbzn](https://github.com/gssbzn))
- chore: Upgrades `advanced_cluster` resource to auto-generated SDK [\#1947](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1947) ([lantoli](https://github.com/lantoli))
- chore: Upgrades `cloud_backup_schedule` resource to auto-generated SDK [\#1946](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1946) ([oarbusi](https://github.com/oarbusi))
- doc: Add doc-preview to the CONTRIBUTING.md [\#1945](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1945) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Upgrades `global_cluster_config` resource to auto-generated SDK [\#1938](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1938) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades cloud\_backup\_snapshot resource to auto-generated SDK [\#1936](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1936) ([lantoli](https://github.com/lantoli))
- chore: Upgrades `cluster_outage_simulation` resource to auto-generated SDK [\#1935](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1935) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades `privatelink_endpoint_service_serverless` resource to auto-generated SDK [\#1932](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1932) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades custom\_dns\_configuration\_cluster\_aws resource to auto-generated SDK [\#1930](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1930) ([lantoli](https://github.com/lantoli))
- chore: Disables tests for `privatelink_endpoint_service_data_federation_online_archive` in CI [\#1929](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1929) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades backup\_compliance\_policy resource to auto-generated SDK [\#1928](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1928) ([lantoli](https://github.com/lantoli))
- chore: Fixes example on privatelink\_endpoint\_service. [\#1926](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1926) ([marcosuma](https://github.com/marcosuma))
- chore: Updates Atlas Go SDK [\#1925](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1925) ([github-actions[bot]](https://github.com/apps/github-actions))
- chore: Upgrades `network_container` resource to auto-generated SDK [\#1920](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1920) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades serverless\_instance resource to auto-generated SDK [\#1913](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1913) ([lantoli](https://github.com/lantoli))
- chore: Upgrades datalake\_pipeline resource to auto-generated SDK [\#1911](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1911) ([lantoli](https://github.com/lantoli))
- chore: Upgrades maintenance\_window resource to auto-generated SDK [\#1886](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1886) ([lantoli](https://github.com/lantoli))
- chore: Bump peter-evans/create-pull-request from 6.0.1 to 6.0.2 [\#2022](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2022) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump actions/checkout from 4.1.1 to 4.1.2 [\#2021](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2021) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump bewuethr/shellcheck-action from 2.1.2 to 2.2.0 [\#2020](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2020) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 79a8cef5d9ef3ab541ee07ef179bd6c3c2d42ecc to 843c0b95f87cd81a2efe729380c6d1f11fb3ea12 [\#2019](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2019) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/hcl/v2 from 2.19.1 to 2.20.0 [\#2018](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2018) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-go from 0.21.0 to 0.22.1 [\#2017](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2017) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/zclconf/go-cty from 1.14.2 to 1.14.3 [\#2016](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/2016) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump dorny/paths-filter from 3.0.1 to 3.0.2 [\#1994](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1994) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from cb704d8c717959cf92ec9db9131372bc20356aa2 to 79a8cef5d9ef3ab541ee07ef179bd6c3c2d42ecc [\#1993](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1993) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump peter-evans/create-pull-request from 6.0.0 to 6.0.1 [\#1992](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1992) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-testing from 1.6.0 to 1.7.0, terraform-plugin-framework from 1.5.0 to 1.6.1 and terraform-plugin-mux from 0.14.0 to 0.15.0 [\#1991](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1991) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.50.27 to 1.50.32 [\#1990](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1990) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.50.22 to 1.50.27 [\#1977](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1977) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.50.7 to 1.50.12 [\#1918](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1918) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump peter-evans/create-pull-request from 5.0.2 to 6.0.0 [\#1915](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1915) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.50.17 to 1.50.22 [\#1965](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1965) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump dorny/paths-filter from 3.0.0 to 3.0.1 [\#1964](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1964) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump golangci/golangci-lint-action from 3.7.0 to 4.0.0 [\#1941](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1941) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 8b28bea118e7723e4672bc7ac323bcd26f271ec4 to cb704d8c717959cf92ec9db9131372bc20356aa2 [\#1940](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1940) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.50.12 to 1.50.17 [\#1939](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1939) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.15.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.15.1) (2024-02-07)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.15.0...v1.15.1)

**Bug Fixes**

- fix: Sets `replication_specs` IDs when updating them. [\#1876](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1876) ([marcosuma](https://github.com/marcosuma))

**Internal Improvements**


- chore: Upgrades `privatelink_endpoint_service_data_federation_online_archive` resource to auto-generated SDK [\#1910](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1910) ([oarbusi](https://github.com/oarbusi))
- chore: Fixes test for `federated_settings_identity_provider` in QA environment [\#1912](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1912) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades `privatelink_endpoint_serverless` resource to auto-generated SDK [\#1908](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1908) ([oarbusi](https://github.com/oarbusi))
- chore: Fixes acceptance and migrations tests not running in CI [\#1907](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1907) ([lantoli](https://github.com/lantoli))
- chore: Upgrades `roles_org_id` resource to auto-generated SDK [\#1906](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1906) ([lantoli](https://github.com/lantoli))
- chore: Upgrades `teams` resource to auto-generated SDK [\#1905](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1905) ([oarbusi](https://github.com/oarbusi))
- doc: Fixes `mongodbatlas_privatelink_endpoint_service_data_federation_online_archives` doc [\#1903](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1903) ([nsmith78660](https://github.com/nsmith78660))
- doc: Fixes some of the typos within the `README.MD` for the PIT example [\#1902](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1902) ([nsmith78660](https://github.com/nsmith78660))
- chore: Upgrades `private_link_endpoint` resource to auto-generated SDK. [\#1901](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1901) ([marcosuma](https://github.com/marcosuma))
- test: Enables Acceptance test in CI for `mongodbatlas_federated_settings_identity_provider` [\#1895](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1895) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades `x509authentication_database_user` resource to auto-generated SDK [\#1884](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1884) ([lantoli](https://github.com/lantoli))
- chore: Bump marocchino/sticky-pull-request-comment from 2.8.0 to 2.9.0 [\#1916](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1916) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files [\#1914](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1914) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.15.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.15.0) (2024-02-01)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.14.0...v1.15.0)

**Breaking changes:**

- remove!: Removes `page_num` and `items_per_page` attributes in `mongodbatlas_search_indexes` data source [\#1880](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1880) ([AgustinBettati](https://github.com/AgustinBettati))
- remove!: Removes `cloud_provider_access`  resource and data source [\#1804](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1804) ([andreaangiolillo](https://github.com/andreaangiolillo))

**Enhancements**

- feat: Adds support to new Federated Auth parameters for OIDC  [\#1874](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1874) ([oarbusi](https://github.com/oarbusi))
- feat: Adds new `ip_addresses` computed attribute in `mongodbatlas_project` resource and data sources [\#1850](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1850) ([AgustinBettati](https://github.com/AgustinBettati))
- feat: Adds `mongodbatlas_organization` New Parameters Support  [\#1835](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1835) ([maastha](https://github.com/maastha))

**Bug Fixes**

- `mongodbatlas_project_ip_access_list` Unexpected replacement of CIDR with IP address [\#1571](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1571)
- fix: Adds toUpperCase to provider and region fields in cluster resources [\#1837](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1837) ([andreaangiolillo](https://github.com/andreaangiolillo))
- fix: Improves error message when improperly setting `provider_region_name` field [\#1815](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1815) ([maastha](https://github.com/maastha))
- fix: Provider produces inconsistent result after importing `mongodbatlas_encryption_at_rest` [\#1813](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1813) ([andreaangiolillo](https://github.com/andreaangiolillo))
- fix: Update attributes `copy_protection_enabled`, `pit_enabled` and `encryption_at_rest_enabled` in the  resource `mongodbatlas_backup_compliance_policy` to be Optional [\#1803](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1803) ([andreaangiolillo](https://github.com/andreaangiolillo))
- fix: Incompatible schema defined for `mongodbatlas_backup_compliance_policy` [\#1799](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1799) ([andreaangiolillo](https://github.com/andreaangiolillo))
- fix: Fixes `mongodbatlas_clusters` plural data source to set `auto_scaling_disk_gb_enabled` attribute correctly [\#1722](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1722) ([maastha](https://github.com/maastha))

**Closed Issues**

- CANNOT_DISABLE_PIT_WITH_BACKUP_COMPLIANCE_POLICY  [\#1855](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1855)
- Error creating user `mongodbatlas_database_user` [\#1852](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1852)
- Serverless `aws_vpc_endpoint` creation fails [\#1826](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1826)
- [Bug]: Changing Scope block in the databaseuser resource results in a replacement [\#1821](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1821)
- [Bug]: `oplog_min_retention_hours` is not expected here [\#1818](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1818)
- Unable to add encryption at rest through terraform [\#1766](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1766)
- Provider produced inconsistent result after apply [\#1708](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1708)
- Error: json: cannot unmarshal number 9501614080 into Go struct field CloudProviderSnapshot.storageSizeBytes of type int  [\#1341](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1341)
- [Query] Unable to find the `endpoint_service_id` [\#1281](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1281)

**Internal Improvements**

- chore: Upgrades `mongodbatlas_project_invitation` resource to auto-generated SDK [\#1900](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1900)
- chore: Upgrades `mongodbatlas_org_invitation` resource to auto-generated SDK [\#1897](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1897)
- chore: Changes from env variable to input in Import GPG key [\#1898](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1898) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades `mongodbatlas_custom_db_role` resource to auto-generated SDK [\#1896](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1896) ([lantoli](https://github.com/lantoli))
- chore: New released Atlas Go SDK can't be used until several hours later [\#1885](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1885) ([gssbzn](https://github.com/gssbzn))
- doc: Change documentation for new attributes to support OIDC Identity providers [\#1883](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1883) ([oarbusi](https://github.com/oarbusi))
- chore: Upgrades auditing resource to auto-generated SDK [\#1881](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1881) ([lantoli](https://github.com/lantoli))
- chore: Upgrades `mongodbatlas_api_key` resource to auto-generated SDK [\#1879](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1879) ([lantoli](https://github.com/lantoli))
- chore: Upgrades  `mongodbatlas_access_list_api_key` resource to auto-generated SDK [\#1877](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1877) ([lantoli](https://github.com/lantoli))
- chore: Updates `mongodbatlas_online_archive` resource with new SDK [\#1875](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1875) ([lantoli](https://github.com/lantoli))
- chore: Removes scheduled tests for Terraform 1.6.x [\#1871](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1871) ([lantoli](https://github.com/lantoli))
- chore: Updates Atlas Go SDK [\#1865](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1865) ([github-actions[bot]](https://github.com/apps/github-actions))
- chore: Unifies SDK connection getting in tests [\#1864](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1864) ([lantoli](https://github.com/lantoli))
- doc: Updates documentation for `mongodbatlas_database_user` resource to guide users to set database name `admin` for custom roles  [\#1862](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1862) ([maastha](https://github.com/maastha))
- doc: Remove Extra Bracket from API Key Example [\#1861](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1861) ([Zuhairahmed](https://github.com/Zuhairahmed))
- chore: Upgrades go toolchain to 1.21.6, Terraform latest version to 1.7.x [\#1860](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1860) ([lantoli](https://github.com/lantoli))
- chore: Increase parallelism in Test Suite [\#1859](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1859) ([lantoli](https://github.com/lantoli))
- chore: Ping GH actions to a GitSHA [\#1858](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1858) ([andreaangiolillo](https://github.com/andreaangiolillo))
- test: Adds unit tests for `mongodbatlas_federated_settings_identity_provider` [\#1857](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1857) ([oarbusi](https://github.com/oarbusi))
- chore: Updates Atlas Go SDK [\#1856](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1856) ([github-actions[bot]](https://github.com/apps/github-actions))
- chore: Updates `mongodbatlas_federated_database_instance` resource with new SDK [\#1854](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1854) ([lantoli](https://github.com/lantoli))
- chore: Updates `mongodbatlas_search_deployment` resource with new SDK [\#1853](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1853) ([lantoli](https://github.com/lantoli))
- chore: Updates `mongodbatlas_alert_configuration` resource with new SDK [\#1851](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1851) ([lantoli](https://github.com/lantoli))
- doc: Update templates/README.md [\#1849](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1849) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Updates `mongodbatlas_search_index` resource with new SDK [\#1848](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1848) ([lantoli](https://github.com/lantoli))
- chore: Define default resource and data source documentation template [\#1847](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1847) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Updates project resource with new SDK [\#1843](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1843) ([lantoli](https://github.com/lantoli))
- chore: Use new attributes from new SDK in Federated Auth [\#1842](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1842) ([oarbusi](https://github.com/oarbusi))
- chore: Updates `mongodbatlas_database_user` resource with new SDK [\#1840](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1840) ([lantoli](https://github.com/lantoli))
- chore: Adjusts scaffold command to generate config file for schema scaffolding [\#1839](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1839) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Updates the Advanced Cluster and Federation tests to use EU\_WEST\_1 and EU\_WEST\_2 regions [\#1838](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1838) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Addressed followup comments in #1833 [\#1836](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1836) ([andreaangiolillo](https://github.com/andreaangiolillo))
- refactor: Migrates `mongodbatlas_organization` resource & data sources to new Atlas SDK [\#1834](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1834) ([maastha](https://github.com/maastha))
- doc: Adjust documentation of `mongodbatlas_search_deployment` resource and data source [\#1833](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1833) ([andreaangiolillo](https://github.com/andreaangiolillo))
- doc: Update CONTRIBUTING.md to explain how to generate doc [\#1832](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1832) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Uses new Atlas Go SDK version [\#1831](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1831) ([lantoli](https://github.com/lantoli))
- refactor: Adjusts `mongodbatlas_search_deployment` schema definitions with structure of new schema scaffolding command [\#1830](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1830) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Define a new makefile command for running tfplugindocs [\#1829](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1829) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Defines new schema scaffolding command using terraform code generation tools [\#1827](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1827) ([AgustinBettati](https://github.com/AgustinBettati))
- doc: Project and Org API Key Support Clarification  [\#1822](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1822) ([Zuhairahmed](https://github.com/Zuhairahmed))
- chore: Runs nightly tests for Terraform CLI version 1.2.x  [\#1820](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1820) ([lantoli](https://github.com/lantoli))
- chore: Rename folders within `examples/` to match the resource names used in each example [\#1819](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1819) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Add templates for `mongodbatlas_search_deployment` [\#1816](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1816) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Adjust name of acceptance test runner workflow [\#1814](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1814) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Adjusts release process to run tests against QA before releasing [\#1812](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1812) ([AgustinBettati](https://github.com/AgustinBettati))
- doc: DOCSP-34594 clarify use case for deleting network containers [\#1811](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1811) ([kanchana-mongodb](https://github.com/kanchana-mongodb))
- chore: Skips Stream test group when running in QA [\#1810](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1810) ([AgustinBettati](https://github.com/AgustinBettati))
- test: Add unit tests to `mongodbatlas_advanced_cluster` [\#1809](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1809) ([oarbusi](https://github.com/oarbusi))
- chore: Migrates `mongodbatlas_federated_settings_identity_provider` to new auto-generated SDK [\#1808](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1808) ([oarbusi](https://github.com/oarbusi))
- doc: Fixes documentation of required arguments in api key data sources [\#1807](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1807) ([AgustinBettati](https://github.com/AgustinBettati))
- doc: Clarifies lack of support for import statement in organization resource [\#1806](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1806) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Add tfplugindocs & tfplugingen-openapi to the dev tools [\#1798](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1798) ([andreaangiolillo](https://github.com/andreaangiolillo))
- test: Fixes flaky behaviour of atlas user data source test [\#1797](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1797) ([AgustinBettati](https://github.com/AgustinBettati))
- doc: Clean up doc bug of resource `mongodbatlas_access_list_api_key` resource name [\#1796](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1796) ([Zuhairahmed](https://github.com/Zuhairahmed))
- chore: Replace PlanOnly tests with PreApply `plancheck.ExpectEmptyPlan()` as advised in terraform-plugin-testing 1.6.0 [\#1795](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1795) ([AgustinBettati](https://github.com/AgustinBettati))
- refactor: Extract acceptance test into separate workflow to handle configuration of env variables [\#1794](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1794) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Enable Network peering acc tests using AWS to run as part of CI [\#1793](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1793) ([andreaangiolillo](https://github.com/andreaangiolillo))
- test: Fixes tests after updating terraform-plugin-testing from 1.5.1 to 1.6.0 [\#1792](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1792) ([AgustinBettati](https://github.com/AgustinBettati))
- doc: Updated Guidance for when users are unable to delete cluster and backup schedule with backup compliance policy enabled [\#1790](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1790) ([Zuhairahmed](https://github.com/Zuhairahmed))
- chore: Allows to to run acceptance tests for a published Atlas provider version [\#1789](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1789) ([lantoli](https://github.com/lantoli))
- chore: Uses Mockery instead of manually created mocks in project unit tests [\#1788](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1788) ([oarbusi](https://github.com/oarbusi))
- test: Adds migration tests for stream resources [\#1787](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1787) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Use Mockery instead of manually created mocks for `mongodbatlas_encryption_at_rest` unit tests [\#1786](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1786) ([oarbusi](https://github.com/oarbusi))
- chore: Removes coverage reports [\#1785](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1785) ([lantoli](https://github.com/lantoli))
- chore: Run external depended tests in CI `mongodbatlas_project_ip_access_list` [\#1784](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1784) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: unit test failure missing permission [\#1778](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1778) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Uses only latest provider version for Test Suite workflow [\#1777](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1777) ([lantoli](https://github.com/lantoli))
- chore: update `bug_report.yml` [\#1776](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1776) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: CLOUDP-214958: Run acc tests faster in local reusing cluster database\_user [\#1775](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1775) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: PR title check ensures uppercase is used in first character [\#1774](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1774) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Uses ExternalProviders helper functions for tests [\#1773](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1773) ([lantoli](https://github.com/lantoli))
- chore: Enable running acceptance tests against QA environment [\#1772](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1772) ([AgustinBettati](https://github.com/AgustinBettati))
- doc: added `CODE_OF_CONDUCT.md` [\#1771](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1771) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Runs tests for latest minor versions of each major TF version [\#1769](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1769) ([lantoli](https://github.com/lantoli))
- chore: CLOUDP-215956 - Convert our bug report issue on GitHub to forms [\#1767](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1767) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Creates mock generation infrastructure [\#1763](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1763) ([lantoli](https://github.com/lantoli))
- chore: update the alert notifications regions for Datadog [\#1761](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1761) ([jfmainville](https://github.com/jfmainville))
- chore: Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.31.0 to 2.32.0 [\#1893](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1893) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.50.3 to 1.50.7 [\#1892](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1892) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-mux from 0.13.0 to 0.14.0 [\#1891](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1891) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-go from 0.20.0 to 0.21.0 [\#1890](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1890) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump peter-evans/create-or-update-comment from 3.1.0 to 4.0.0 [\#1889](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1889) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 58f5ac78e19e6cc3fb9d4048ae1a13bf364fa983 to 5ef175f2fd84957530d0fdd1384a541069e403f2 [\#1888](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1888) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump dorny/paths-filter from 2.11.1 to 3.0.0 [\#1887](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1887) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.49.22 to 1.50.3 [\#1873](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1873) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump crazy-max/ghaction-import-gpg from 2.1.0 to 6.1.0 [\#1870](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1870) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from d9a97a5b5231f455f0390017ccd706727b45e287 to 58f5ac78e19e6cc3fb9d4048ae1a13bf364fa983 [\#1869](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1869) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/zclconf/go-cty from 1.14.1 to 1.14.2 [\#1867](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1867) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump go.mongodb.org/atlas-sdk/v20231115004 from 20231115004.0.0 to 20231115004.1.0 [\#1866](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1866) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.49.17 to 1.49.22 [\#1846](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1846) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/gruntwork-io/terratest from 0.46.9 to 0.46.11 [\#1845](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1845) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-framework from 1.4.2 to 1.5.0 [\#1844](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1844) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 2a93ea6f3eef7ba31bb9b4fe83dab787a45356fb to d9a97a5b5231f455f0390017ccd706727b45e287 [\#1825](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1825) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/gruntwork-io/terratest from 0.46.8 to 0.46.9 [\#1824](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1824) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.49.13 to 1.49.17 [\#1823](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1823) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 [\#1817](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1817) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.49.10 to 1.49.13 [\#1802](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1802) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 08975f08f935b937e16554ebd18f713b5263248a to 2a93ea6f3eef7ba31bb9b4fe83dab787a45356fb [\#1801](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1801) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.49.6 to 1.49.10 [\#1791](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1791) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.49.1 to 1.49.6 [\#1783](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1783) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-mux from 0.12.0 to 0.13.0 [\#1782](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1782) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-go from 0.19.1 to 0.20.0 [\#1781](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1781) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.30.0 to 2.31.0 [\#1780](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1780) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 5eff60fda839b96c3e34d8239dffb116c900582c to 08975f08f935b937e16554ebd18f713b5263248a [\#1779](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1779) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.14.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.14.0) (2023-12-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.13.1...v1.14.0)

**Breaking changes:**

- fix!: Replaces .String\(\) method with internal method .TimeToString\(\) to align formatting with the Atlas API. [\#1699](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1699) ([marcosuma](https://github.com/marcosuma))
   - Fixes date fields that were not compliant to the ISO 8601 timestamp format in UTC, in line with the documentation.
- feat!: New required attributes `authorized_user_first_name` and `authorized_user_last_name` in `mongodbatlas_backup_compliance_policy` resource. [\#1655](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1655) ([marcosuma](https://github.com/marcosuma))
   - Adds first and last name in `mongodbatlas_backup_compliance_policy` resource to reflect recent changes in the Atlas API.

**Enhancements**

- feat: New `mongodbatlas_stream_instance` resource [\#1685](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1685) ([AgustinBettati](https://github.com/AgustinBettati))
- feat: New `mongodbatlas_stream_instance` and `mongodbatlas_stream_instances` data sources [\#1689](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1689) [\#1701](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1701) ([AgustinBettati](https://github.com/AgustinBettati))
- feat: New `mongodbatlas_stream_connection` resource [\#1736](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1736) ([AgustinBettati](https://github.com/AgustinBettati))
- feat: New `mongodbatlas_stream_connection` and `mongodbatlas_stream_connections` data sources [\#1757](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1757) ([AgustinBettati](https://github.com/AgustinBettati))

**Bug Fixes**

- fix: Doesn't disable X.509 in the project when `mongodbatlas_x509_authentication_database_user` resource is deleted [\#1760](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1760) ([lantoli](https://github.com/lantoli))
- fix: Converts root `project_id` attribute to optional in `mongdbatlas_project_api_key` resource [\#1664](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1664) ([AgustinBettati](https://github.com/AgustinBettati))
- fix: Defines `project_assignment` block in `mongodbatlas_project_api_key` as required to avoid plugin crash [\#1663](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1663) ([AgustinBettati](https://github.com/AgustinBettati))
- fix: Fixes cluster update when adding replication specs [\#1755](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1755) ([lantoli](https://github.com/lantoli))
- fix: Provider crashes when an invalid role name is specified to `mongodbatlas_project_api_key` [\#1720](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1720) ([andreaangiolillo](https://github.com/andreaangiolillo))
- fix: Fixes string representation of id for project delete function [\#1733](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1733) ([oarbusi](https://github.com/oarbusi))

**Bug Fixes for HashiCorp Terraform Version 1.0.8**

- fix: Fixes `mongodbatlas_cloud_provider_access`, `mongodbatlas_org_invitation`, `mongodbatlas_project_api_key`, `mongodbatlas_third_party_integration` in older Terraform versions [\#1748](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1748) ([lantoli](https://github.com/lantoli))
- fix: Fixes `mongodbatlas_network_container`, `mongodbatlas_private_endpoint_regional_mode` in older Terraform versions [\#1741](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1741) ([lantoli](https://github.com/lantoli))
- fix: Fixes `mongodbatlas_serverless_instance` in older Terraform versions [\#1740](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1740) ([lantoli](https://github.com/lantoli))
- fix: Fixes some tests in `search_index` test group in older Terraform versions [\#1758](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1758) ([lantoli](https://github.com/lantoli))

**Deprecations and Removals**

- deprecate: Deprecates optional root `project_id` attribute in `mongdbatlas_project_api_key` [\#1665](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1665) ([AgustinBettati](https://github.com/AgustinBettati))

**Closed Issues**

- `mongodbatlas_alert_configuration` failing to be created on apply due to `METRIC_TYPE_UNSUPPORTED` for `DISK_PARTITION` alerts 400 error [\#1716](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1716)
- Provider produced inconsistent result after apply [\#1707](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1707)
- Error: Provider produced inconsistent final plan for `mongodbatlas_privatelink_endpoint_service` [\#1690](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1690)
- Provider not working with Secret Manager [\#1683](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1683)
- In the `mongodbatlas_advanced_cluster` ressource, forbidden characters for the  values of the tags are allowed by the provider and fail on apply. [\#1668](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1668)

**Internal Improvements**

- test: Unit test for `project_ip_access_list` resource [\#1756](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1756) ([oarbusi](https://github.com/oarbusi))
- chore: CLOUDP-215162: Update Jira GitHub Action to update Ticket status based on the issue [\#1754](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1754) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Migrate `project_ip_access_list` to new SDK [\#1753](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1753) ([oarbusi](https://github.com/oarbusi))
- chore: Fix examples job by defining beta flag [\#1751](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1751) ([AgustinBettati](https://github.com/AgustinBettati))
- test: Unit test `encryption_at_rest` resource [\#1750](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1750) ([oarbusi](https://github.com/oarbusi))
- chore: Bump github.com/aws/aws-sdk-go from 1.48.13 to 1.49.1 [\#1744](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1744) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: New scaffold command for creating resources/data sources [\#1739](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1739) ([AgustinBettati](https://github.com/AgustinBettati))
- test: Add unit test to `database_user` resource [\#1738](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1738) ([oarbusi](https://github.com/oarbusi))
- chore: Lint error [\#1734](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1734) ([andreaangiolillo](https://github.com/andreaangiolillo))
- doc: CLOUDP-215923: Removes sunset resources from documentation [\#1732](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1732) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Define atlas stream functionality under beta flag environment variable [\#1726](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1726) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: `migrate encryption_at_rest` resource to new SDK [\#1725](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1725) ([oarbusi](https://github.com/oarbusi))
- chore: Migrate database user resource to new SDK [\#1723](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1723) ([oarbusi](https://github.com/oarbusi))
- doc: CLOUDP-216288 - Update the warning message in the cluster resource to mention a bug on the container ids [\#1719](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1719) ([andreaangiolillo](https://github.com/andreaangiolillo))
- chore: Uses official CLI GH action in Terraform clean-up workflow [\#1718](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1718) ([lantoli](https://github.com/lantoli))
- chore: Change region used in online archive process region to one supported in cloud dev [\#1703](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1703) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: CLOUDP-215192 - Update Automation That Creates Github Issues To Create CLOUDP rather than INTMDB [\#1696](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1696) ([andreaangiolillo](https://github.com/andreaangiolillo))
- test: Adds unit test to project resource [\#1694](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1694) ([oarbusi](https://github.com/oarbusi))
- chore: Create Test Suite workflow [\#1687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1687) ([lantoli](https://github.com/lantoli))
- chore: Project resource migration to new sdk [\#1686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1686) ([oarbusi](https://github.com/oarbusi))
- refactor: Remove redundancy in creation functions for resources/data sources separated in packages [\#1682](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1682) ([AgustinBettati](https://github.com/AgustinBettati))
- refactor: Adjust `project_api_key` to new file structure [\#1676](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1676) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Remove references of legacy mongodbatlas package in make file [\#1675](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1675) ([AgustinBettati](https://github.com/AgustinBettati))
- chore: Updates Atlas Go SDK [\#1674](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1674) ([github-actions[bot]](https://github.com/apps/github-actions))
- chore: Increases project list limit for clean-up Github action [\#1673](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1673) ([lantoli](https://github.com/lantoli))
- test: Adds unit test to `alert_configuration` [\#1670](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1670) ([oarbusi](https://github.com/oarbusi))
- chore: Adjust CI change detection file paths after file restructure [\#1667](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1667) ([AgustinBettati](https://github.com/AgustinBettati))
- doc: Updates contributing file with code and test best practices [\#1666](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1666) ([lantoli](https://github.com/lantoli))
- chore: Adds a website make goal to preview doc changes [\#1662](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1662) ([lantoli](https://github.com/lantoli))
- test: Includes unit testing for search deployments state transition logic and model conversions [\#1653](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1653) ([AgustinBettati](https://github.com/AgustinBettati))
- feat: Migrates some resources to new file structure [\#1705](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1705) [\#1704](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1704) [\#1702](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1702) [\#1700](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1700) ([lantoli](https://github.com/lantoli))
- feat: Cleans up before and after Test Suite so they don't interfere [\#1695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1695) ([lantoli](https://github.com/lantoli))
- feat: Runs Migration tests for different provider versions [\#1691](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1691) ([lantoli](https://github.com/lantoli))
- feat: Allows to run acc/mig tests in GH action choosing version and test group. [\#1717](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1717) ([lantoli](https://github.com/lantoli))
- feat: Restructures files [\#1657](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1657) ([lantoli](https://github.com/lantoli))
- feat: Runs Acceptance and Migration tests for different Terraform CLI versions [\#1688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1688) ([lantoli](https://github.com/lantoli))
- feat: Allows to choose Terraform version for Acceptance and Migration Tests, default to latest version [\#1684](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1684) ([lantoli](https://github.com/lantoli))
- fix: Fixes config flaky tests [\#1680](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1680) ([lantoli](https://github.com/lantoli))
- fix: Fixes project flaky tests [\#1669](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1669) ([lantoli](https://github.com/lantoli))
- fix: Runs only test group if specified in Github Action acceptance or migration tests [\#1730](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1730) ([lantoli](https://github.com/lantoli))
- fix: Fixes migration tests for `backup_online_archive` [\#1724](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1724) ([lantoli](https://github.com/lantoli))
- chore: Bump tj-actions/verify-changed-files from 7b7a3b8db9077729f56bd82fced85f4b0ee67bcd to 5eff60fda839b96c3e34d8239dffb116c900582c [\#1747](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1747) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump actions/stale from 8 to 9 [\#1746](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1746) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump actions/setup-go from 4 to 5 [\#1745](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1745) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 2e93a84b325e5c4d3544924aee8afb7e1ffe189f to 1e75cac4ffa7ea5879addde1869f8fca09fce4c1 [\#1679](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1679) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.48.2 to 1.48.7 [\#1678](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1678) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/spf13/cast from 1.5.1 to 1.6.0 [\#1677](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1677) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 1e75cac4ffa7ea5879addde1869f8fca09fce4c1 to 7b7a3b8db9077729f56bd82fced85f4b0ee67bcd [\#1713](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1713) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.48.7 to 1.48.13 [\#1712](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1712) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/gruntwork-io/terratest from 0.46.7 to 0.46.8 [\#1711](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1711) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/hashicorp/terraform-plugin-testing from 1.5.1 to 1.6.0 [\#1710](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1710) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Updates Atlas Go SDK [\#1706](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1706) ([github-actions[bot]](https://github.com/apps/github-actions))

## [v1.13.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.13.1) (2023-11-23)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.13.0...v1.13.1)

**Enhancements**

- feat: Unit test resource and data source schemas in [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework-benefits) [\#1646](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1646) ([lantoli](https://github.com/lantoli))

**Bug Fixes**

- fix: uses SchemaConfigModeAttr for list attributes in `mongodbatlas_cluster` resource. [\#1654](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1654) ([marcosuma](https://github.com/marcosuma))
- fix: handles attributes as computed in `mongodbatlas_cluster` resource. [\#1642](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1642) ([marcosuma](https://github.com/marcosuma))
- fix: avoids error when removing project api key assignment for deleted project [\#1641](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1641) ([AgustinBettati](https://github.com/AgustinBettati))
- **Breaking Change**: fix!: handles paused clusters with errors when updating. [\#1640](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1640) ([marcosuma](https://github.com/marcosuma))
- fix: adds `data_process_region` field to `mongodbatlas_online_archive` resource [\#1634](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1634) ([marcosuma](https://github.com/marcosuma))

**Internal Improvements**

- fix: Update issues.yml to remove assignee [\#1649](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1649) ([andreaangiolillo](https://github.com/andreaangiolillo))
- fix: \(INTMDB-1312\) It is not possible to add breaking change label to PR [\#1647](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1647) ([andreaangiolillo](https://github.com/andreaangiolillo))
- doc: addresses outstanding comments from 1.13.0 docs release [\#1648](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1648) ([oarbusi](https://github.com/oarbusi))
- chore: migrates `mongodbatlas_alert_configuration` to new SDK [\#1630](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1630) ([oarbusi](https://github.com/oarbusi))
- chore: Bump github.com/hashicorp/terraform-plugin-go from 0.19.0 to 0.19.1 [\#1652](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1652) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump github.com/aws/aws-sdk-go from 1.47.11 to 1.48.2 [\#1651](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1651) ([dependabot[bot]](https://github.com/apps/dependabot))
- chore: Bump tj-actions/verify-changed-files from 82a523f60ad6816c35b418520f84629024d70e1e to 2e93a84b325e5c4d3544924aee8afb7e1ffe189f [\#1650](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1650) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.13.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.13.0) (2023-11-21)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.12.3...v1.13.0)

**Enhancements**

- Adds support for [MongoDB Atlas Search Node](https://www.mongodb.com/docs/atlas/atlas-search/atlas-search-overview/#search-nodes-architecture) management with `mongodbatlas_search_deployment` resource and data source [\#1633](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1633) ([AgustinBettati](https://github.com/AgustinBettati))
- Adds `type` and `fields` attributes in resource and data sources for `mongodbatlas_search_index` [\#1605](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1605) ([lantoli](https://github.com/lantoli))

**Bug Fixes**

- Fixes Terraform encryption-at-rest error when upgrading to Terraform Provider for MongoDB Atlas v1.12.2 [\#1617](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1617) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Defines `ConfigMode` for the `endpoints`  attribute, enabling use of computed nested blocks in `mongodbatlas_privatelink_endpoint_service` to support HashiCorp Terraform version 1.0.8 [\#1629](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1629) ([AgustinBettati](https://github.com/AgustinBettati))
- Makes `disk_iops` a computed attribute in `mongodbatlas_advanced_cluster` resource [\#1620](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1620) ([AgustinBettati](https://github.com/AgustinBettati))

**Closed Issues**

- `tags` not working for cluster [\#1619](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1619)
- Terraform bug in updating serverless project [\#1611](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1611)
- Changes to `mongodbatlas_project_ip_access_list` comments force a replacement [\#1600](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1600)
- Rate limit of 10 invitations per 1 minutes exceeded [\#1589](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1589)
- The `terraform-provider-mongodbatlas_v1.12.2` plugin crashed! [\#1567](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1567)
- Breaking change to drop deprecated fields made in minor version release [\#1493](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1493)

**Internal Improvements**

- Updates PR action to automatically add labels based on the PR title [\#1637](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1637) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Removes unused code in resource and cleanup in `project_api_key` docs [\#1636](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1636) ([AgustinBettati](https://github.com/AgustinBettati))
- Improves testing of `search_index` resource [\#1635](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1635) ([lantoli](https://github.com/lantoli))
- Updates Atlas Go SDK [\#1632](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1632) ([lantoli](https://github.com/lantoli))
- Updates PR template for further verifications. [\#1628](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1628) ([marcosuma](https://github.com/marcosuma))
- Removes all references to Flowdock and New Relic third-party integrations [\#1616](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1616) ([oarbusi](https://github.com/oarbusi))
- Fixes documentation errors in cloud provider access [\#1615](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1615) ([milh0use](https://github.com/milh0use))
- Fixes `mongodbatlas_search_index` acceptance tests flow [\#1610](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1610) ([lantoli](https://github.com/lantoli))
- Updates Atlas Go SDK [\#1604](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1604) ([lantoli](https://github.com/lantoli))
- Disables `event_trigger` from acceptance test due to missing cluster. [\#1603](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1603) ([marcosuma](https://github.com/marcosuma))
- Reduces stale days to 5 and close after 2 day of stale [\#1602](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1602) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes Update SDK GitHub action [\#1596](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1596) ([lantoli](https://github.com/lantoli))
- Disables `data_source` `event_trigger` tests since they are failing [\#1595](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1595) ([marcosuma](https://github.com/marcosuma))
- Changes naming convention for `data_source` event trigger test. [\#1594](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1594) ([marcosuma](https://github.com/marcosuma))
- Changes naming convention for `data_source` event trigger tests. [\#1593](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1593) ([marcosuma](https://github.com/marcosuma))
- Updates RELEASING.md [\#1592](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1592) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Adds code health to report in merge queue [\#1588](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1588) ([AgustinBettati](https://github.com/AgustinBettati))
- Adds online archive migration test to github action [\#1587](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1587) ([maastha](https://github.com/maastha))
- Changes naming convention for event trigger tests. [\#1586](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1586) ([marcosuma](https://github.com/marcosuma))
- Enables `assume_role` acceptance tests with temporary credentials [\#1585](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1585) ([marcosuma](https://github.com/marcosuma))
- Explicitly states defaults for project flags [\#1547](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1547) ([garyburgmann](https://github.com/garyburgmann))
- Bumps tj-actions/verify-changed-files from 6d688963a73d28584e163b6f62cf927a282c4d11 to 82a523f60ad6816c35b418520f84629024d70e1e [\#1626](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1626) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.47.5 to 1.47.11 [\#1625](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1625) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.29.0 to 2.30.0 [\#1624](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1624) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/gruntwork-io/terratest from 0.46.6 to 0.46.7 [\#1623](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1623) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps hashicorp/setup-terraform from 2 to 3 [\#1579](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1579) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.47.4 to 1.47.5 [\#1608](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1608) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/gruntwork-io/terratest from 0.46.1 to 0.46.6 [\#1607](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1607) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.46.3 to 1.47.4 [\#1606](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1606) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.12.3](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.12.3) (2023-11-03)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.12.2...v1.12.3)

**Enhancements**

- Adds `acceptDataRisksAndForceReplicaSetReconfig` parameter in `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` resources [\#1575](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1575) ([lantoli](https://github.com/lantoli))

**Bug Fixes**

- Overrides to attribute behavior for resource elems. in `mongodbatlas_cluster` resource [\#1572](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1572) ([marcosuma](https://github.com/marcosuma))
- Fixes `computed` and `default` usage in `mongodbatlas_cluster` resource based on the documentation [\#1564](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1564) ([marcosuma](https://github.com/marcosuma))
- Fixes global `num_shards` adding it as computed and removing the default in `mongodbatlas_cluster` resource [\#1548](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1548) ([marcosuma](https://github.com/marcosuma))

**Deprecations and Removals**

- Deprecates `page_num` and `items_per_page` in data source `mongodbatlas_search_indexes` [\#1538](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1538) ([lantoli](https://github.com/lantoli))

**Closed Issues**

- `replication_specs` do not support deep diff [\#1544](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1544)
- There is no `team_id` argument/attribute reference in `mongodbatlas_project_invitation` resource block. [\#1535](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1535)
- Provides a `mongodbatlas_privatelink_endpoint` by region when using data source `mongodbatlas_privatelink_endpoint`  [\#1525](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1525)
- Error message: Either Atlas Programmatic API Keys or AWS Secrets Manager attributes must be set [\#1483](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1483)
- Unable to update Slack `alert_configurations` via Oauth integration  [\#1074](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1074)

**Internal Improvements**

- Updates to Go 1.21.3 [\#1550](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1550) ([lantoli](https://github.com/lantoli))
- Disables `assume_role` acceptance test workflow [\#1583](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1583) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes import documentation for `project_api_key` resource [\#1582](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1582) ([lantoli](https://github.com/lantoli))
- Fixes attributes `paused`, `version_release_system` and `tags` in advanced cluster resource [\#1581](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1581) ([lantoli](https://github.com/lantoli))
- Updates run condition in migration tests github action [\#1580](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1580) ([maastha](https://github.com/maastha))
- Does not delete project for trigger acctest. [\#1573](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1573) ([marcosuma](https://github.com/marcosuma))
- Updates migration tests to run separately and use last released version of provider for plan checks [\#1565](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1565) ([maastha](https://github.com/maastha))
- Fixes aws region and aws account to be used for trigger acceptance test [\#1558](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1558) ([marcosuma](https://github.com/marcosuma))
- Adds sdk autoupdates [\#1557](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1557) ([wtrocki](https://github.com/wtrocki))
- Fixes linter cache [\#1555](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1555) ([lantoli](https://github.com/lantoli))
- Sets format for AWS region value in the provider definition [\#1549](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1549) ([marcosuma](https://github.com/marcosuma))
- Adds file .tool-versions for asdf [\#1546](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1546) ([lantoli](https://github.com/lantoli))
- Fixes the realm URL when it is set. [\#1545](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1545) ([marcosuma](https://github.com/marcosuma))
- Changes interface{} to any [\#1543](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1543) ([lantoli](https://github.com/lantoli))
- Fixes small doc bug in CHANGELOG [\#1539](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1539) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Fixes setting of authentication realm url. [\#1537](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1537) ([marcosuma](https://github.com/marcosuma))
- Migrates search index resource and data sources to new SDK [\#1536](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1536) ([lantoli](https://github.com/lantoli))
- Bumps tj-actions/verify-changed-files from 78dc414e915e0664bcf0d2b42465a86cd47bcc3c to 6d688963a73d28584e163b6f62cf927a282c4d11 [\#1562](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1562) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/hashicorp/terraform-plugin-framework from 1.4.1 to 1.4.2 [\#1561](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1561) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps google.golang.org/grpc from 1.57.0 to 1.57.1 [\#1570](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1570) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps go.mongodb.org/atlas-sdk/v20231001001 from 20231001001.0.0 to 20231001001.1.0 [\#1533](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1533) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.46.0 to 1.46.3 [\#1560](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1560) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/gruntwork-io/terratest from 0.46.0 to 0.46.1 [\#1559](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1559) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/hashicorp/hcl/v2 from 2.19.0 to 2.19.1 [\#1542](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1542) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.45.27 to 1.46.0 [\#1541](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1541) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.12.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.12.2) (2023-10-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.12.1...v1.12.2)

**Enhancements**

- Supports `data_expiration_rule` parameter in `mongodbatlas_online_archive` [\#1528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1528) ([AgustinBettati](https://github.com/AgustinBettati))
- Supports new `notifier_id` parameter in `mongodbatlas_alert_configuration` [\#1514](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1514) ([AgustinBettati](https://github.com/AgustinBettati))

**Bug Fixes**

- Fixes issue where Encryption at rest returns inconsistent plan when setting secret access key [\#1529](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1529) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes issue where alert configuration data source for third party notifications returns nil pointer [\#1513](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1513) ([AgustinBettati](https://github.com/AgustinBettati))
- Adjusts format of database user resource id as defined in previous versions [\#1506](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1506) ([AgustinBettati](https://github.com/AgustinBettati))
- Removes delete `partition_fields` statements [\#1499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1499) ([marcosuma](https://github.com/marcosuma))
- Changes validation of empty provider credentials from Error to Warning [\#1501](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1501) ([AgustinBettati](https://github.com/AgustinBettati))
- Uses `container_id` from created cluster in example [\#1475](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1475) ([mckornfield](https://github.com/mckornfield))
- Adjusts time for stale github issues to close after 1 week of inactivity [\#1512](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1512) ([AgustinBettati](https://github.com/AgustinBettati))
- Updates 1.10.0-upgrade-guide.html.markdown [\#1511](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1511) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Updates template issue with clearer guidelines [\#1510](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1510) ([marcosuma](https://github.com/marcosuma))
- Avoids including `provider_disk_type_name` property in cluster update request if attribute was removed [\#1508](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1508) ([AgustinBettati](https://github.com/AgustinBettati))

**Deprecations and Removals**

- Removes the data source `mongodbatlas_privatelink_endpoint_service_adl` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the data source `mongodbatlas_privatelink_endpoints_service_adl` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes mongodbatlas_privatelink_endpoint_service_adl [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the resource `mongodbatlas_privatelink_endpoints_service_adl` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the data source `mongodbatlas_data_lake` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the data source `mongodbatlas_data_lakes` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))
- Removes the resource `mongodbatlas_data_lake` [\#1503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1503) ([marcosuma](https://github.com/marcosuma))

**Closed Issues**

- Error changing user [\#1509](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1509)
- Provider "registry.terraform.io/mongodb/mongodbatlas" planned an invalid value [\#1498](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1498)
- IP allowlist range force replacement on 1.12.0 [\#1495](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1495)
- Importing Online Archive resources is missing parameter partition\_fields in terraform state [\#1492](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1492)
- `mongodbatlas_network_container` faulty optional variable regions [\#1490](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1490)
- Matcher not allowing null [\#1489](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1489)
- Provider version 1.12.0 is breaking the resource mongodbatlas\_database\_user \(1.11.1 works correctly\) [\#1485](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1485)
- `mongodbatlas_project_ip_access_list` causes invalid plans [\#1484](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1484)
- Changes to oplog\_min\_retention\_hours not being applied when set to null [\#1481](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1481)
- Create alert with more than 1 notification [\#1473](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1473)

**Internal Improvements**

- Migrates online archive resource and data sources to new SDK [\#1523](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1523) ([AgustinBettati](https://github.com/AgustinBettati))
- Fixes cleanup-test-env script continues if delete of one project fails [\#1516](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1516) ([AgustinBettati](https://github.com/AgustinBettati))
- Updates atlas-sdk to v20231001001 [\#1515](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1515) ([AgustinBettati](https://github.com/AgustinBettati))
- Fixes module naming convention [\#1500](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1500) ([govardhanpagidi](https://github.com/govardhanpagidi))
- Updates cluster update handler to update advanced\_configuration first and make oplog\_min\_retention\_hours non-computed [\#1497](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1497) ([maastha](https://github.com/maastha))
- Adds coverage report to PRs [\#1496](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1496) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Adds sagemaker quickstart to repo [\#1494](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1494) ([govardhanpagidi](https://github.com/govardhanpagidi))
- Closes code block in "Resource: Cloud Provider Access Configuration Paths" documentation page [\#1487](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1487) ([mattt416](https://github.com/mattt416))
- Bump github.com/gruntwork-io/terratest from 0.43.13 to 0.44.0 [\#1482](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1482) ([dependabot[bot]](https://github.com/apps/dependabot))
- Uses retry.StateChangeConf for encryption-at-rest resource. [\#1477](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1477) ([marcosuma](https://github.com/marcosuma))
- Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.28.0 to 2.29.0, github.com/hashicorp/terraform-plugin-framework from 1.3.4 to 1.4.0, github.com/hashicorp/terraform-plugin-go from 0.18.0 to 0.19.0, github.com/hashicorp/terraform-plugin-mux from 0.11.2 to 0.12.0 [\#1468](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1468) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-framework-validators from 0.10.0 to 0.12.0 [\#1466](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1466) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump golang.org/x/net from 0.13.0 to 0.17.0 [\#1524](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1524) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-framework from 1.4.0 to 1.4.1 [\#1521](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1521) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/zclconf/go-cty from 1.14.0 to 1.14.1 [\#1520](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1520) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gruntwork-io/terratest from 0.45.0 to 0.46.0 [\#1519](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1519) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.45.21 to 1.45.24 [\#1518](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1518) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/hcl/v2 from 2.18.0 to 2.18.1 [\#1517](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1517) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.45.8 to 1.45.21 [\#1505](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1505) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gruntwork-io/terratest from 0.44.0 to 0.45.0 [\#1504](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1504) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.12.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.12.1) (2023-09-22)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.12.0...v1.12.1)

**Internal Improvements**

- Updates 1.12.0 release guide and Changelog [\#1488](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1488) ([maastha](https://github.com/maastha))
- Adjusts PR template so we ensure removals and deprecations are made in isolated PRs [\#1480](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1480) ([AgustinBettati](https://github.com/AgustinBettati))

**Bug Fixes**

- Adds missing DatabaseRegex field when creating FederatedDataSource [\#1486](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1486) ([marcosuma](https://github.com/marcosuma))

**Closed Issues**

- `tags` vs. `labels` usage in `mongodbatlas_cluster` resource[\#1370](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1370)

## [v1.12.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.12.0) (2023-09-20)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.11.1...v1.12.0)

**Enhancements**

- Support for `tags` attribute in `mongodbatlas_cluster`, `mongodbatlas_advanced_cluster`, and `mongodbatlas_serverless_instance`. See [Atlas Resource Tags](https://www.mongodb.com/docs/atlas/tags/) to learn more. [\#1461](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1461) ([AgustinBettati](https://github.com/AgustinBettati))
- Support for new `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users` data sources [\#1432](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1432) ([AgustinBettati](https://github.com/AgustinBettati))

**Bug Fixes**

- **Breaking Change**: Fixes an issue where removing `collectionName` from user role doesn't work [\#1471](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1471) ([marcosuma](https://github.com/marcosuma)).
  - Note: As a result, `mongodbatlas_database_user` no longer requires `roles.collection_name` attribute and doesn't support an empty `collection_name`. You should remove any usage of `roles.collection_name = ""` in configurations for this resource when you upgrade to this version. For more details see:  https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.12.0-upgrade-guide. 
- Populates `total_count` in `mongodbatlas_alert_configurations` data source  [\#1476](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1476) ([lantoli](https://github.com/lantoli))
- Improves error handling for `cloud_backup_schedule` resource. [\#1474](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1474) ([marcosuma](https://github.com/marcosuma))
- Handles incorrect ids when importing `alert_configuration` or `project_ip_access_list` resources [\#1472](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1472) ([lantoli](https://github.com/lantoli))
- Changelog Spelling Fixes  [\#1457](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1457) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Adds `mongodbatlas_database_user` username parameter OIDC footnote in docs [\#1458](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1458) ([Zuhairahmed](https://github.com/Zuhairahmed))

**Deprecations and Removals**

- Deprecation of `labels` attribute in `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` resources.
- Remove deprecated fields in `mongodbatlas_alert_configuration` resource [\#1385](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1385) ([AgustinBettati](https://github.com/AgustinBettati))
- Removal of `api_keys` attribute from `mongodbatlas_project` [\#1365](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1365) ([maastha](https://github.com/maastha))
- Removal of attributes in `mongodbatlas_encryption_at_rest` resource: aws_kms, azure_key_vault, google_cloud_kms [\#1383](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1383) ([maastha](https://github.com/maastha))
- Removal of MongoDB Atlas Terraform Provider v1.12.0 deprecated fields. [\#1418](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1418) ([marcosuma](https://github.com/marcosuma))
  - provider: duration_seconds
  - advanced_cluster resource: bi_connector
  - cloud_backup_snapshot_restore_job resource: delivery_type
  - cloud_provider_access_setup resource: aws
  - cluster resource: bi_connector, provider_backup_enabled, aws_private_link, aws_private_link_srv
  - database_user resource: provider_backup_enabled
  - project_api_key resource: role_names
  - cluster and clusters data sources: bi_connector
  - project_key and project_keys data sources: role_names

**Closed Issues**

- Alert notification interval\_min not working [\#1464](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1464)
- Changing DB user role from readwrite@DB.Col1 to readwrite@DB doesn't work [\#1462](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1462)
- Unable to deploy a DISK\_PARTITION\_UTILIZATION\_DATA AlertConfiguration [\#1410](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1410)
- Error: The terraform-provider-mongodbatlas\_v1.11.0 plugin crashed [\#1396](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1396)
- Trigger resource does not force replacement when app id changes [\#1310](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1310)

**Internal Improvements**

- Bump goreleaser/goreleaser-action from 4 to 5 [\#1470](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1470) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.45.4 to 1.45.8 [\#1469](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1469) ([dependabot[bot]](https://github.com/apps/dependabot))
- Merge feature branch into master [\#1460](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1460) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Bump github.com/aws/aws-sdk-go from 1.45.2 to 1.45.4 [\#1459](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1459) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-testing from 1.4.0 to 1.5.1 [\#1455](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1455) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/zclconf/go-cty from 1.13.3 to 1.14.0 [\#1454](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1454) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/hcl/v2 from 2.17.0 to 2.18.0 [\#1453](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1453) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump actions/checkout from 3 to 4 [\#1452](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1452) ([dependabot[bot]](https://github.com/apps/dependabot))
- Fix docs and example fix for project\_api\_key resource after removing role\_names deprecated field [\#1441](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1441) ([AgustinBettati](https://github.com/AgustinBettati))
- Add breaking changes strategy for Terraform [\#1431](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1431) ([wtrocki](https://github.com/wtrocki))
- Extract Configure and Metadata framework functions into single implementation [\#1424](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1424) ([AgustinBettati](https://github.com/AgustinBettati))
- Fix INTMDB-1017 - Updated alert configuration schema with required params [\#1421](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1421) ([andreaangiolillo](https://github.com/andreaangiolillo))
- IP Access List doc updates for Terraform Resources, Data Sources [\#1414](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1414) ([zach-carr](https://github.com/zach-carr))
- Avoid diff in state after import for undefined optional attribute in alert config notification [\#1412](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1412) ([AgustinBettati](https://github.com/AgustinBettati))
- Migrate Resource: mongodbatlas\_project\_ip\_access\_list to Terraform Plugin Framework [\#1411](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1411) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Remove validation of empty public\_key and private\_key attributes in provider config to avoid breaking change [\#1402](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1402) ([AgustinBettati](https://github.com/AgustinBettati))
- Migrate DataSource mongodbatlas\_alert\_configuration to Terraform Plugin Framework [\#1397](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1397) ([AgustinBettati](https://github.com/AgustinBettati))
- Migrate DataSource: mongodbatlas\_project\_ip\_access\_list to Terraform Plugin Framework [\#1395](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1395) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Migrate Resource: mongodbatlas\_database\_user to Terraform Plugin Framework [\#1388](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1388) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update project documentation to remove api\_keys references [\#1386](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1386) ([maastha](https://github.com/maastha))
- Migrates `mongodbatlas_alert_configuration` resource and removes deprecated fields [\#1385](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1385) ([AgustinBettati](https://github.com/AgustinBettati))
- Prepares migration to Terraform framework [\#1384](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1384) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Migrates `mongodbatlas_encryption_at_rest` resource to Terraform Plugin Framework [\#1383](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1383) ([maastha](https://github.com/maastha))
- Adds new framework provider, main and acceptance tests to use mux server with existing sdk v2 provider [\#1366](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1366) ([AgustinBettati](https://github.com/AgustinBettati))
- Migrates `mongodbatlas_project` resource to Terraform Plugin Framework and remove api\_keys attribute [\#1365](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1365) ([maastha](https://github.com/maastha))

## [v1.11.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.11.1) (2023-09-06)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.11.0...v1.11.1)

**Enhancements**

- Adds Atlas OIDC Database User support to `mongodbatlas_database_user` [\#1382](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1382) ([andreaangiolillo](https://github.com/andreaangiolillo))
  - Note: This feature is only available starting in [MongoDB 7.0](https://www.mongodb.com/evolved#mdbsevenzero) or later. To learn more see https://www.mongodb.com/docs/atlas/security-oidc/  
- Adds Atlas `datasetNames` support in `mongodbatlas_federated_database_instance` [\#1439](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1439) ([maastha](https://github.com/maastha))
- Improves `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` documentation to highlight that `provider_volume_type=STANDARD` is not available for NVMe clusters [\#1430](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1430) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Adds a new example for `mongodbatlas_online_archive` [\#1372](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1372) ([rsolovev](https://github.com/rsolovev))
- Adds a new example for `mongodbatlas_cloud_backup_schedule` to create policies for multiple clusters [\#1403](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1403) ([maastha](https://github.com/maastha))


**Bug Fixes**

- Updates `tag_sets` to `storage_stores.read_preference` in `mongodbatlas_federated_database_instance` [\#1440](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1440) ([maastha](https://github.com/maastha))
- Updates cluster documentation about labels field [\#1425](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1425) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes null pointer error in `mongodbatlas_alert_configuration` [\#1419](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1419) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Updates `mongodbatlas_event_trigger` resource to force replacement when app id changes [\#1387](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1387) ([maastha](https://github.com/maastha))
- Updates deprecation message to 1.12.0 [\#1381](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1381) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes null pointer error in `mongodbatlas_project` data source [\#1377](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1377) ([andreaangiolillo](https://github.com/andreaangiolillo))


**Closed Issues**

- Provider registry registry.terraform.io does not have a provider named registry.terraform.io/hashicorp/mongodbatlas [\#1389](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1389)
- Event Trigger resource doesn't support wildcard collection name [\#1374](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1374)
- 1.11.0 - terraform provider `mongodbatlas_projects` access denied  [\#1371](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1371)
- Terraform-provider-mongodbatlas\_v1.10.2 plugin crashes when including backup schedule [\#1368](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1368)
- Error: Plugin did not respond - panic: interface conversion: interface is nil, not map[string]interface [\#1337](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1337)
- "Error: error deleting organization information" When importing organization [\#1327](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1327)
- `instance_size` for advance cluster marked as optional in the documentation [\#1311](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1311)
- `region_configs.#._specs.instance_size` in `mongodbatlas_advanced_cluster` is required [\#1288](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1288)

**Internal Improvements**

- Updates the release flow to remove the acceptance steps [\#1443](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1443) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Bumps github.com/aws/aws-sdk-go from 1.44.334 to 1.45.2 [\#1442](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1442) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.27.0 to 2.28.0 [\#1429](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1429) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/zclconf/go-cty from 1.13.2 to 1.13.3 [\#1428](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1428) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/gruntwork-io/terratest from 0.43.12 to 0.43.13 [\#1427](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1427) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.44.329 to 1.44.334 [\#1426](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1426) ([dependabot[bot]](https://github.com/apps/dependabot))
- Removes 3rd shard from 2 shard global cluster example [\#1423](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1423) ([BassT](https://github.com/BassT))
- Updates issue.yml to use issue number as Ticket title [\#1422](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1422) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Migrates to new Atlas SDK - `mongodbatlas_federated_database_instance` resource [\#1415](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1415) ([maastha](https://github.com/maastha))
- Updates broken links to the Atlas Admin API docs [\#1413](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1413) ([zach-carr](https://github.com/zach-carr))
- Self document make [\#1407](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1407) ([wtrocki](https://github.com/wtrocki))
- Adds instructions for updates of the Atlas SDK [\#1406](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1406) ([wtrocki](https://github.com/wtrocki))
- Bumps github.com/mongodb-forks/digest from 1.0.4 to 1.0.5 [\#1405](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1405) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.44.324 to 1.44.329 [\#1404](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1404) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps golangci/golangci-lint-action from 3.6.0 to 3.7.0 [\#1393](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1393) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/aws/aws-sdk-go from 1.44.319 to 1.44.324 [\#1392](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1392) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps github.com/gruntwork-io/terratest from 0.43.11 to 0.43.12 [\#1391](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1391) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bumps go.mongodb.org/atlas from 0.32.0 to 0.33.0 [\#1390](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1390) ([dependabot[bot]](https://github.com/apps/dependabot))
- Improves the release process [\#1380](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1380) ([andreaangiolillo](https://github.com/andreaangiolillo))
- clenaup-test-env.yml [\#1379](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1379) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Creates github action to delete projects in the test env [\#1378](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1378) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Bumps github.com/aws/aws-sdk-go from 1.44.314 to 1.44.319 [\#1375](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1375) ([dependabot[bot]](https://github.com/apps/dependabot))
- Adds githooks [\#1373](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1373) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Fixes cluster outage tests [\#1364](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1364) ([andreaangiolillo](https://github.com/andreaangiolillo))

## [v1.11.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.11.0) (2023-08-04)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.10.2...v1.11.0)

**Enhancements**

- [Azure Service Principles](https://learn.microsoft.com/en-us/azure/active-directory/develop/app-objects-and-service-principals?tabs=browser) support in `mongodbatlas_cloud_provider_access_setup` and `mongodbatlas_cloud_provider_access_authorization` [\#1343](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1343) ([andreaangiolillo](https://github.com/andreaangiolillo)) - INTMDB-545
- Atlas [Shared Cluster Backups](https://www.mongodb.com/docs/atlas/backup/cloud-backup/shared-cluster-backup/) support in `mongodbatlas_shared_tier_snapshot` and `mongodbatlas_shared_tier_restore_job` [\#1324](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1324) ([andreaangiolillo](https://github.com/andreaangiolillo)) and [\#1323](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1323) ([andreaangiolillo](https://github.com/andreaangiolillo)) - INTMDB-546
- Atlas Project `limits` support in `mongodbatlas_project` [\#1347](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1347) ([AgustinBettati](https://github.com/AgustinBettati)) - INTMDB-554
- New example for Encryption at Rest using Customer Key Management and multi-region cluster [\#1349](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1349) ([andreaangiolillo](https://github.com/andreaangiolillo)) - INTMDB-340

**Deprecations and Removals**   

- Marking `cloud_provider_access` resource and data source as deprecated [\#1355](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1355) ([AgustinBettati](https://github.com/AgustinBettati)) - INTMDB-967	

**Bug Fixes**

- Update `mongodbatlas_cloud_backup_schedule` to add the ID field to policyItems [\#1357](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1357) ([andreaangiolillo](https://github.com/andreaangiolillo))
- `project_api_key` data source missing `project_assignment` attribute [\#1356](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1356) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Support update of description for project api key resource [\#1354](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1354) ([AgustinBettati](https://github.com/AgustinBettati))
- Null pointer in `resource_mongodbatlas_cloud_backup_schedule` [\#1353](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1353) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Plugin did not respond - panic: interface conversion: interface is nil, not map[string]interface [\#1342](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1342) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Error deleting organization information when importing organization [\#1352](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1352) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Mark project api key resource as destroyed if not present [\#1351](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1351) ([AgustinBettati](https://github.com/AgustinBettati))
- `mongodbatlas_privatelink_endpoint_service` data source doc bug fix [\#1334](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1334) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Make region atributed optional computed in third-party-integration [\#1332](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1332) ([maastha](https://github.com/maastha))

**Closed Issues**

- json: cannot unmarshal number 4841168896 into Go struct field CloudProviderSnapshot.storageSizeBytes of type int [\#1333](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1333)
- Labels is not creating tags in the MongoAtlas UI [\#1319](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1319)
- `mongodbatlas_online_archive` `schedule` parameter update causing crashing in `terraform apply` [\#1318](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1318)
- Update Pager Duty integration fails with INTEGRATION\_FIELDS\_INVALID [\#1316](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1316)
- mongodbatlas\_event\_trigger is not updated if config\_match is added [\#1302](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1302)
- Updating the 'name' field of a 'mongodbatlas\_project' recreates a new Project [\#1296](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1296)
- mongodbatlas\_org\_invitation is missing ORG\_BILLING\_READ\_ONLY role support [\#1280](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1280)
- mongodbatlas\_alert\_configuration notification microsoft\_teams\_webhook\_url is always updated [\#1275](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1275)
- Provider not destroying API keys [\#1261](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1261)
- Has `project_assignment` of `mongodbatlas_api_key` not been implemented? [\#1249](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1249)
- Invalid attribute providerBackupEnabled specified. [\#1245](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1245)

**Internal Improvements**

- Fix documentation for `mongodbatlas_api_key` [\#1363](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1363) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Improve self-managed x509 database user docs [\#1336](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1336) ([andreaangiolillo](https://github.com/andreaangiolillo))
- add prefix to dependabot PR [\#1361](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1361) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update README.md with supported OS/Arch [\#1350](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1350) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Add PR lint to repo [\#1348](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1348) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Mark `instance_size` in electable specs required in `advanced_cluster` documentation [\#1339](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1339) ([maastha](https://github.com/maastha))
- Update RELEASE.md Github issue [\#1331](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1331) ([maastha](https://github.com/maastha))
- Update privatelink endpoint service resources timeout config [\#1329](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1329) ([maastha](https://github.com/maastha))
- Use go-version-file in github actions [\#1315](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1315) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Add autogenerated SDK to terraform [\#1309](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1309) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Migrate to terraform-plugin-testing [\#1301](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1301) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Bump github.com/aws/aws-sdk-go from 1.44.308 to 1.44.314 [\#1360](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1360) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gruntwork-io/terratest from 0.43.10 to 0.43.11 [\#1358](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1358) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/terraform-plugin-testing from 1.3.0 to 1.4.0 [\#1346](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1346) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gruntwork-io/terratest from 0.43.8 to 0.43.10 [\#1345](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1345) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.44.304 to 1.44.308 [\#1344](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1344) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.44.302 to 1.44.304 [\#1335](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1335) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/aws/aws-sdk-go from 1.44.299 to 1.44.302 [\#1330](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1330) ([dependabot[bot]](https://github.com/apps/dependabot))


## [v1.10.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.10.2) (2023-07-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.10.1...v1.10.2)

**Bug Fixes:**

- `mongodbatlas_advanced_cluster` doc is not formatted correctly [\#1326](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1326) - INTMDB-941 ([andreaangiolillo](https://github.com/andreaangiolillo))
- `mongodbatlas_event_trigger` is not updated if `config_match` is added [\#1305](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1305) - INTMDB-922 ([andreaangiolillo](https://github.com/andreaangiolillo))
- `mongodbatlas_online_archive` `schedule` parameter update causing crashing in terraform apply - INTMDB-935 [\#1320](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1320) ([andreaangiolillo](https://github.com/andreaangiolillo))

**Internal Improvements:**

- Fix `mongodbatlas_online_archive` tests - INTMDB-938 [\#1321](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1321) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.298 to 1.44.299 [\#1312](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1312) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.10.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.10.1) (2023-7-13)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.10.0...v1.10.1)

**Enhancements:**

- Support for updating the name field of `mongodbatlas_project` without recreating a new Project - INTMDB-914 [\#1298](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1298) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Support for `federation_settings_id` parameter in `mongodbatlas_organization` to enable linking to an existing federation upon Create - INTMDB-838 [\#1289](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1289) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Support for `schedule` parameter in resource `mongodbatlas_online_archive` - INTMDB-828 [\#1272](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1272) ([andreaangiolillo](https://github.com/andreaangiolillo))
- New `mongodbatlas_advanced_cluster` doc examples for Multi-Cloud Clusters and Global Clusters - INTMDB-442 [\#1256](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1256) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Support for `transactionLifetimeLimitSeconds` parameter in `mongodbatlas_cluster` and `mongodbatlas_advanced_cluser` - INTMDB-874 [\#1252](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1252) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Add Troubleshooting.md to include issue with using `dynamic` in Terraform - INTMDB-855 [\#1240](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1240) ([maastha](https://github.com/maastha))

**Bug Fixes:**
- Remove default value to [retainBackups parameter](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Clusters/operation/deleteCluster) in `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` - INTMDB-932 [\#1314](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1314) ([andreaangiolillo](https://github.com/andreaangiolillo))
- `mongodbatlas_cloud_backup_snapshot_restore_job` extend guards for delivery type deletions - INTMDB-919 [\#1300](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1300) ([andreaangiolillo](https://github.com/andreaangiolillo))
- `mongodbatlas_org_invitation` is missing `ORG_BILLING_READ_ONLY` role - INTMDB-904 [\#1287](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1287) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Change warning to error for org key delete - INTMDB-889 [\#1283](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1283) ([martinstibbe](https://github.com/martinstibbe))
- Add MicrosoftTeamsWebhookURL to values that are based on schema vs API - INTMDB-896 [\#1279](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1279) ([martinstibbe](https://github.com/martinstibbe))
- Update `group_id` -\> `project_id` for backup snapshots DOCSP-30798 [\#1273](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1273) ([jwilliams-mongo](https://github.com/jwilliams-mongo))
- Update example documentation for `mongodbatlas_project_api_key` - INTMDB-876 [\#1265](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1265) ([martinstibbe](https://github.com/martinstibbe))
- Make sure failed Terraform run rolls back properly - INTMDB-433 [\#1264](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1264) ([martinstibbe](https://github.com/martinstibbe))
- Fix "pause" during cluster and `mongodbatlas_advanced_cluster` update [\#1248](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1248) ([maastha](https://github.com/maastha))
- Add ForceNew to audit if the project id changes - INTMDB-435 [\#1247](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1247) ([martinstibbe](https://github.com/martinstibbe))

**Closed Issues:**

- `mongodbatlas_alert_configuration` failing to be created on apply due to `METRIC_TYPE_UNSUPPORTED` 400 error [\#1242](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1242)
- `app_id` property doesn't work in `mongodbatlas_event_trigger` resource [\#1224](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1224)
- Random change in `region_configs` order of `mongodbatlas_advanced_cluster` [\#1204](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1204)
- Problem returning some fields from `mongodbatlas_advanced_cluster` [\#1189](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1189)
- `node_count` documented as read-only for `mongodbatlas_advanced_cluster` [\#1187](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1187)
- In `mongodbatlas_third_party_integration` the `microsoft_teams_webhook_url` parameter keeps updating on every apply [\#1135](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1135)
- Point in Time Restore is not enabled when `should_copy_oplogs` is set to `true`, when copying backups to other regions [\#1134](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1134)
- Documentation: `analyzer` argument in Atlas search index is required [\#1132](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1132)
- Serverless Instance wants to do an in-place update on every run [\#1070](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1070)

**Internal Improvements:**

- INTMDB-912: Generate the CHANGELOG.md [\#1307](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1307) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Use GET one for cloud provider access to improve existing workflow - INTMDB-137 [\#1246](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1246) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.7 to 0.43.8 [\#1306](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1306) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.295 to 1.44.298 [\#1304](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1304) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.6 to 0.43.7 [\#1303](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1303) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-913: \[Terraform\] Enable fieldalignment linter [\#1297](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1297) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.294 to 1.44.295 [\#1293](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1293) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-911: \[Terraform\] Remove unused secret from code-health workflow [\#1291](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1291) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-910: \[Terraform\] Remove Automated Tests workflow [\#1290](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1290) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.291 to 1.44.294 [\#1286](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1286) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.26.1 to 2.27.0 [\#1285](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1285) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.5 to 0.43.6 [\#1284](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1284) ([dependabot[bot]](https://github.com/apps/dependabot))
- Remove slack key from repo [\#1282](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1282) ([maastha](https://github.com/maastha))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.3 to 0.43.5 [\#1277](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1277) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.287 to 1.44.291 [\#1276](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1276) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-895: Third\_Party\_Integrations region field is required parameter in Terraform [\#1274](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1274) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Update RELEASING.md [\#1271](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1271) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-881: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - "Acceptance Tests" [\#1269](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1269) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-892: \[Terraform\] Add APIx-Integration as a reviewer of dependabot PR [\#1268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1268) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.286 to 1.44.287 [\#1267](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1267) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.2 to 0.43.3 [\#1266](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1266) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-879: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - Config [\#1263](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1263) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-878: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - Network [\#1260](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1260) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.284 to 1.44.286 [\#1259](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1259) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.43.0 to 0.43.2 [\#1258](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1258) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-880: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - Project [\#1257](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1257) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-883: Fix "Create JIRA ticket" Action [\#1255](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1255) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update README.md [\#1254](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1254) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update RELEASING.md [\#1253](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1253) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.283 to 1.44.284 [\#1251](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1251) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-871: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev - backup [\#1250](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1250) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Update acceptance-tests.yml [\#1244](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1244) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.282 to 1.44.283 [\#1243](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1243) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.272 to 1.44.282 [\#1237](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1237) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump golangci/golangci-lint-action from 3.4.0 to 3.6.0 [\#1235](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1235) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-843: \[Terraform\] Improve acceptance test setup to run in parallel & against cloud-dev. Clusters tests [\#1234](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1234) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-249: \[Terraform\] Lint and fix linting for examples [\#1221](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1221) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-808: Using vars instead of secrets for not sensitive info [\#1220](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1220) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump octokit/request-action from 2.1.7 to 2.1.9 [\#1211](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1211) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/hcl/v2 from 2.16.2 to 2.17.0 [\#1206](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1206) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.10.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.10.0) (2023-6-15)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.9.0...v1.10.0)

**Enhancements:**

- New: [Organizations Management](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/createOrganization) including Create (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1176) - INTMDB-533  
- New: [Federated Database Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createFederatedDatabase) resource and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1163) - INTMDB-801
- New: [Query Limit for Database Instance](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createOneDataFederationQueryLimit) resource and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1173) - INTMDB-802
- New: [Private Endpoint](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Federation/operation/createDataFederationPrivateEndpoint) resources and data sources for Federated Database Instance and Online Archive (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1182) - INTMDB-803
- New: [Data Lake Pipelines](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/createPipeline) resource and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1174) - INTMDB-804
- New: [Data Lake Pipelines Run](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Data-Lake-Pipelines/operation/getPipelineRun) data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1177) - INTMDB-805
- New: [Cluster Outage Simulation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cluster-Outage-Simulation/operation/startOutageSimulation) resource and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1188) - INTMDB-835
- Feature Add: [Cluster Enable Extended Storage](https://www.mongodb.com/docs/atlas/customize-storage/#minimum-disk-capacity-to-ram-ratios) Sizes in `mongodbatlas_project` (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1128) - INTMDB-466
- Feature Add: [srvShardOptimizedConnectionString parameter](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Multi-Cloud-Clusters/operation/createCluster) to `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` Data Sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1157) - INTMDB-694
- Feature Add: [retainBackups parameter](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Multi-Cloud-Clusters/operation/createCluster) to `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1210) - INTMDB-781
- [Programmatic API Key](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Programmatic-API-Keys) Resource Updates (`mongodbatlas_api_key`, `mongodbatlas_project_api_key` and `mongodbatlas_project_ip_access_list_key`) + Doc Cleanup (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1208) - INTMDB-655
- Release.md File Updates with Action Items (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1203) - INTMDB-690
- ChangeLog Generator (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1165) - INTMDB-720
- Upgrade to [Go 1.20](https://go.dev/blog/go1.20) (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1153) - INTMDB-658

**Deprecations and Removals:**

- `mongodbatlas_data_lake` and `mongodbatlas_privatelink_endpoint_service_adl` (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1190) - INTMDB-806
-  Remove `mongodbatlas_private_ip_mode` and NEW_RELIC and FLOWDOCK in `mongodbatlas_third_party_integration` resources and data sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1159) - INTMDB-408
-  Remove mongodbatlas_cloud_provider (access, snapshot, snapshot_backup_policy, snapshot_restore_job) resources and datas sources (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1159) - INTMDB-408

**Bug Fixes:**

- `mongodbatlas_serverless_instance` wants to do an in-place update on every run (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1152) - INTMDB-710
- Documentation bug: analyzer argument in `mongodbatlas_search_index` is required (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1158) - INTMDB-780
- Point in Time Restore is not enabled when should_copy_oplogs is set to true, when copying backups to other regions (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1150) - INTMDB-783
- `mongodbatlas_third_party_integration` - microsoft_teams_webhook_url keeps updating on every apply (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1148) - INTMDB-784
- In `mongodbatlas_database_user` usernames with spaces breaks state due to URL encoding (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1164) - INTMDB-809
- `mongodbatlas_backup_compliance_policy` causing `mongodbatlas_cloud_backup_schedule` resource to fail (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1209) - INTMDB-827 
- `mongodbatlas_advanced_cluster` `node_count` parameter doc bug fix (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1193) - INTMDB-844
- Fix typos in docs for `mongodbatlas_network_peering` and  `mongodbatlas_network_container` resource imports (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1200) 

**Closed issues:**

- Online Archive: "Specified cloud provider is not supported in Data Lake Storage" but works in UI [\#1216](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1216)
- Asymmetric hardware is not supported by the provider [\#1214](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1214)
- `region` argument missing from `mongodbatlas_third_party_integration` for integration with PagerDuty [\#1180](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1180)
- Correct docs for importing network peering resources [\#1179](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1179)
- Terraform destroy produces a 500 \(UNEXPECTED\_ERROR\) on the underlying API call [\#1162](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1162)
- produced an unexpected new value: Root resource was present, but now │ absent [\#1160](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1160)
- Failed to respond to the plugin.\(\*GRPCProvider\).PlanResourceChange call [\#1136](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1136)
- Error: error creating MongoDB Cluster: unexpected EOF [\#674](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/674)

**Internal Improvements:**

- Chore\(deps\):  Bump actions/stale from 7 to 8 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1117)
- Chore\(deps\):  Bump github.com/zclconf/go-cty from 1.13.1 to 1.13.2 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1184)
- Chore\(deps\):  Bump github.com/aws/aws-sdk-go from 1.44.268 to 1.44.272 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1196)
- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.42.0 to 0.43.0 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1197)
- Chore(deps): Bump github.com/spf13/cast from 1.5.0 to 1.5.1 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1195)
- Chore(deps): Bump go.mongodb.org/atlas from 0.25.0 to 0.28.0 (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1194)
- corrected documentation for advanced cluster and cluster (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1167) 
- Update component field to "Terraform" GitHub action (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1178)
- Add action to create JIRA ticket for a new Github issue (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1166)
- Terraform Provider Secrets Audit (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1202) 
- Add Code-health action (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1144)  
- Fix TestAccConfigDSSearchIndexes_basic (https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1227) 

## [v1.9.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.9.0) (2023-4-27)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.8.2...v1.9.0)

**Enhancements:**

- New Feature: [Backup Compliance Policy](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/) Support  [\#1127](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1127) - INTMDB-587

**Bug Fixes:**

- Update resource [mongodbatlas_project](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project) to detect null boolean values in project settings [\#1145](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1145) - INTMDB-789
- Update on resource [mongodbatlas_search_index]([https://www.mongodb.com/docs/atlas/atlas-search/create-index/](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/search_index)) resource docs [\#1137](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1137) - DOCSP-28948
- Removing resource [mongodbatlas_cluster](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cluster) `disk_size_gb` examples [\#1133](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1133)

**Closed issues:**

- mongodbatlas\_search\_index does not change name [\#1096](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1096)
- Unhelpful error when importing advanced cluster using mongodbatlas\_cluster resource [\#1089](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1089)
- Update Slack alert configuration fails with INTEGRATION\_FIELDS\_INVALID [\#1086](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1086)
- Upgrade to terraform-plugin-sdk v2.25.0 [\#1080](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1080)
- mongodbatlas\_project\_ip\_access\_list.comment should be optional [\#1079](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1079)
- Can't unset auto\_scaling [\#1072](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1072)
- mongodbatlas\_access\_list\_api\_key fails to import [\#1064](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1064)
- Terraform plan fails if API key created by `mongodbatlas_api_key` resource is deleted outside of Terraform [\#1057](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1057)
- mongodbatlas\_search\_index does not recreate when cluster\_name and project\_id fields change. [\#1053](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1053)
- ERROR: Pager Duty API key must consist of 32 hexadecimal digits [\#1049](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1049)
- mongodbatlas\_alert\_configuration not detecting drift [\#999](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/999)
- Provider insists on changing a sub-parameter even when no changes are necessary [\#997](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/997)
- Type: TEAM alert notification not saved properly [\#971](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/971)
- `app_id` In documentation is ambiguous for MongoDB Atlas Event Trigger [\#957](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/957)
- Provider panic with `authentication_enabled=true` input [\#873](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/873)
- Schema error when creating event\_trigger referring to Atlas App Services function  [\#858](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/858)

## [v1.8.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.8.2) (2023-3-30)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.8.1...v1.8.2)

**Enhancements:**

- Support for "TIMESERIES" Collection Type in [`mongodbatlas_online_archive`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/online_archive) [\#1114](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1114) - INTMDB-648
- Support for new "DATADOG" regions  in [`mongodbatlas_third_party_integration`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/third_party_integration) [\#1105](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1105) - INTMDB-638

**Bug Fixes:**

- Error in unsetting auto_scaling in [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) [\#1112](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1112) - INTMDB-617
- Update Status Codes in  [`mongodbatlas_search_index`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/search_index) [\#1104](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1104) - INTMDB-687
- [`mongodbatlas_project_ip_access_list`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project_ip_access_list) comment should be optional [\#1103](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1103) - INTMDB-637
- Plan fails if API key created by [`mongodbatlas_api_key`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/api_key) resource is deleted outside of Terraform [\#1097](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1097) - INTMDB-581
- Google Cloud Terraform Provider Test Version Upgrade Refactoring [\#1098](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1098) - INTMDB-359


**Closed Issues:**

- mongodbatlas\_federated\_settings\_org\_role\_mapping INVALID\_ATTRIBUTE [\#1110](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1110)
- Errors when creating or importing timeseries online archive [\#1081](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1081)
- Unexpected EOF [\#1083](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1083)
- mongodbatlas\_access\_list\_api\_key creation fails after api\_key creation [\#1075](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1075)
- Panic when creating AWS privatelink endpoint [\#1067](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1067)
- error getting project 403 \(request IP\_ADDRESS\_NOT\_ON\_ACCESS\_LIST\) even if whitelisted IP [\#1048](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1048)
- Crash during the changes for `mongodbatlas_cloud_backup_schedule` interval with dynamic blocks. [\#1041](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1041)
- Regression in 1.8.0: mongodbatlas\_third\_party\_integration marks "type" attribute as deprecated [\#1032](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1032)
- Cannot reference `container_id` for `mongodbatlas_advanced_cluster` [\#1008](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1008)
- ERROR: Datadog API key must consist of 32 hexadecimal digits [\#1001](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1001)
- num\_shards value changed to 1 [\#970](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/970)
- Segmentation Fault in TerraForm Provider [\#969](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/969)
- Import statements are broken in documentation and help commands in the Terraform provider are outdated. [\#956](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/956)
- Invitation handling is not working after user accepted invitation. [\#945](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/945)
- bi\_connector settings don't work in mongodbatlas\_advanced\_cluster [\#893](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/893)
- Sensitive information exposed as resource Id - mongodbatlas\_x509\_authentication\_database\_user  [\#884](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/884)
- plugin crashes during apply: panic: runtime error: invalid memory address or nil pointer dereference [\#866](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/866)

**Internal Improvements:**

- Release staging v.1.8.2 [\#1115](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1115) ([martinstibbe](https://github.com/martinstibbe))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.221 to 1.44.226 [\#1109](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1109) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.25.0 to 2.26.1 [\#1108](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1108) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/zclconf/go-cty from 1.12.1 to 1.13.1 [\#1107](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1107) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump actions/setup-go from 3 to 4 [\#1106](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1106) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/hcl/v2 from 2.16.1 to 2.16.2 [\#1101](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1101) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.24.1 to 2.25.0 [\#1100](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1100) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.44.216 to 1.44.221 [\#1099](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1099) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/aws/aws-sdk-go from 1.40.56 to 1.44.216 [\#1094](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1094) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump go.mongodb.org/atlas from 0.21.0 to 0.23.1 [\#1092](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1092) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump golang.org/x/net from 0.1.0 to 0.7.0 [\#1071](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1071) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.8.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.8.1) (2023-3-7)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.8.0...v1.8.1)

**Enhancements:**

- Upgrade to [go1.19](https://go.dev/blog/go1.19) [\#1031](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1031) - INTMDB-390
- Add configurable timeouts to resources that wait for [`cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cluster) to become IDLE [\#1047](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1047) - INTMDB-508
- Improve [`cloud_provider_access_authorization`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_provider_access) and [`encryption_at_rest`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/encryption_at_rest) (remove need for `time_sleep` arguments) [\#1045](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1045) - INTMDB-560
- Add [`search_index`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/search_index) error handling [\#1077](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1077) - INTMDB-600
- New / Improved Upon [Resource Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples): 
  - [`cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cluster) with NVMe Upgrade ([Non-Volatile Memory Express](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ssd-instance-store.html)) [\#1037](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1037) - INTMDB-32. See example [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-cluster)
  - [`privatelink_endpoint_serverless`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/privatelink_endpoint_serverless) Examples for AWS + Azure [\#1043](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1043) - INTMDB-424. See example for [AWS](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/aws-atlas-privatelink-serverless) and [Azure](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/azure-serverless)
  - Improvement for [`private_link_endpoint`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/privatelink_endpoint) [\#1082](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1082) - INTMDB-410. see example [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/aws-privatelink-endpoint)
  - Improvement for [`encryption_at_rest`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/encryption_at_rest) [\#1060](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1060). see example [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-encryption-at-rest/aws)

**Bug Fixes:**

- Resource [`ldap_configuration`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/ldap_configuration) broken [\#1033](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1033) - INTMDB-440
- [`event_trigger`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/event_trigger) Import statements are broken in documentation [\#1046](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1046) - INTMDB-513
- [`event_trigger`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/event_trigger) Error Handler Update [\#1061](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1061) - INTMDB-517
- TEAM alert notification not saved properly [\#1029](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1029) - INTMDB-529
- [`alert_configuration`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/alert_configuration) not detecting drift [\#1030](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1030) - INTMDB-542
- [`third_party_integration`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/third_party_integration) marks "type" attribute as deprecated erroneously [\#1034](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1034) - INTMDB-556
- Error "Pager Duty API key must consist of 32 hexadecimal digits" [\#1054](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1054) - INTMDB-570
- Terraform provider stuck in changes for [`advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) configuration [\#1066](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1066) - INTMDB-572
- [`search_index`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/search_index) does not recreate when `cluster_name` and `project_id` fields change [\#1078](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1078) - INTMDB-576
- POST Create Access List Entries for One Organization API Key endpoint supports list, but Terraform does not [\#1065](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1065) - INTMDB-579
- Typo in Readme [\#1073](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1073)
- Update project_api_key.html.markdown [\#1044](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1044)
- Doc Clean Up [`cloud_provider_access`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_provider_access) [\#1035](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1035)
- Update alert_configuration.html.markdown [\#1068](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1068)
- Doc Clean Up `provider_backup_enabled` deprecated to `cloud_backup` [\#1036](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1036)

**Closed Issues:**

- Unable to create third party integration of type Datadog with version 1.8.0 [\#1038](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1038)
- mongodbatlas\_third\_party\_integration - api\_token keeps updating on every apply [\#963](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/963)
- Unable to add mongodbatlas provider to CDK [\#952](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/952)
- \[Bug\] `update_snapshots` doesn't save at TF state with [`cloud_backup_schedule`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_backup_schedule) resource [\#904](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/904)
- Cannot ignore changes for replication\_specs when [`auto_scaling`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#region_configs) is enabled [\#888](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/888)
- Warning: Deprecated attribute [\#1042](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1042)
- GCP Network Peering remains pending when created via terraform [\#917](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/917)

**Internal Improvements:**

- Chore\(deps\): Bump github.com/hashicorp/hcl/v2 from 2.16.0 to 2.16.1 [\#1062](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1062) ([dependabot[bot]](https://github.com/apps/dependabot))
- Update access\_list\_api\_key.html.markdown [\#1058](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1058) ([Zuhairahmed](https://github.com/Zuhairahmed))
- Chore\(deps\): Bump github.com/hashicorp/hcl/v2 from 2.15.0 to 2.16.0 [\#1055](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1055) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.9 to 0.41.10 [\#1051](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1051) ([dependabot[bot]](https://github.com/apps/dependabot))
- Update CODEOWNERS to use APIx-Integrations [\#1050](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1050) ([andreaangiolillo](https://github.com/andreaangiolillo))
- Chore\(deps\): Bump golangci/golangci-lint-action from 3.3.1 to 3.4.0 [\#1026](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1026) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.8.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.8.0) (2023-1-26)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.7.0...v1.8.0)

**Enhancements:**

- Snapshot Distribution Support [\#979](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/979) - INTMDB-400
- Programmatically Create API Keys [\#974](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/974) - INTMDB-346
- Retrieve `org_id` from API Keys [\#973](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/973) - INTMDB-454
- `oplogMinRetentionHours` Parameter Support in `advanced_cluster` and `cluster` resources [\#1016](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1016) - INTMDB-397
- Analytics Node Tier New Features Support [\#994](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/994) - INTMDB-488
- Improve Default Alerts and Example Creation [\#993](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/993) - INTMDB-382
- Expand documentation for `cloud_backup_schedule` to include information about valid values for `frequency_interval` [\#1007](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1007) - INTMDB-547

**Depreciations:**

- `cloud_provider_snapshot`, `cloud_provider_snapshot_backup_policy`, `cloud_provider_snapshot_restore_job`, and `private_ip_mode` are now deprecated and will be removed from codebase as of v1.9 release [\#988](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/988) - INTMDB-409
- `NEW_RELIC` and `FLOWDOCK` in `third_party_integration` resource are now deprecated and will be removed from codebase as of v1.9 release [\#989](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/989) - INTMDB-482

**Bug Fixes:**

- Hide `current_certificate` when X.509 Authentication Database Users are Created [\#985](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/985) - INTMDB-468
- Import example added for `encryption_at_rest` resource [\#992](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/992) - INTMDB-530 
- Resource `cloud_backup_snapshot_export_job` variable name change [\#976](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/976) - INTMDB-523
- `update_snapshot` doesn't save at TF state with `cloud_backup_schedule` resource fix [\#974](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/974) - INTMDB-472
- Invitation handling after user accepts invitation fix [\#1012](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1012) - INTMDB-511 
- `bi_connector` settings in `advanced_cluster` fix (breaking changes) [\#1010](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1010) - INTMDB-455
- `third_party_integration` api_token keeps updating on every apply fix [\#1011](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1011) - INTMDB-519
- `custom_db_role` error fix [\#1009](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1009) - INTMDB-448
- `ldap_configuration` and `ldap_verify` resources fix [\#1004](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1004) - INTMDB-543
- `cloud_backup_schedule` resource fix [\#968](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/968) - INTMDB-427
- `search_index_test` fix [\#964](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/964) - INTMDB-341
- Cannot ignore changes for replication_specs when autoscaling enabled fix [\#961](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/961) - INTMDB-464
- BI Connector documentation fix [\#1017](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1017) 
- `federated_settings_org_config` import example fix [\#996](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/996)

**Closed Issues:**

- Documentation: Expand documentation for mongodbatlas\_cloud\_backup\_schedule to include information about valid values for frequency\_interval  [\#1005](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1005)
- Serverless instance returns incorrect connection string [\#934](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/934)
- Terraform apply failed with Error: Provider produced inconsistent final plan This is a bug in the provider, which should be reported in the provider's own issue tracker. [\#926](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/926)

**Internal Improvements:**

- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.7 to 0.41.9 [\#1013](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/1013) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.7.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.7.0) (2023-1-16)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.6.1...v1.7.0)

**Enhancements:**

- AWS Secrets Manager (AWS SM) Authetication for Terraform Atlas Provider [\#975](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/975) - INTMDB-521 

**Bug Fixes:**

- Resource cloud_backup_snapshot_export_job variable name change [#976](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/976) - INTMDB-523
- Deprecate legacy mongodbatlas.erb given Terraform Registry autogeneration [#962](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/962) - INTMDB-477  

**Closed Issues:**

- Terraform plan fail: Asymmetric hardware is not supported by the v1.0 API [\#958](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/958)
- Error importing ressource mongodbatlas\_network\_peering.mongo\_peer [\#906](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/906)
- \[Bug\] `container_id` is unconfigurable  attribute at `mongodbatlas_advanced_cluster` resource [\#890](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/890)
- mongodbatlas\_alert\_configuration - api\_token keeps wanting to change [\#863](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/863)
- Docs - Example - Return a Connection String - Azure Private Endpoint [\#713](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/713)

**Internal Improvements:**

- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.6 to 0.41.7 [\#978](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/978) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump actions/stale from 6 to 7 [\#977](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/977) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.4 to 0.41.6 [\#967](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/967) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/go-test/deep from 1.0.8 to 1.1.0 [\#966](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/966) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump goreleaser/goreleaser-action from 3 to 4 [\#965](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/965) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.6.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.6.1) (2022-12-6)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.6.0...v1.6.1)

**Enhancements:**

- Enable Adv Cluster and Cluster to have configurable timeouts [\#951](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/951) - INTMDB-503 
- Updated Prometheus Example [\#942](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/942) - INTMDB-498 
- Auto-Generate Changelog [\#944](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/944) - INTMDB-478 

**Bug Fixes:**

- Alert Configuration -- Api Token erroneous changes [\#941](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/941) - INTMDB-426
- Fix example private endpoint called out in issue 713 [\#907](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/907) - INTMDB-434 
- Cluster rename is inconsistently rejected by Terraform [\#929](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/929) - INTMDB-438 
- Terraform does not wait for cluster update when creating GCP private endpoints [\#943](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/943) - INTMDB-465
- mongodbatlas_federated_settings_connected_organization customer HELP + doc bug [\#924](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/924) - INTMDB-481  
- Serverless Private Endpoint Connection String Example Fix [\#940](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/940) - INTMDB-493 
- Fix regional mode endpoint test [\#946](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/946) - INTMDB-470 
- Skip tests for OPS GENIE and GOV [\#937](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/937) - INTMDB-484 
- Test Instability around privatelink tests fix [\#895](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/895) - INTMDB-384 
- Shorten test names that are too long to allow for targeting specific tests [\#932](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/932) - INTMDB-368 
- Remove container_id from configurable attribute in advanced_cluster [\#931](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/931) - INTMDB-463 

**Closed Issues:**

- No documented way to get config out of third party integration [\#939](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/939)
- Double checking Terraform Plan before I destroy Production [\#938](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/938)
- Issue: MongoDB Atlas Plugin Failure v 1.5.0 [\#928](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/928)
- MILLION\_RPU unit isn't supported by provider [\#854](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/854)
- MS Teams alert support in terraform provider is missing [\#827](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/827)

**Internal Improvements:**

- v1.6.1 - Conditionally ignore serverless connection string changes [\#953](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/953) ([evertsd](https://github.com/evertsd))
- Swap logic for variable substitution [\#950](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/950) ([martinstibbe](https://github.com/martinstibbe))
- Fix serverless endpoint tests [\#949](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/949) ([evertsd](https://github.com/evertsd))
- Release staging v1.6.1 [\#947](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/947) ([martinstibbe](https://github.com/martinstibbe))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.41.0 to 0.41.3 [\#936](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/936) ([dependabot[bot]](https://github.com/apps/dependabot))
- Serverless Endpoint Service Doc Bug [\#930](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/930) ([Zuhairahmed](https://github.com/Zuhairahmed))

## [v1.6.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.6.0) (2022-11-17)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.5.0...v1.6.0)

**Enhancements:** 

- Termination Protection for Advanced Cluster/Cluster/Serverless Instances [\#912](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/912) - INTMDB-444
- AWS/Azure Serverless Private Endpoints [\#913](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/913) - INTMDB-364

**Internal Improvements:**

- docs(website): fix federated_settings_org_config resource name by removing the misleading trailing s [\#908](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/908)
- chore(github): add link to contribution guidelines in PR template [\#910](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/910#issuecomment-1310007413)
- docs(resource/role_mapping): indent sub-elements of role_assignments [\#918](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/918)
- docs(resource/role_mapping): add link to reference of available role IDs [\#919](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/919)
- federated settings plural fix [\#914](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/914)  
- Chore(deps): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.24.0 to 2.24.1 [\#922](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/922)
- Chore(deps): Bump golangci/golangci-lint-action from 3.3.0 to 3.3.1 [\#925](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/925)
- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.40.24 to 0.41.0 [\#923](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/923)
- Workaround to handle serverless endpoint tests failing due to provider name missing from API [\#927](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/927)
- Release staging v1.6.0 [\#921](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/921)

## [v1.5.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.5.0) (2022-11-07)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.6...v1.5.0)

**Enhancements:** 
- INTMDB-224 - Support AtlasGov with Terraform [\#865](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/865)
- INTMDB-314 - Feature add: Add ability to upgrade shared/TENANT tiers for clusters and advanced clusters [\#874](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/874)
- INTMDB-349 - New AtlasGov parameter to tell the provider to use the Atlas gov base URL [\#865](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/865)
- INTMDB-373 - Add new notification parameters to the mongodbatlas_alert_config resource [\#883](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/883)	
- INTMDB-378 - Document for users how to get a pre-existing container id [\#883](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/883)
- INTMDB-377 - Release 1.5 (both pre and then GA) [\#887](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/887)
- INTMDB-394 - MS Teams alert support [\#320](https://github.com/mongodb/go-client-mongodb-atlas/pull/320)

**Bug Fixes:** 
- INTMDB-326 - Review code/tests and docs for resource_mongodbatlas_search_index [\#891](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/891)
- INTMDB-334 - Determine best path forward for GCP PSC timeouts and implement [\#859](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/859)
- INTMDB-370 - Docs not complete for cloud_backup_snapshot_restore_job	[\#870](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/870)
- INTMDB-403 - Update third_party_integration.markdown [\#851](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/851)
- INTMDB-405 - Add cluster label to advanced clusters	[\#857](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/857)
- INTMDB-406 - MILLION_RPU unit isn't supported by provider #854 [\#854](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/854)
  
**Closed Issues:**
- MS Teams alert support in terraform provider is missing [\#827](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/827)
- MILLION_RPU unit isn't supported by provider not_stale  [\#854](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/854)

**Merged Pull Requests:**
- INTMDB-378: Add link for How To Guide for existing container ID [\#883](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/883)
- INTMDB-403: Update third_party_integration.markdown [\#851](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/851)
- INTMDB-404: mongodbatlas_advanced_cluster doc updates [\#852](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/852)
- INTMD-428: doc update to "mongodbatlas_projects" [\#869](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/869)
- Docs: fix custom_dns_configuration_cluster_aws [\#860](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/860)
- Relying on atlas api for unit validation on alert configuration [\#862](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/862)
- Adding a github actions to automatically close stale issues/PRs based on CLOUDP-79100 [\#872](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/872)
- Encryption_at_rest M10+ limit doc update [\#886](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/886)
- Update cluster.html.markdown [\#878](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/878)
- Add APIx1 CODEOWNER [\#894](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/894)
- Chore(deps): Bump octokit/request-action from 2.1.6 to 2.1.7 [\#868](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/868)
- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.40.22 to 0.40.24 [\#875](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/875)
- Chore(deps): Bump golangci/golangci-lint-action from 3.2.0 to 3.3.0 [\#897](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/897)

## [v1.4.6](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.6) (2022-09-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.5...v1.4.6)

**Enhancements and Bug Fixes:** 
- INTMDB-387 - Enable Azure NVME for Atlas Dedicated clusters [\#833](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/833)
- INTMDB-342 - Update TestAccDataSourceMongoDBAtlasPrivateEndpointRegionalMode_basic test to use a new project to prevent conflicts  [\#837](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/837)
- INTMDB-347 - Cloud_backup is not being correctly imported - issue [\#768](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/768)
- INTMDB-354 - Update docs around what requires an API key access list [\#834](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/834)
- INTMDB-363 - [Updated Feature] Add serverless backup to mongodbatlas_serverless_instance [\#830](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/830)
- INTMDB-379 - Release 1.4.6 (both pre and then GA)	
- INTMDB-381 - Customer is unable to disable backup auto export [\#823](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/823)
- INTMDB-383 - Update the BYOK/encryption at rest resource [\#805](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/805)
- INTMDB-385 - use_org_and_group_names_in_export_prefix is not working for a customer
- INTMDB-386 - Add new role types to invitation verification	[\#840](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/840)
- INTMDB-371 - Timeout when creating privatelink_endpoint [\#806](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/806)
- INTMDB-372 - Fix failing test for testAccMongoDBAtlasAlertConfigurationConfigWithMatchers	[\#836](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/836)
- INTMDB-358 - Upgrade to go1.18 [\#835](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/835)
- INTMDB-391 - Doc Fix for teams.html.markdown [\#838](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/838)

**Closed Issues:**
-  importing existing cluster does not populate backup status #768 [\#768](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/768)


**Merged Pull Requests:**
- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.40.21 to 0.40.22 [\#842](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/842) ([dependabot[bot]](https://github.com/apps/dependabot))

- Rename team.html.markdown into teams.html.markdown [\#838](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/838) ([leo-ferlin-sutton](https://github.com/leo-ferlin-sutton))

- Chore(deps): Bump github.com/gruntwork-io/terratest from 0.40.20 to 0.40.21 [\#825](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/825) ([dependabot[bot]](https://github.com/apps/dependabot))

- Fix create index error msg[\#824](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/824) ([ebouther](https://github.com/ebouther))


## [v1.4.5](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.5)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.4...v1.4.5)

**Fixed**

- INTMDB-369: Fix parsing of `delivery_type_config` when using `point_in_time` for `cloud_backup_snapshot_restore_job`, in [\#813](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/813)
- INTMDB-322: Validated serverless alert_configurations and improved documentation on usage, addressing issue [\#722](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/722) in [\#819](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/819)

## [v1.4.4](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.4) (2022-08-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.3...v1.4.4)

**Fixed**

- INTMDB320 - Fix Global Cluster import documentation, in [\#796](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/796)
- INTMDB-331 - Update GCP documentation, issue [\#753](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/753), in [\#793](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/793)
- INTMDB-351 - Project data_source reads name, issue [\#788](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/788), in [\#795](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/795)
- INTMDB-362: Header Clarifications "Resource" vs "Data Source" in Documentation, in [\#803])(https://github.com/mongodb/terraform-provider-mongodbatlas/pull/803)
- INTMDB-343: Update go from 1.16 to 1.17 add windows arm64 build support, in [\#797](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/797)

## [v1.4.4-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.4-pre.1) (2022-08-17)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.3...v1.4.4-pre.1)

**Closed Issues:**

- Unable to update members in an existing "mongodbatlas\_teams" as the provider attempts to remove all users first [\#790](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/790)
- Please elaborate how to acquire PROJECTID and PEERINGID and PROVIDERNAME for import of network peering [\#789](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/789)
- error: error reading cloud provider access cloud provider access role not found in mongodbatlas, please create it first [\#781](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/781)

**Merged Pull Requests:**

- Update CONTRIBUTING.md [\#798](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/798) ([themantissa](https://github.com/themantissa))
- Fix federated\_settings\_identity\_provider attribute name [\#791](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/791) ([florenp](https://github.com/florenp))

## [v1.4.3](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.3) (2022-07-12)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.2...v1.4.3)

**Fixed:**

* INTMDB-335: Add option for multiple weekly monthly schedules @martinstibbe in [\#784](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/784)
* INTMDB-348: autoexport parameter not being set via provider [\#784](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/784)
* INTMDB-323: Removed the requirement to set `MONGODB_ATLAS_ENABLE_BETA` to use serverless and update the docs to match. [\#783](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/783)
* INTMDB-330 Fixed Serverless Instance Import Documentation. Closes [\#754](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/754)
* Fix typo in custom_db_role documentation [\#780](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/780)
* Fix typo in federated_settings_org_configs documentation [\#779](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/779)
 

## [v1.4.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.2) (2022-07-7)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.1...v1.4.2)

**Fixed:**

* INTMDB-313: Update project settings default flags by @martinstibbe in https://github.com/mongodb/terraform-provider-mongodbatlas/pull/778

## [v1.4.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.1) (2022-07-7)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.4.0...v1.4.1)

**Fixed:**

* Update CHANGELOG.md by @martinstibbe in https://github.com/mongodb/terraform-provider-mongodbatlas/pull/771
* INTMDB-313: Update project settings default flags by @martinstibbe in https://github.com/mongodb/terraform-provider-mongodbatlas/pull/773


## [v1.4.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.0) (2022-07-5)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1...v1.4.0) 

**Closed Issues:**

Note: the binary executable for windows/arm64 is not available for this release.  Next release will include.
- Fix for Add support for cloud export backup to mongodbatlas_cloud_backup_schedule [\#740](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/740)
- Feature Add: Update the project resource with new settings [\#741](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/741)
- Fix for  Potential bug when disabling auditing [\#705](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/705)
- Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings [\#706](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/706)
- Fix for Correct import function for snapshot export bucket #714 [\#715](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/715)
- Fix for Add support for schema migration [\#717](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/717)
- Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings
- Fix for Cannot import export bucket - bad state id encoding [\#708](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/708)
- Error missing expected { when updating the provider [\#697](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/697)

**Merged Pull Requests:**

- INTMDB-321: Add support for cloud export backup to mongodbatlas_cloud_backup_schedule [\#740](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/740) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-313: Update the project resource with new settings [\#741](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/741) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-301: Feature add: Add support for management of federated authentication configuration [\#742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/742) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-307: Add Regionalized Private Endpoint Settings [\#718](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/718) ([evertsd](https://github.com/evertsd))
- INTMDB-310: Potential bug when disabling auditing [\#705](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/705) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-311: Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings [\#706](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/706) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-315: Correct import function for snapshot export bucket [\#715](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/715) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-309: Add support for schema migration [\#717](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/717) ([martinstibbe](https://github.com/martinstibbe))

## [v1.4.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.4.0-pre.1) (2022-06-29)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1...v1.4.0-pre.1) 

**Closed Issues:**

Note: the binary executable for windows/arm64 is not available for this release.  Next release will include.
- Fix for Add support for cloud export backup to mongodbatlas_cloud_backup_schedule [\#740](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/740)
- Feature Add: Update the project resource with new settings [\#741](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/741)
- Fix for  Potential bug when disabling auditing [\#705](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/705)
- Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings [\#706](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/706)
- Fix for Correct import function for snapshot export bucket #714 [\#715](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/715)
- Fix for Add support for schema migration [\#717](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/717)
- Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings
- Fix for Cannot import export bucket - bad state id encoding [\#708](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/708)
- Error missing expected { when updating the provider [\#697](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/697)

**Merged Pull Requests:**

- INTMDB-321: Add support for cloud export backup to mongodbatlas_cloud_backup_schedule [\#740](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/740) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-313: Update the project resource with new settings [\#741](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/741) ([martinstibbe](https://github.com/martinstibbe)) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-301: Feature add: Add support for management of federated authentication configuration [\#742](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/742) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-307: Add Regionalized Private Endpoint Settings [\#718](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/718) ([evertsd](https://github.com/evertsd))
- INTMDB-310: Potential bug when disabling auditing [\#705](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/705) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-311: Feature Add: Prometheus and Microsoft Team to the Third Party Integration Settings [\#706](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/706) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-315: Correct import function for snapshot export bucket [\#715](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/715) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-309: Add support for schema migration [\#717](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/717) ([martinstibbe](https://github.com/martinstibbe))

## [v1.3.1-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1-pre.1) (2022-02-23)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.3.1...v1.3.1-pre.1)

**Closed issues:**
- Advance Cluster resource ignoring the autoscaling options [\#686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/686)
- Ensure we handle new flow for project deletion well #688  [\#688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/688) 
- Provider did not catch 400 error returned from the API [\#687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/687)
- Update timing of autodefer [\#695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/695)

**Merged pull requests:**

- INTMDB-300: Advance Cluster resource ignoring the autoscaling options [\#686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/686) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-302: Ensure we handle new flow for project deletion well #688  [\#688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/688) ([evertsd](https://github.com/evertsd))
-  INTMDB-303: Provider did not catch 400 error returned from the API [\#687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/687) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-305: Update timing of autodefer [\#695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/695) ([martinstibbe](https://github.com/martinstibbe))

## [v1.3.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1) (2022-03-28)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.3.1...v1.3.1-pre.1)

- INTMDB-306: [Terraform] Release bug fix version 1.3.1
## [v1.3.1-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.1-pre.1) (2022-02-23)

**Closed issues:**
- Advance Cluster resource ignoring the autoscaling options [\#686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/686)
- Ensure we handle new flow for project deletion well #688  [\#688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/688) 
- Provider did not catch 400 error returned from the API [\#687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/687)
- Update timing of autodefer [\#695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/695)

**Merged pull requests:**

- INTMDB-300: Advance Cluster resource ignoring the autoscaling options [\#686](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/686) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-302: Ensure we handle new flow for project deletion well #688  [\#688](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/688) ([evertsd](https://github.com/evertsd))
-  INTMDB-303: Provider did not catch 400 error returned from the API [\#687](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/687) ([martinstibbe](https://github.com/martinstibbe))
- INTMDB-305: Update timing of autodefer [\#695](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/695) ([martinstibbe](https://github.com/martinstibbe))

## [v1.3.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.0) (2022-02-23)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.3.0-pre.1...v1.3.0)

**Merged pull requests:**

- Create 1.3.0-upgrade-guide.html.markdown [\#682](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/682) ([themantissa](https://github.com/themantissa))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.40.2 to 0.40.3 [\#681](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/681) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.3.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.3.0-pre.1) (2022-02-22)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.2.0...v1.3.0-pre.1)

**Closed issues:**

- Auto scaling of storage cannot be disabled for mongodbatlas\_advanced\_cluster via "disk\_gb\_enabled" [\#677](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/677)
- Can't create M0 free tier on Azure [\#675](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/675)
- Error: error creating MongoDB Cluster: unexpected EOF [\#674](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/674)
- attempting to assign api key to resource `mongodbatlas_project` results in `Error: Unsupported block type` [\#671](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/671)
- Error in documentation page : aws\_kms -\> aws\_kms\_config [\#666](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/666)
- Accepting organization invitation causes 404 [\#636](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/636)
- Alert configuration state is not stable [\#632](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/632)
- mongodbatlas\_cloud\_backup\_schedule with Azure results in "restore window days" mandatory [\#625](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/625)
- incorrect values from `mongodbatlas_cluster` and `mongodbatlas_clusters` datasources [\#618](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/618)
- mongodb atlas network container atlas cidr block value and real  value is not mached [\#617](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/617)
- config\_full\_document\_before not working in EventTriggers [\#616](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/616)
- Hitting timing issue when trying to integrate with `aws` provider's `aws_iam_access_key` resource [\#127](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/127)

**Merged pull requests:**

- INTMDB-291: pre-release v1.3.0 [\#680](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/680) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.40.1 to 0.40.2 [\#679](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/679) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.40.0 to 0.40.1 [\#678](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/678) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-299: Support Cloud Backup Export Jobs [\#673](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/673) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump octokit/request-action from 2.1.0 to 2.1.4 [\#672](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/672) ([dependabot[bot]](https://github.com/apps/dependabot))
- update the documentations for new changes [\#670](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/670) ([nikhil-mongo](https://github.com/nikhil-mongo))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.9 to 0.40.0 [\#669](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/669) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-298: fixes a bug where you couldn't delete a team in team resource [\#668](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/668) ([coderGo93](https://github.com/coderGo93))
- INTMDB-297: set the container id to avoid null in state of data source cluster\(s\) [\#667](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/667) ([coderGo93](https://github.com/coderGo93))
- fixed typo in mongodbatlas\_teams sidebar [\#665](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/665) ([MartinCanovas](https://github.com/MartinCanovas))
- INTMDB-295: Fixes a bug about unauthorized error in project resource [\#664](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/664) ([coderGo93](https://github.com/coderGo93))
- INTMDB-293: Added container\_id in advanced cluster [\#663](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/663) ([coderGo93](https://github.com/coderGo93))
- INTMDB-294: \[Terraform\] Address security warnings from dependabot [\#661](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/661) ([thetonymaster](https://github.com/thetonymaster))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.8 to 0.38.9 [\#660](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/660) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-290: Added advanced configuration for datasource/resource of advanced cluster [\#658](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/658) ([coderGo93](https://github.com/coderGo93))
- Fix 1.2 upgrade/info guide formatting error [\#657](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/657) ([themantissa](https://github.com/themantissa))

## [v1.2.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.2.0) (2022-01-14)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.2.0-pre.1...v1.2.0)

**Merged pull requests:**

- INTMDB-268: Release v1.2.0 [\#656](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/656) ([abner-dou](https://github.com/abner-dou))

## [v1.2.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.2.0-pre.1) (2022-01-13)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.1.1...v1.2.0-pre.1)

**Closed issues:**

- mongodbatlas\_teams provides a Team resource [\#649](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/649)
- terraform-provider-mongodbatlas\_v0.8.0 plugin: panic: runtime error: invalid memory address or nil pointer dereference [\#644](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/644)
- backup snapshot and restore not working automated  [\#642](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/642)
- Delete default alerts [\#628](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/628)
- Cannot set multiple notifications for an alert [\#626](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/626)
- PRIVATE\_ENDPOINT\_SERVICE\_ALREADY\_EXISTS\_FOR\_REGION [\#590](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/590)
- INVALID\_CLUSTER\_CONFIGURATION when modifying a cluster to use replication\_specs \(eg for multi-region\) [\#588](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/588)

**Merged pull requests:**

- Re-branch from 651 due to conflict of docs fixes [\#654](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/654) ([themantissa](https://github.com/themantissa))
- INTMDB-268: Pre-release v1.2.0 [\#650](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/650) ([abner-dou](https://github.com/abner-dou))
- INTMDB-5: added parameter team name for alert configurations  [\#648](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/648) ([coderGo93](https://github.com/coderGo93))
- Fix markdown formatting in network\_container.html.markdown [\#647](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/647) ([pmacey](https://github.com/pmacey))
- INTMDB-15: Added parameter advanced conf for cluster datasource [\#646](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/646) ([coderGo93](https://github.com/coderGo93))
- INTMDB-284: Updated docs in alert configuration resource and datasource [\#645](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/645) ([coderGo93](https://github.com/coderGo93))
- INTMDB-285: Fix org\_invitations issue [\#643](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/643) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.9.0 to 2.10.1 [\#641](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/641) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-263: Create Resource and Datasource for private\_link\_endpoint\_adl [\#640](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/640) ([abner-dou](https://github.com/abner-dou))
- INTMDB-287: Fixes the issues in project api keys [\#639](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/639) ([coderGo93](https://github.com/coderGo93))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.5 to 0.38.8 [\#638](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/638) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-283: Fixes a bug about optional parameters in Cloud Backup Schedule [\#631](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/631) ([coderGo93](https://github.com/coderGo93))
- INTMDB-281: Fix realm event trigger issue [\#630](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/630) ([abner-dou](https://github.com/abner-dou))
- INTMDB-282: Updated test and docs for alert configuration using notifications [\#629](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/629) ([coderGo93](https://github.com/coderGo93))
- INTMDB-280: Fix cluster datasource scaling issue [\#627](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/627) ([abner-dou](https://github.com/abner-dou))
- INTMDB-272: Validate using interval\_min for PagerDuty, VictorOps, GenieOps [\#624](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/624) ([coderGo93](https://github.com/coderGo93))
- INTMDB-276: Added VersionReleaseSystem parameter for resource/datasource\(s\) of cluster [\#623](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/623) ([coderGo93](https://github.com/coderGo93))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.4 to 0.38.5 [\#622](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/622) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.8.0 to 2.9.0 [\#621](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/621) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-24: Change computed paused to optional [\#620](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/620) ([coderGo93](https://github.com/coderGo93))
- INTMDB-257: Changed 'hcl' markdown tag to 'terraform' tag [\#619](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/619) ([abner-dou](https://github.com/abner-dou))
- Fix link in update guide and add version [\#615](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/615) ([themantissa](https://github.com/themantissa))
- INTMDB-279: Fixes a bug where it crashes when importing a trigger [\#614](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/614) ([coderGo93](https://github.com/coderGo93))
- mongodbatlas-project resource: add api\_keys attribute [\#504](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/504) ([tetianakravchenko](https://github.com/tetianakravchenko))

## [v1.1.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.1.1) (2021-11-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.1.0...v1.1.1)

**Closed issues:**

- Cannot update the default backup schedule policy without defining API Key access IPs. [\#610](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/610)

**Merged pull requests:**

- Release v1.1.1 [\#613](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/613) ([abner-dou](https://github.com/abner-dou))
- Fix documentation v1.1.0 [\#612](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/612) ([abner-dou](https://github.com/abner-dou))

## [v1.1.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.1.0) (2021-11-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.1.0-pre.1...v1.1.0)

**Merged pull requests:**

- INTMDB-264: Release v1.1.0 [\#611](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/611) ([abner-dou](https://github.com/abner-dou))
- Guide and minor main page changes for 1.1.0 [\#609](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/609) ([themantissa](https://github.com/themantissa))

## [v1.1.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.1.0-pre.1) (2021-11-17)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.2...v1.1.0-pre.1)

**Fixed bugs:**

- Plugin crash when changing value of iam\_assumed\_role\_arn for resource mongodbatlas\_cloud\_provider\_access\_authorization [\#565](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/565)
- terraform-provider-mongodbatlas\_v1.0.1 crashes after creating mongodbatlas\_cloud\_provider\_access\_authorization resource [\#554](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/554)
- INVALID\_CLUSTER\_CONFIGURATION when adding new regions [\#550](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/550)

**Closed issues:**

- update\_snapshots is not getting set on Atlas while using "mongodbatlas\_cloud\_backup\_schedule" tf resource  [\#594](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/594)
- Race condition when destroying cluster and disabling encryption at rest on the project-level [\#518](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/518)
- Scaling max/min is applied each time with disabled autoscaling [\#482](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/482)
- Documentation: Update contribution readme for developing the provider for terraform +14 newer [\#466](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/466)

**Merged pull requests:**

- INTMDB-264: pre-release v1.1.0 [\#608](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/608) ([abner-dou](https://github.com/abner-dou))
- INTMDB:273: Fix replication\_specs update error [\#607](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/607) ([abner-dou](https://github.com/abner-dou))
- Update cloud\_provider\_snapshots.html.markdown [\#605](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/605) ([nhuray](https://github.com/nhuray))
- Fix docs for third party data source. [\#604](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/604) ([jkodroff](https://github.com/jkodroff))
- Fix timeout in acctest [\#602](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/602) ([abner-dou](https://github.com/abner-dou))
- INTMDB-270: fix issue with project resource importer test [\#601](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/601) ([abner-dou](https://github.com/abner-dou))
- Update MDB version info [\#600](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/600) ([themantissa](https://github.com/themantissa))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.2 to 0.38.4 [\#599](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/599) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-271: Fixing a bug and improving for custom zone mappings in global cluster config [\#597](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/597) ([coderGo93](https://github.com/coderGo93))
- INTMDB-275: Changed the pointer in some paremeters for custom db role [\#596](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/596) ([coderGo93](https://github.com/coderGo93))
- INTMDB-270: Added  'with\_default\_alerts\_settings' to project resource [\#595](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/595) ([abner-dou](https://github.com/abner-dou))
- Fix backup option from provider\_backup\_enabled to cloud\_backup [\#592](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/592) ([paikend](https://github.com/paikend))
- INTMDB-222: Added Synonyms to Search Index RS and DS [\#591](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/591) ([abner-dou](https://github.com/abner-dou))
- Fix typo: mongodbatlast =\> mongodbatlas [\#589](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/589) ([sonlexqt](https://github.com/sonlexqt))
- Chore\(deps\): Bump github.com/go-test/deep from 1.0.7 to 1.0.8 [\#587](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/587) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-260: Added GCP feature for Private Endpoint [\#586](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/586) ([coderGo93](https://github.com/coderGo93))
- INTMDB-227:Create new Resource and Datasource for Serverless Instance [\#585](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/585) ([abner-dou](https://github.com/abner-dou))
- INTMDB-269: Fix issue with default auto\_scaling\_disk\_gb\_enabled value [\#584](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/584) ([abner-dou](https://github.com/abner-dou))
- fixes failing snapshots because timeout is too short [\#583](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/583) ([ebostijancic](https://github.com/ebostijancic))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.38.1 to 0.38.2 [\#582](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/582) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-239: Added new resource/datasource and deprecate for cloud backup snapshot and restore job [\#581](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/581) ([coderGo93](https://github.com/coderGo93))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.12 to 0.38.1 [\#580](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/580) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-216: Added resource and datasource\(s\) of Advanced Cluster [\#570](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/570) ([coderGo93](https://github.com/coderGo93))
- Add Organisation and Project invitations [\#560](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/560) ([beergeek](https://github.com/beergeek))

## [v1.0.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.2) (2021-10-07)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.2-pre.1...v1.0.2)

**Closed issues:**

- gcp cluster doc issue [\#568](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/568)
- mongodbatlas\_auditing documentation mismatch [\#555](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/555)

**Merged pull requests:**

- INTMDB-246: Release v1.0.2 [\#579](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/579) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.11 to 0.37.12 [\#578](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/578) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.0.2-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.2-pre.1) (2021-10-04)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.1...v1.0.2-pre.1)

**Fixed bugs:**

- Error: error getting search index information: json: cannot unmarshal array into Go struct field IndexMapping.mappings.fields of type mongodbatlas.IndexField [\#545](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/545)
- Regression: Data Source for mongodbatlas\_cluster makes terraform hang indefinitely using version 1.0 [\#521](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/521)
- Enabling encryption at rest with any provider for a Cluster will throw error [\#517](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/517)
- realm trigger causes provider to authenticate with atlas even if realm triggers are not in use [\#512](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/512)
- Adding IP to access List failed when lot of entries [\#470](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/470)

**Closed issues:**

- Datalake configuration at creation time [\#561](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/561)
- Nested map variable works if defined in module consumer but not if defined in module itself [\#559](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/559)
- Getting blocked by IP when trying to create a project / cluster [\#557](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/557)

**Merged pull requests:**

- INTMDB-246: pre-release 1.0.2-pre.1 [\#577](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/577) ([abner-dou](https://github.com/abner-dou))
- Roll up for documentation fixes [\#576](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/576) ([themantissa](https://github.com/themantissa))
- INTMDB-259: Fix issue when create a tenant cluster without auto\_scaling\_disk\_gb [\#575](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/575) ([abner-dou](https://github.com/abner-dou))
- INTMDB-203: Fix IOPS restriction on NVME clusters [\#574](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/574) ([abner-dou](https://github.com/abner-dou))
- INTMDB-254: Fix replication\_specs behaviour when update cluster [\#573](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/573) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.10 to 0.37.11 [\#572](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/572) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.7.1 to 2.8.0 [\#571](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/571) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.8 to 0.37.10 [\#569](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/569) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-256: Fixes a bug for updated a role in cloud access authorization [\#567](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/567) ([coderGo93](https://github.com/coderGo93))
- INTMDB-245: Added an example for encryption at rest using azure with a cluster [\#566](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/566) ([coderGo93](https://github.com/coderGo93))
- INTMDB-221: Added projectOwnerID to project resource [\#564](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/564) ([abner-dou](https://github.com/abner-dou))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.37.7 to 0.37.8 [\#563](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/563) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-211: Add new advanced shard key options in global cluster resource [\#562](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/562) ([abner-dou](https://github.com/abner-dou))
- INTMDB-252: Added two parameters for cluster advanced configuration [\#558](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/558) ([coderGo93](https://github.com/coderGo93))
- Fix typo in import search index error [\#556](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/556) ([stefanosala](https://github.com/stefanosala))
- INTMDB-230: added  property to maintenance window rs ds [\#552](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/552) ([abner-dou](https://github.com/abner-dou))
- INTMDB-249: Lint and fix linting for examples [\#538](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/538) ([gssbzn](https://github.com/gssbzn))

## [v1.0.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.1) (2021-09-02)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.1-pre.1...v1.0.1)

**Merged pull requests:**

-  tag version 1.0.1 for release [\#553](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/553) ([coderGo93](https://github.com/coderGo93))

## [v1.0.1-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.1-pre.1) (2021-09-01)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.0...v1.0.1-pre.1)

**Fixed bugs:**

- Cannot define a mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy and enable provider\_backup\_enabled for an existing cluster in the same apply [\#350](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/350)

**Closed issues:**

- Not able to obtain CSRS hostnames from mongodbatlas\_cluster resource [\#543](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/543)
- Cloud Provider Access Setup for Datalake [\#486](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/486)
- Terraform: Unable to fetch connection strings when using 'data' resource for existing cluster [\#422](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/422)
- SDK framework update v2.0.0+ [\#408](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/408)
- Issues with resource and API design for Cloud Provider Snapshot Backup Policy [\#222](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/222)

**Merged pull requests:**

- Chore\(deps\): Bump go.mongodb.org/realm from 0.0.1 to 0.1.0 [\#551](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/551) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.7.0 to 2.7.1 [\#549](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/549) ([dependabot[bot]](https://github.com/apps/dependabot))
- INTMDB-251: Update search rs and ds to use go-client v0.12.0 [\#548](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/548) ([abner-dou](https://github.com/abner-dou))
- tag version 1.0.1 for pre release [\#546](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/546) ([coderGo93](https://github.com/coderGo93))
- test: skip instead of fatal for team ids missing [\#544](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/544) ([gssbzn](https://github.com/gssbzn))
- Chore\(deps\): Bump github.com/spf13/cast from 1.3.1 to 1.4.1 [\#542](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/542) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/gruntwork-io/terratest from 0.32.20 to 0.37.7 [\#541](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/541) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump github.com/mongodb-forks/digest from 1.0.1 to 1.0.3 [\#540](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/540) ([dependabot[bot]](https://github.com/apps/dependabot))
- Chore\(deps\): Bump octokit/request-action from 2.0.0 to 2.1.0 [\#539](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/539) ([dependabot[bot]](https://github.com/apps/dependabot))
- feat: add dependabot [\#537](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/537) ([gssbzn](https://github.com/gssbzn))
- feat: mcli integration [\#536](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/536) ([gssbzn](https://github.com/gssbzn))
- docs: fix typo cluster.html.markdown [\#535](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/535) ([gssbzn](https://github.com/gssbzn))
- docs: update README [\#534](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/534) ([gssbzn](https://github.com/gssbzn))
- INTMDB-226 - Added forceNew to vpc\_id in network\_peering [\#533](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/533) ([abner-dou](https://github.com/abner-dou))
- Modified workflow to trigger the automated tests [\#532](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/532) ([coderGo93](https://github.com/coderGo93))
- INTMDB-247: Fixes a bug where it's taking 3 minutes to read a cluster [\#530](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/530) ([coderGo93](https://github.com/coderGo93))
- task: check examples are formatted correctly [\#529](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/529) ([gssbzn](https://github.com/gssbzn))
- feat: use golangci lint action [\#528](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/528) ([gssbzn](https://github.com/gssbzn))
- task: remove misspell as a dependency [\#527](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/527) ([gssbzn](https://github.com/gssbzn))
- INTMDB-236: Updated the cluster configuration [\#526](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/526) ([coderGo93](https://github.com/coderGo93))
- INTMDB-244: add deprecation notes for cloud backup documentation [\#525](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/525) ([abner-dou](https://github.com/abner-dou))
- INTMDB-237: fix word in private\_endpoint resource documentation [\#523](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/523) ([abner-dou](https://github.com/abner-dou))
- INTMDB-243: Fixes a bug for encryption at rest with new parameters [\#522](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/522) ([coderGo93](https://github.com/coderGo93))
- INTMDB-235: Added example of ldap configuration docs [\#520](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/520) ([coderGo93](https://github.com/coderGo93))
- INTMDB-242: Fixes the bug when if you don set public/private key it would fail for getting realm client [\#519](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/519) ([coderGo93](https://github.com/coderGo93))
- Add stronger warning against attempting a shared tier upgrade [\#516](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/516) ([themantissa](https://github.com/themantissa))
- INTMDB-219: Fixed cluster scaling issue [\#515](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/515) ([abner-dou](https://github.com/abner-dou))
- INTMDB-218: fixes the bug when you try to add more than 100 ip whitelist [\#514](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/514) ([coderGo93](https://github.com/coderGo93))
- INTMDB-225: Fixed network peering resource for Azure [\#513](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/513) ([abner-dou](https://github.com/abner-dou))

## [v1.0.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.0) (2021-08-11)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v1.0.0-pre.1...v1.0.0)

**Merged pull requests:**

-  INTMDB-215: tag version 1.0.0 for release [\#511](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/511) ([coderGo93](https://github.com/coderGo93))
- Update 1.0.0-upgrade-guide.html.markdown [\#510](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/510) ([themantissa](https://github.com/themantissa))
- update the privatelink doc with Azure example [\#509](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/509) ([nikhil-mongo](https://github.com/nikhil-mongo))

## [v1.0.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.0.0-pre.1) (2021-08-10)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.9.1...v1.0.0-pre.1)

**Closed issues:**

- Multi cloud not supported? [\#497](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/497)
- Unsupported argument `bi_connector_config`  [\#491](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/491)
- Support for Mongo DB Cluster 4.4 [\#487](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/487)
- Backup policy ID requirement is a catch 22 [\#485](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/485)
- Updating from terraform 0.14.5 tp 0.15.0 or further version \(up to 1.0\) mongodbatlas started to add database\_name="" outsides roles part [\#480](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/480)
- resource mongodbatlas\_auditing audit\_filter param doesn't ignore whitespace changes in json [\#477](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/477)
- documentation: Azure private link example [\#469](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/469)
- Terraform not able to detect all the changes from the mongodb .tf files [\#465](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/465)
- mongodbatlas\_private\_endpoint and mongodbatlas\_private\_endpoint\_interface\_link gets re-created everytime [\#464](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/464)
- mongodbatlas\_database\_user lifecycle ignore\_changes is ignored [\#462](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/462)
- Unable to manage LDAP groups due to forced incorrect auth\_database\_name [\#447](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/447)
- Request for Docs - Migration from Ahkaryle-s provider [\#26](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/26)

**Merged pull requests:**

- INTMDB-215: tag version 1.0.0 for pre-release [\#508](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/508) ([coderGo93](https://github.com/coderGo93))
- INTDMB-223: Updated Cloud Backup to SDK v2 [\#507](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/507) ([abner-dou](https://github.com/abner-dou))
- INTMDB-233: Update linter version [\#506](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/506) ([MihaiBojin](https://github.com/MihaiBojin))
- INTMDB-232: Fix user agent version [\#505](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/505) ([MihaiBojin](https://github.com/MihaiBojin))
- update resources documentation to address INTMDB-225 [\#503](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/503) ([nikhil-mongo](https://github.com/nikhil-mongo))
- INTMDB-202: Changed to TypeSet for replication specs [\#502](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/502) ([coderGo93](https://github.com/coderGo93))
- INTDMB-223: update search index to sdk v2 [\#501](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/501) ([abner-dou](https://github.com/abner-dou))
- INTMDB-17: fixed import state method in search index resource [\#500](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/500) ([abner-dou](https://github.com/abner-dou))
- Adding autodefer parameter to automatically defer any maintenance [\#499](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/499) ([vgarcia-te](https://github.com/vgarcia-te))
- docs: fix typo in GCP network\_peering.network\_name [\#498](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/498) ([daquinoaldo](https://github.com/daquinoaldo))
- INTMDB-180: file env variable spelling error [\#495](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/495) ([abner-dou](https://github.com/abner-dou))
- INTMDB-188: fixed issue related with read non-existing resource [\#494](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/494) ([abner-dou](https://github.com/abner-dou))
- add example for atlas-aws vpc peering [\#493](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/493) ([nikhil-mongo](https://github.com/nikhil-mongo))
- MongoDB Atlas-GCP VPC Peering [\#492](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/492) ([nikhil-mongo](https://github.com/nikhil-mongo))
- MongoDB Atlas - GCP VPC Peering example [\#490](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/490) ([nikhil-mongo](https://github.com/nikhil-mongo))
- INTMDB-183: Migrate to TF SDK 2 [\#489](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/489) ([coderGo93](https://github.com/coderGo93))
- INTMDB-17:  Resource/Data Source Atlas Search  [\#488](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/488) ([abner-dou](https://github.com/abner-dou))
- INTMDB-214: Deprecation of private endpoint [\#484](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/484) ([coderGo93](https://github.com/coderGo93))
- INTMDB-179: added more examples for connection strings [\#483](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/483) ([coderGo93](https://github.com/coderGo93))
- INTMDB-198: Fixes a bug where it appears empty private endpoint in cluster [\#481](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/481) ([coderGo93](https://github.com/coderGo93))
- INTMDB-201: Added to detect changes for name of cluster in update func [\#479](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/479) ([coderGo93](https://github.com/coderGo93))
- INTMDB-28: Added Event Triggers Realm [\#476](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/476) ([coderGo93](https://github.com/coderGo93))
- INTMDB-145: Cloud backup schedule  [\#475](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/475) ([coderGo93](https://github.com/coderGo93))
- Starter example improvements and doc update [\#474](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/474) ([themantissa](https://github.com/themantissa))
- INTMDB-212: Deprecation of Project IP Whitelist [\#473](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/473) ([coderGo93](https://github.com/coderGo93))
- INTMDB-18-Test for Online Archive and sync attribute for discussion [\#472](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/472) ([leofigy](https://github.com/leofigy))
- INTMDB-128: Modified design when you can get .id from various resources [\#471](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/471) ([coderGo93](https://github.com/coderGo93))
- update README - added plugin dev override [\#468](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/468) ([abner-dou](https://github.com/abner-dou))
- CLOUDP-90710: Expose BASE\_URL so that we can test terraform with a custom server [\#467](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/467) ([andreaangiolillo](https://github.com/andreaangiolillo))
- INTMDB-200: Fixes a bug about updating a region name with GCP [\#463](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/463) ([coderGo93](https://github.com/coderGo93))
- INTMDB-19: Added resource and datasource\(s\) of data lake [\#414](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/414) ([coderGo93](https://github.com/coderGo93))
- INTMDB-18 : DataSource and Resource support for Online Archive [\#413](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/413) ([leofigy](https://github.com/leofigy))

## [v0.9.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.9.1) (2021-05-17)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.9.1-pre.1...v0.9.1)

**Merged pull requests:**

- chore v0.9.1 changelog update [\#461](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/461) ([leofigy](https://github.com/leofigy))
- Missing formatting backtick in the documentation [\#457](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/457) ([lescactus](https://github.com/lescactus))

## [v0.9.1-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.9.1-pre.1) (2021-05-14)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.9.0...v0.9.1-pre.1)

**Fixed bugs:**

- mongodbatlas\_cluster bi\_connector state changes on terraform CLI 0.14.2 even without any bi\_connector configuration - terraform 14  [\#364](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/364)
- Update the CA certificate with a os environment [\#442](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/442) ([pitakill](https://github.com/pitakill))

**Closed issues:**

- New single apply cloud provider access requires encryption\_at\_rest\_provider set in mongodbatlas\_cluster [\#452](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/452)
- Migration to mongodbatlas\_cloud\_provider\_access\_setup / authorization [\#451](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/451)
- GCP can't set region for cluster [\#450](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/450)
- Error verifying GPG signature for provider "mongodbatlas" [\#448](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/448)
- "CLUSTER\_DISK\_IOPS\_INVALID" related error/unexpected update-in-place [\#439](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/439)
- encryption\_at\_rest failing with UNEXPECTED ERROR \(and discussion of Cloud Provider Access possible improvement\) [\#409](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/409)
- Test update - Update test certificate  [\#407](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/407)
- mongodbatlas\_private\_endpoint and mongodbatlas\_private\_endpoint\_interface\_link not working as expected in version 0.7 [\#406](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/406)

**Merged pull requests:**

- INTMDB-207 chore: Doc update for changelog v0.9.1-pre.1 [\#460](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/460) ([leofigy](https://github.com/leofigy))
- INTMDB-208: Updated docs for upgrading private endpoints [\#458](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/458) ([coderGo93](https://github.com/coderGo93))
- INTMDB-205 fixing client update side effects [\#456](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/456) ([leofigy](https://github.com/leofigy))
-  INTMDB-205-client-update bumping the client version up [\#455](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/455) ([leofigy](https://github.com/leofigy))
- Test config update [\#454](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/454) ([leofigy](https://github.com/leofigy))
- INTMDB-206 Documentation and example updates [\#453](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/453) ([leofigy](https://github.com/leofigy))
- updated cluster doc and examples  for the new IOPS change [\#446](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/446) ([nikhil-mongo](https://github.com/nikhil-mongo))
- fix page title and sidebar [\#445](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/445) ([themantissa](https://github.com/themantissa))
- chore v0.9.0 changelog [\#444](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/444) ([leofigy](https://github.com/leofigy))

## [v0.9.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.9.0) (2021-04-22)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.9.0-pre.1...v0.9.0)

**Implemented enhancements:**

- Test or TestAccResourceMongoDBAtlasDatabaseUser\_withAWSIAMType with [\#432](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/432) ([pitakill](https://github.com/pitakill))
- INTMDB 186 - Added authorization resource to split the cloud access provider config [\#420](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/420) ([leofigy](https://github.com/leofigy))

**Closed issues:**

- Outdated usage example about "mongodbatlas\_encryption\_at\_rest" [\#424](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/424)

**Merged pull requests:**

- Remove IOPS and adjust parameter description [\#443](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/443) ([themantissa](https://github.com/themantissa))
- remove unnecessary variables and use roles instead of keys for AWS Encryption-AtRest [\#441](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/441) ([nikhil-mongo](https://github.com/nikhil-mongo))
- Update default IOPS and 0.9.0 info guide [\#440](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/440) ([themantissa](https://github.com/themantissa))
- \[Azure VNET Peering\] changed the incorrect parameter used for role assignment and role definition [\#438](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/438) ([nikhil-mongo](https://github.com/nikhil-mongo))
- chore changelog for v0.9.0 prerelease [\#437](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/437) ([leofigy](https://github.com/leofigy))
- Update release.yml [\#436](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/436) ([leofigy](https://github.com/leofigy))
- INTMDB-199: Fixes the error when updating an replication specs after removed one zone [\#434](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/434) ([coderGo93](https://github.com/coderGo93))
- Examples of terratest upgrade [\#431](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/431) ([coderGo93](https://github.com/coderGo93))
- Fix: small doc error [\#428](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/428) ([lescactus](https://github.com/lescactus))
- INTMDB-194: Added func to get db major version for testing [\#427](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/427) ([coderGo93](https://github.com/coderGo93))
- Add  examples creating user with aws\_iam\_type [\#426](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/426) ([alexsergeyev](https://github.com/alexsergeyev))
- INTMDB-155: Fixes a bug related to bi\_connector cluster by deprecating [\#423](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/423) ([coderGo93](https://github.com/coderGo93))
- INTMDB-168: updated docs the format of using dependencies [\#421](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/421) ([coderGo93](https://github.com/coderGo93))
- INTMDB-185: Added parameter regions for GCP network container [\#418](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/418) ([coderGo93](https://github.com/coderGo93))
- TeamsUpdate - fixing small bug, again missing update [\#417](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/417) ([leofigy](https://github.com/leofigy))
- Fixes test about ca certificate x509 [\#416](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/416) ([coderGo93](https://github.com/coderGo93))
- Working example for Atlas-encryptionAtRest-roles with a single tf apply [\#415](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/415) ([zohar-mongo](https://github.com/zohar-mongo))
- INTMDB-181: Detects unnecessary changes changes for azure/gcp encryption at rest [\#412](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/412) ([coderGo93](https://github.com/coderGo93))
- corrected the title by making the variable name plural [\#404](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/404) ([crosbymichael1](https://github.com/crosbymichael1))
- INTMDB-154: Deprecation for provider\_encrypt\_ebs\_volume  [\#403](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/403) ([leofigy](https://github.com/leofigy))
- INTMDB-133: Vendor removal to include terratest samples [\#395](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/395) ([leofigy](https://github.com/leofigy))
- INTMDB-114/115: Added resource, datasource and tests for LDAP configuration and verify [\#379](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/379) ([coderGo93](https://github.com/coderGo93))
- INTMDB-116: Added parameter ldap auth type for resource and datasource\(s\) of database user [\#376](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/376) ([coderGo93](https://github.com/coderGo93))
- INTMDB-16: Added resource and datasource for Custom DNS Configuration for Atlas Clusters on AWS [\#370](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/370) ([coderGo93](https://github.com/coderGo93))
- INTMDB-133: Examples for encryption at rest with roles  [\#369](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/369) ([coderGo93](https://github.com/coderGo93))

## [v0.9.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.9.0-pre.1) (2021-04-21)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.2...v0.9.0-pre.1)

**Closed issues:**

- TF support for creating api keys [\#433](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/433)
- Apple Silicon \(darwin/arm64\) support [\#430](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/430)
- Tenant Provider Acceptance tests are failing [\#419](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/419)
- 500 \(request "UNEXPECTED\_ERROR"\) [\#411](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/411)
- Error Creating cluster GCP - 500 UNEXPECTED\_ERROR [\#410](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/410)
- r/mongodbatlas\_third\_party\_integration fails on read after create [\#354](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/354)

## [v0.8.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.2) (2021-02-03)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.1...v0.8.2)

**Closed issues:**

- Issues with AWS/Azure Privatelink in v0.8.0 and v0.8.1 [\#401](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/401)
- mongodbatlas-0.4.2: panic: runtime error: invalid memory address or nil pointer dereference [\#399](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/399)
- Seemingly unable to set provider source to mongodb/mongodbatlas for terraform v0.14? [\#396](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/396)
- Missing connection\_strings for Azure Private Link in resource mongodbatlas\_cluster [\#390](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/390)
- Error in Docs [\#387](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/387)

**Merged pull requests:**

- INTMDB-177: chore: release changelog [\#402](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/402) ([leofigy](https://github.com/leofigy))
- INTMDB-174: updated an example for cluster [\#400](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/400) ([coderGo93](https://github.com/coderGo93))
- INTMDB-175: Added azure status and fixes the error about target state [\#398](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/398) ([coderGo93](https://github.com/coderGo93))

## [v0.8.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.1) (2021-01-28)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.0...v0.8.1)

**Fixed bugs:**

- Removal of user scopes is not detected by the provider [\#363](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/363)

**Closed issues:**

- Parameter mismatch in mongodbatlas\_privatelink\_endpoint\_service [\#391](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/391)
- Can't add a team to a project [\#389](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/389)
- failed to create policy items while using mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy [\#386](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/386)
- Unable to import resources with "terraform import" - 401 \(request "Unauthorized"\) [\#385](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/385)

**Merged pull requests:**

-  INTMDB-172: chore changelog update for v0.8.1 [\#397](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/397) ([leofigy](https://github.com/leofigy))
- INTMDB-169: delete encoding url path in private endpoint service [\#393](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/393) ([coderGo93](https://github.com/coderGo93))
- INTMDB-158: Added private endpoint in connection strings [\#392](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/392) ([coderGo93](https://github.com/coderGo93))
- INTMDB-163: Wrong order for PrivateLink Endpoint Service and detects unnecessary changes [\#388](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/388) ([coderGo93](https://github.com/coderGo93))

## [v0.8.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.0) (2021-01-20)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.0-pre.2...v0.8.0)

**Closed issues:**

- Managing encryption at rest using iam roles fails [\#382](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/382)
- Unable to use mongodbatlas\_network\_peering data source [\#377](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/377)

**Merged pull requests:**

- INTMDB-153: Create 0.8.0-upgrade-guide.html.markdown [\#384](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/384) ([themantissa](https://github.com/themantissa))

## [v0.8.0-pre.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.0-pre.2) (2021-01-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.8.0-pre.1...v0.8.0-pre.2)

**Closed issues:**

- Bad Release Practice [\#381](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/381)

**Merged pull requests:**

- INTMDB-162: Fixes bug about detecting changes and make sensitive values [\#383](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/383) ([coderGo93](https://github.com/coderGo93))

## [v0.8.0-pre.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.8.0-pre.1) (2021-01-15)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.7.0...v0.8.0-pre.1)

**Fixed bugs:**

- Unexpected behaviour from resource `mongodbatlas_teams` when adding username for user not yet part of/Pending to join Organisation [\#329](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/329)
- mongodbatlas\_alert\_configuration - metric\_threshold.threshold is not being passed when value is zero [\#311](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/311)
- After manually deletion of a billing alert, no more plan oder apply will succeed, becuase of an 404 during plan [\#305](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/305)

**Closed issues:**

- Deleting a user from mongo atlas results in a 404 error, not that user getting re-created [\#360](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/360)
- Replace "mongodbatlas\_project\_ip\_whitelist" resource/datasource/docs references with "mongodbatlast\_project\_ip\_accesslist" to reflect API/UI change. [\#358](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/358)
- Quick start for provider is not quick and comes with side effect about `replication_factor` field [\#356](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/356)
- mongodbatlas\_database\_user resource's id attribute does not have the username value [\#348](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/348)
- Mongodbatlas documentation issue with Data Sources [\#347](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/347)
- Please add support for Azure Private Link as a private endpoint [\#346](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/346)
- mongodbatlas\_maintenance\_window fails with BadRequest: Invalid Day of Week [\#289](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/289)

**Merged pull requests:**

- INTMDB-160: Resetting an encryption at rest [\#380](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/380) ([coderGo93](https://github.com/coderGo93))
- INTMDB-149: tag version 0.8.0 for pre-release [\#378](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/378) ([coderGo93](https://github.com/coderGo93))
- Fix typo "requirments" in the PR template [\#375](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/375) ([KeisukeYamashita](https://github.com/KeisukeYamashita))
- Path escape import id of database user [\#373](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/373) ([KeisukeYamashita](https://github.com/KeisukeYamashita))
- Update cluster to match Atlas doc [\#372](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/372) ([themantissa](https://github.com/themantissa))
- INTMDB-147: Changed to required in schema of roles for database users [\#371](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/371) ([coderGo93](https://github.com/coderGo93))
- INTMDB-144: Updated for scopes database users [\#368](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/368) ([coderGo93](https://github.com/coderGo93))
- Fix database user resource broken indent [\#367](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/367) ([KeisukeYamashita](https://github.com/KeisukeYamashita))
- INTMDB-142: Fixes the bug for alertconfiguration using data dog [\#366](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/366) ([coderGo93](https://github.com/coderGo93))
- INTMDB-133: Updated Encryption At Rest to work with IAM Roles [\#365](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/365) ([coderGo93](https://github.com/coderGo93))
- INTMDB-141: Fixing 404 for existing database user [\#362](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/362) ([leofigy](https://github.com/leofigy))
- INTMDB-121: Prevents removing existing users [\#361](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/361) ([leofigy](https://github.com/leofigy))
- update wording [\#359](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/359) ([themantissa](https://github.com/themantissa))
- update the documentation and examples for adding the replication spec… [\#357](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/357) ([nikhil-mongo](https://github.com/nikhil-mongo))
- fix: code cleanup [\#355](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/355) ([AlphaWong](https://github.com/AlphaWong))
- Cloud Access Provider Datasources, Resources, and Documentation \(INTMDB 131\) [\#352](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/352) ([leofigy](https://github.com/leofigy))
- doc fix for db users data source [\#351](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/351) ([themantissa](https://github.com/themantissa))
- AWS and AZURE Private Endpoints \(INTMDB-123 & INTMDB-124\) [\#349](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/349) ([coderGo93](https://github.com/coderGo93))
- Basicexample for starting with Atlas [\#345](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/345) ([nikhil-mongo](https://github.com/nikhil-mongo))
- Fix update function for DB users [\#341](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/341) ([EricZaporzan](https://github.com/EricZaporzan))

## [v0.7.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.7.0) (2020-10-23)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.5...v0.7.0)

**Fixed bugs:**

- X509 user creation and update throws error  [\#312](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/312)
- inherited\_roles are not correctly removed from custom\_db\_roles [\#280](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/280)
- alertConfigs fix field update in terraform state [\#334](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/334) ([leofigy](https://github.com/leofigy))

**Closed issues:**

- Warning when installing the provider on Terraform 0.13 [\#342](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/342)
- mongodbatals\_network\_container [\#336](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/336)
- Typo in documentation  [\#335](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/335)
- Crash when refreshing TF State for a single cluster [\#330](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/330)
- 500 response on /groups/GROUP-ID/peers [\#320](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/320)
- An invalid enumeration value M5 was specified. [\#318](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/318)
- Container ID on the cluster data source is always empty [\#314](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/314)
- email\_enabled always reported as a change for mongodbatlas\_alert\_configuration [\#306](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/306)

**Merged pull requests:**

- Quick docs for 3rd party [\#344](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/344) ([themantissa](https://github.com/themantissa))
- chore: changelog v0.7.0 [\#343](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/343) ([leofigy](https://github.com/leofigy))
- documentation fix \#335 and examples added for the Azure VNET peering and AWS Private Link [\#340](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/340) ([nikhil-mongo](https://github.com/nikhil-mongo))
- mongodbatlas\_alert\_configuration - reset ID if was deleted and it's already in the plan [\#333](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/333) ([leofigy](https://github.com/leofigy))
- New resource and datasource for Project IP Access list [\#332](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/332) ([coderGo93](https://github.com/coderGo93))
- Client upgrade to fix metric threshold value set as 0 [\#331](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/331) ([leofigy](https://github.com/leofigy))
- docs: add mongo SLA link [\#328](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/328) ([gssbzn](https://github.com/gssbzn))
- Example added for database user scope [\#327](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/327) ([nikhil-mongo](https://github.com/nikhil-mongo))
- Add "Sensitive: true" for securing sensitive data in state [\#325](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/325) ([KeisukeYamashita](https://github.com/KeisukeYamashita))
- Create README and examples directory [\#324](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/324) ([themantissa](https://github.com/themantissa))
- fix: fixes a bug for issue 289 [\#323](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/323) ([coderGo93](https://github.com/coderGo93))
- Third party integrations  [\#321](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/321) ([leofigy](https://github.com/leofigy))
- changed from running on PR to manually trigger acceptance tests [\#319](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/319) ([coderGo93](https://github.com/coderGo93))
- Cluster docs [\#317](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/317) ([nikhil-mongo](https://github.com/nikhil-mongo))
- chore: changelog for v0.6.5 [\#316](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/316) ([leofigy](https://github.com/leofigy))
- Chore: Fix the ProviderVersion in the useragent string [\#309](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/309) ([MihaiBojin](https://github.com/MihaiBojin))

## [v0.6.5](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.5) (2020-09-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.6...v0.6.5)

## [v0.6.6](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.6) (2020-09-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.4...v0.6.6)

**Fixed bugs:**

- X509 is using the wrong authentication database when updating an existing user  [\#292](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/292)
- `mongodbatlas_cloud_provider_snapshot_backup_policy` `restore_window_days` \(optional value\) is being set even when omitted in resource config [\#290](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/290)

**Closed issues:**

- "mongodbatlas\_alert\_configuration" prints Slack API token in plain text [\#310](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/310)
- Can we create search index using terraform-provider-mongodbatlas? [\#308](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/308)
- Error: rpc error: code = Unavailable desc = transport is closing [\#302](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/302)
- Can't create alerts with using many event\_type  [\#232](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/232)

**Merged pull requests:**

- mongo atlas client update fix \#292 \#312 [\#315](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/315) ([leofigy](https://github.com/leofigy))
- DB user creation error because bad encoding in path [\#313](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/313) ([leofigy](https://github.com/leofigy))
- Database user scopes [\#307](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/307) ([coderGo93](https://github.com/coderGo93))
- Setting deterministic encoding id output, just sorting the keys [\#303](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/303) ([leofigy](https://github.com/leofigy))

## [v0.6.4](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.4) (2020-09-02)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.3...v0.6.4)

**Fixed bugs:**

- Unable to import $external auth users [\#285](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/285)
- custom\_db\_roles cannot be created with only inherited roles [\#279](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/279)
- mongodbatlas\_team data provider team\_id null after successful API query [\#277](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/277)

**Closed issues:**

- There is no parity between the Atlas API documentation and the provider doc in regards to alert event\_type values [\#295](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/295)
- Renaming a custom\_db\_role with attached users is not possible [\#284](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/284)
- changing cluster to \_NVME fails on commented-out IOPS [\#283](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/283)
- Error when assigning a custom db role to a database user.  [\#273](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/273)
- Error when creating `mongodbatlas_project_ip_whitelist` resource [\#266](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/266)
- Can't create a alert for Replication Oplog Window [\#227](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/227)

**Merged pull requests:**

- chore: add Changelog for 0.6.4 [\#301](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/301) ([marinsalinas](https://github.com/marinsalinas))
- fix: added a validation for restore\_window\_days [\#300](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/300) ([PacoDw](https://github.com/PacoDw))
- update auth\_database\_name [\#299](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/299) ([themantissa](https://github.com/themantissa))
- Fix \#227 \#232: Added a new Threshold attribute [\#298](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/298) ([PacoDw](https://github.com/PacoDw))
- Fix \#285: Unable to import $external auth users [\#297](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/297) ([PacoDw](https://github.com/PacoDw))
- Fixes many testacc [\#296](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/296) ([coderGo93](https://github.com/coderGo93))
- Fix \#279 \#280 [\#294](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/294) ([PacoDw](https://github.com/PacoDw))
- GitHub actions tests [\#293](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/293) ([coderGo93](https://github.com/coderGo93))
- Changed the harcoded links from hashicorp repo to mongodb repo [\#288](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/288) ([coderGo93](https://github.com/coderGo93))
- add note about container creation [\#287](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/287) ([themantissa](https://github.com/themantissa))
- Correct cluster labels documentation [\#286](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/286) ([themantissa](https://github.com/themantissa))
- Add templates to repo [\#282](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/282) ([themantissa](https://github.com/themantissa))
- Fix \#277: mongodbatlas\_team data provider team\_id null after successful API query [\#281](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/281) ([PacoDw](https://github.com/PacoDw))
- Create SECURITY.md [\#278](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/278) ([themantissa](https://github.com/themantissa))
- Update README.md [\#276](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/276) ([themantissa](https://github.com/themantissa))
- Release configuration [\#275](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/275) ([MihaiBojin](https://github.com/MihaiBojin))
- Improvement for 503 error response while creating a cluster [\#274](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/274) ([coderGo93](https://github.com/coderGo93))
- Cleaned vendored deps [\#272](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/272) ([MihaiBojin](https://github.com/MihaiBojin))
- Replaced the digest auth library with one that supports SHA-256 [\#271](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/271) ([MihaiBojin](https://github.com/MihaiBojin))
- Updated changelog v0.6.3 [\#270](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/270) ([PacoDw](https://github.com/PacoDw))
- fix: fix golangci lint and travis [\#269](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/269) ([gssbzn](https://github.com/gssbzn))
- feat: add a unique user agent [\#268](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/268) ([gssbzn](https://github.com/gssbzn))
- fix: added validation for autoscaling compute enabled and when true a… [\#267](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/267) ([coderGo93](https://github.com/coderGo93))
- Added a field AwsIAMType for database user [\#264](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/264) ([coderGo93](https://github.com/coderGo93))
- Updated Backup Policy documentation [\#259](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/259) ([PacoDw](https://github.com/PacoDw))

## [v0.6.3](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.3) (2020-07-27)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.2...v0.6.3)

**Fixed bugs:**

- Can't create a new cluster \(M2/M5\) after 0.6.2 version [\#265](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/265)

**Closed issues:**

- Creating cluster eventually returns 503 [\#256](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/256)

## [v0.6.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.2) (2020-07-16)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.1...v0.6.2)

**Fixed bugs:**

- Adding 16 whitelist entries at the same time causes an error [\#252](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/252)
- Error when create or import cluster - panic: runtime error: invalid memory address or nil pointer dereference [\#243](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/243)
- Cannot re-apply config when M2/M5 `disk_size_gb` is specified incorrectly [\#115](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/115)
- accepter\_region\_name not required for AWS on read/import/update [\#53](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/53)
- fix: resource/project\_ip\_whitelist - modify ip whitelist entry valida… [\#257](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/257) ([marinsalinas](https://github.com/marinsalinas))

**Closed issues:**

- In recommendations, prevent export of keys appearing in OS history [\#261](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/261)

**Merged pull requests:**

- Small change to recommendations [\#263](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/263) ([themantissa](https://github.com/themantissa))
- Updated changelog v0.6.2 [\#262](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/262) ([PacoDw](https://github.com/PacoDw))
- Updated go version to v1.14 [\#260](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/260) ([PacoDw](https://github.com/PacoDw))
- Fix auto scaling attributes [\#255](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/255) ([PacoDw](https://github.com/PacoDw))
- add: project\_ip\_whitelist datasource [\#254](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/254) ([gmlp](https://github.com/gmlp))
- imp: team datasource add team name option [\#253](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/253) ([gmlp](https://github.com/gmlp))
- fix: fixes \#115  issue with disk size for shared tiers [\#251](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/251) ([gmlp](https://github.com/gmlp))
- Added golangci configuration and travis fix [\#248](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/248) ([PacoDw](https://github.com/PacoDw))
- Updated the name of module client mongodb atlas [\#244](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/244) ([coderGo93](https://github.com/coderGo93))
- fix: fixes \#53 accepter\_region\_name not required for AWS on read/import/update [\#242](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/242) ([gmlp](https://github.com/gmlp))

## [v0.6.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.1) (2020-06-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.6.0...v0.6.1)

**Fixed bugs:**

- Error when use provider\_name = TENANT on 0.6.0 mongodbatlas provider version. [\#246](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/246)

**Closed issues:**

- Add MongoDB Collection Data Source [\#250](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/250)

**Merged pull requests:**

- Updated changelog v0.6.1 [\#249](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/249) ([PacoDw](https://github.com/PacoDw))
- Fix \#246: Error when use provider\_name = TENANT on 0.6.0 mongodbatlas provider version [\#247](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/247) ([PacoDw](https://github.com/PacoDw))
- Fix \#243: Error when create or import cluster - panic: runtime error: invalid memory address or nil pointer dereference [\#245](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/245) ([PacoDw](https://github.com/PacoDw))

## [v0.6.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.6.0) (2020-06-11)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.5.1...v0.6.0)

**Implemented enhancements:**

- mongodbatlas\_database\_user can not be imported when they contain dashes "-" in the name [\#179](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/179)

**Fixed bugs:**

- Changes to mongodbatlas\_database\_user.role.collection\_name are ignored [\#228](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/228)
- Hour and minute properties don't update when they are zero for mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy [\#211](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/211)
- Issues with advanced\_configuration section on mongodbatlas\_cluster [\#210](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/210)
- Changes are not detected when changing Team's role\_names attribute on mongodbatlas\_project [\#209](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/209)
- terraform plan and apply fails after upgrading this module to 0.5.0 [\#200](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/200)
- Issues upgrading cluster to an AWS NVME tier. [\#132](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/132)

**Closed issues:**

- Updating Snapshot Backup Policy: This resource requires access through a whitelist of ip ranges. [\#235](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/235)
- Cannot import mongodbatlas\_database\_user if username contains a hyphen [\#234](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/234)
- How to create a custom db role using built-in and connection action [\#226](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/226)
- connection\_strings returning empty private values [\#220](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/220)
- Documentation incorrect about accessing connection\_strings from clusters? [\#219](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/219)
- Incorrect description for atlas\_cidr\_block in mongodbatlas\_network\_peering documentation [\#215](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/215)
- RESOURSE or RESOURCE? Spelling change for readme.md [\#185](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/185)
- mongodbatlas\_encryption\_at\_rest key rotation impossible to perform with Azure KeyVault [\#80](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/80)

**Merged pull requests:**

- chore: updated changelog to v0.6.0 [\#241](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/241) ([PacoDw](https://github.com/PacoDw))
- Documentation Improvements and New Guide for 0.6.0 [\#240](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/240) ([themantissa](https://github.com/themantissa))
- fixes \#210: Issues with advanced\_configuration section on mongodbatlas\_cluster [\#238](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/238) ([gmlp](https://github.com/gmlp))
- New parameters about pagination for datasources [\#237](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/237) ([coderGo93](https://github.com/coderGo93))
- fix: fixes \#132 issues upgrading cluster to an AWS NVME tier [\#236](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/236) ([gmlp](https://github.com/gmlp))
- Cluster autoscaling [\#233](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/233) ([coderGo93](https://github.com/coderGo93))
- Fix \#228: Changes to mongodbatlas\_database\_user.role.collection\_name are ignored [\#231](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/231) ([PacoDw](https://github.com/PacoDw))
- fixes \#211: Hour and minute properties don't update when they are zero for mongodbatlas\_cloud\_provider\_snapshot\_backup\_policy [\#230](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/230) ([gmlp](https://github.com/gmlp))
- Terraform sdk [\#229](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/229) ([PacoDw](https://github.com/PacoDw))
- Fix \#209: Changes are not detected when changing Team's role\_names attribute on mongodbatlas\_project [\#225](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/225) ([PacoDw](https://github.com/PacoDw))
- New fields for snapshot restore jobs [\#224](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/224) ([coderGo93](https://github.com/coderGo93))
- Improve connection string doc [\#223](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/223) ([themantissa](https://github.com/themantissa))
- Update network\_peering.html.markdown [\#217](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/217) ([themantissa](https://github.com/themantissa))
- fix: fixed DatabaseUserID to allows names with multiple dashes [\#214](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/214) ([PacoDw](https://github.com/PacoDw))
- Fix \#80 - Update for GCP Encryption at rest [\#212](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/212) ([coderGo93](https://github.com/coderGo93))
- Added field container\_id in resource cluster [\#208](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/208) ([coderGo93](https://github.com/coderGo93))

## [v0.5.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.5.1) (2020-04-27)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.5.0...v0.5.1)

**Implemented enhancements:**

- Support new private and privateSrv connection strings [\#183](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/183)

**Closed issues:**

- Alert configuration roles array should not be required [\#201](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/201)
- Can't get PrivateLink-aware mongodb+srv address when using privatelink [\#147](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/147)

**Merged pull requests:**

- chore: updated changelog file for v0.5.1 [\#207](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/207) ([PacoDw](https://github.com/PacoDw))
- Fix travis, remove google cookie [\#204](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/204) ([marinsalinas](https://github.com/marinsalinas))
- Fix: improved validation to avoid error 404 [\#203](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/203) ([PacoDw](https://github.com/PacoDw))
- Changed roles to computed [\#202](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/202) ([PacoDw](https://github.com/PacoDw))
- Fixed the documetation menu [\#199](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/199) ([PacoDw](https://github.com/PacoDw))

## [v0.5.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.5.0) (2020-04-22)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.4.2...v0.5.0)

**Implemented enhancements:**

- mongodbatlas\_encryption\_at\_rest outputs IAM secrets to stdout  [\#93](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/93)
- Cloud Provider Snapshot Backup Policy [\#180](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/180) ([PacoDw](https://github.com/PacoDw))

**Fixed bugs:**

- TERRAFORM CRASH on importing mongodbatlas\_alert\_configuration [\#171](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/171)

**Closed issues:**

- Problem using Cross Region Replica Set in GCP [\#188](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/188)
- Delete this please. [\#187](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/187)
- Confusing output when modifying a cluster [\#186](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/186)
- Cluster auto-scaling [\#182](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/182)
- Docs with wrong resource type [\#175](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/175)
- On upgrade from 0.4.1 to 0.4.2 start getting errors [\#174](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/174)
- Can't create SYSTEM\_NORMALIZED\_CPU\_IOWAIT alert [\#172](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/172)
- mongodbatlas\_alert\_configuration - not able to specify ROLE for type\_name = "GROUP" [\#153](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/153)

**Merged pull requests:**

- chore: update Changelog file for v0.5.0 version [\#197](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/197) ([marinsalinas](https://github.com/marinsalinas))
- Add CONTRIBUTING file [\#196](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/196) ([themantissa](https://github.com/themantissa))
- Update MongoSDK to v0.2.0 [\#195](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/195) ([marinsalinas](https://github.com/marinsalinas))
- Doc update for private\_ip\_mode [\#194](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/194) ([themantissa](https://github.com/themantissa))
- Peering Container documentation fix [\#193](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/193) ([themantissa](https://github.com/themantissa))
- Update backup documenation [\#191](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/191) ([themantissa](https://github.com/themantissa))
- Fix documentation of roles block role\_name [\#184](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/184) ([fbreckle](https://github.com/fbreckle))
- Connection strings [\#181](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/181) ([coderGo93](https://github.com/coderGo93))
- Typo in `provider_disk_type_name` description [\#178](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/178) ([caitlinelfring](https://github.com/caitlinelfring))
- added roles in schema of notifications for alert configurations [\#177](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/177) ([coderGo93](https://github.com/coderGo93))
- fix-\#175 - missing word in resource name [\#176](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/176) ([themantissa](https://github.com/themantissa))
- Fix \#171: added validation to avoid nil type error [\#173](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/173) ([PacoDw](https://github.com/PacoDw))
- Fix Attributes Reference bullet points [\#168](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/168) ([brunopadz](https://github.com/brunopadz))

## [v0.4.2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.4.2) (2020-03-12)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.4.1...v0.4.2)

**Fixed bugs:**

- mongodbatlas\_cluster fails to redeploy manually deleted cluster [\#159](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/159)

**Closed issues:**

- mongodbatlas\_alert\_configuration - not able to generate any alerts with event\_type = "OUTSIDE\_METRIC\_THRESHOLD" and matcher.fieldName != "TYPE\_NAME" [\#164](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/164)
- Still cannot create cluster in region ME\_SOUTH\_1 on plugin version 0.4.1 [\#161](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/161)
- mongoatlas\_cluster fails to create  - invalid enumeration value M2 was specified  [\#160](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/160)
- Can't create cluster ME\_SOUTH\_1 region [\#157](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/157)

**Merged pull requests:**

- chore: fix linting issues [\#169](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/169) ([marinsalinas](https://github.com/marinsalinas))
- chore: add changelog file for 0.4.2 version [\#167](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/167) ([marinsalinas](https://github.com/marinsalinas))
- Doc: Fix import for mongodbatlas\_project\_ip\_whitelist [\#166](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/166) ([haidaraM](https://github.com/haidaraM))
- chore: removed wrong validation for matchers.value [\#165](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/165) ([PacoDw](https://github.com/PacoDw))
- feature: add default label to clusters [\#163](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/163) ([marinsalinas](https://github.com/marinsalinas))
- Cleaned Cluster state when it isn't found to allow create it again [\#162](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/162) ([PacoDw](https://github.com/PacoDw))
- cluster: removed array of regions due to they could be changed [\#158](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/158) ([PacoDw](https://github.com/PacoDw))

## [v0.4.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.4.1) (2020-02-26)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.4.0...v0.4.1)

**Fixed bugs:**

- Add name argument in mongodbatlas\_project datasource [\#140](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/140)

**Closed issues:**

- Delete timeout for mongodbatlas\_private\_endpoint resource too short [\#151](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/151)
- mongodbatlas\_project name [\#150](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/150)
- mongodbatlas\_custom\_db\_role not waiting for resource creation [\#148](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/148)
- Cannot use mongodbatlas\_maintenance\_window - Error provider does not support [\#145](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/145)
- Error creating users with mongodbatlas\_database\_user \(following documentation examples\) [\#144](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/144)
- Auto Scale Cluster Tier Missing [\#141](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/141)

**Merged pull requests:**

- chore: add changelog file for 0.4.1 version [\#156](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/156) ([marinsalinas](https://github.com/marinsalinas))
- Custom DB Roles: added refresh function to allow to create/remove multiple custom roles at the same time [\#155](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/155) ([PacoDw](https://github.com/PacoDw))
- chore: increase timeout when delete in private\_endpoint resource [\#154](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/154) ([marinsalinas](https://github.com/marinsalinas))
- add upgrade guide [\#149](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/149) ([themantissa](https://github.com/themantissa))
- Correct `mongodbatlas_teams` resource name in docs [\#143](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/143) ([mattt416](https://github.com/mattt416))
- Project data source [\#142](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/142) ([PacoDw](https://github.com/PacoDw))

## [v0.4.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.4.0) (2020-02-18)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.3.1...v0.4.0)

**Implemented enhancements:**

- expose 'paused' as an argument for mongodbatlas\_cluster [\#105](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/105)
- Add pitEnabled feature of mongodbatlas\_cluster resource [\#104](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/104)
- Support for AWS security groups in mongodbatlas\_project\_ip\_whitelist [\#67](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/67)

**Fixed bugs:**

- Cannot update GCP network peer [\#86](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/86)
- Cluster fails to build on 0.3.1 when mongo\_db\_major\_version is not specified [\#81](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/81)
- Crash \(panic, interface conversion error\) when creating mongodbatlas\_encryption\_at\_rest in Azure [\#74](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/74)
- Creating M2 cluster without specifying disk\_size\_gb results in 400 Bad Request [\#72](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/72)

**Closed issues:**

- add mongodbatlas\_project datasource [\#137](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/137)
- Inconsistent documentation between GitHub repo and Terraform site [\#136](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/136)
- Cloud provider snapshot management [\#124](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/124)
- Add support in cluster-tier autoscaling [\#123](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/123)
- Continuous Backup is not supported for \(new\) AWS clusters [\#121](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/121)
- If don't specify IOPS when creating M10 or M20 cluster a 0 value is passed in causing failure [\#120](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/120)
- When is mongodbatlas\_project\_ip\_whitelist security group feature going to be released? [\#114](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/114)
- Error creating MongoDB Cluster: unexpected EOF [\#110](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/110)
- Issue with import  mongodbatlas\_cloud\_provider\_snapshot\_restore\_job [\#109](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/109)
- mongodbatlas\_network\_container Already exists [\#88](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/88)
- mongodbatlas\_network\_container doesn't form a valid json request [\#83](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/83)
- mongodbatlas\_network\_containers datasource doesn't work with Azure [\#71](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/71)
- missing schema for provider "mongodbatlas" resource type mongodbatlas\_ip\_whitelist [\#70](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/70)
- Whitelisted Project IPs when manually deleted causes failure at next plan/apply [\#68](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/68)
- Modifying project ip whitelist destroy and re-create all resources [\#51](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/51)

**Merged pull requests:**

- Changelog for v0.4.0 [\#138](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/138) ([marinsalinas](https://github.com/marinsalinas))
- Readme: Updated env variables [\#135](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/135) ([PacoDw](https://github.com/PacoDw))
- Database Users: updated Read Function to avoid plugin error when it upgrades [\#133](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/133) ([PacoDw](https://github.com/PacoDw))
- Fix snapshot import with hyphened cluster\_name [\#131](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/131) ([marinsalinas](https://github.com/marinsalinas))
- Spelling and grammer [\#130](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/130) ([CMaylone](https://github.com/CMaylone))
- chore: added database\_name as deprecated attribute [\#129](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/129) ([PacoDw](https://github.com/PacoDw))
- Encryption At Rest: fixed issues and added an enhancement [\#128](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/128) ([PacoDw](https://github.com/PacoDw))
- Add PIT enabled argumento to Cluster Resource and Data Source [\#126](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/126) ([marinsalinas](https://github.com/marinsalinas))
- X509 Authentication Database User [\#125](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/125) ([PacoDw](https://github.com/PacoDw))
- Database users: added x509\_type attribute [\#122](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/122) ([PacoDw](https://github.com/PacoDw))
- Shared tier doc edits [\#119](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/119) ([themantissa](https://github.com/themantissa))
- Private endpoints [\#118](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/118) ([PacoDw](https://github.com/PacoDw))
- Update cluster doc [\#117](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/117) ([themantissa](https://github.com/themantissa))
- Update backup, add links [\#116](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/116) ([themantissa](https://github.com/themantissa))
- Projects: adding teams attribute [\#113](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/113) ([PacoDw](https://github.com/PacoDw))
- Update cluster.html.markdown [\#112](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/112) ([themantissa](https://github.com/themantissa))
- Fix DiskSizeGB missing [\#111](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/111) ([marinsalinas](https://github.com/marinsalinas))
- Terraform resource for MongoDB Custom Roles [\#108](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/108) ([PacoDw](https://github.com/PacoDw))
- Fix peering resource [\#107](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/107) ([PacoDw](https://github.com/PacoDw))
- Fix \#68: Added the ability to re-create the whitelist entry when it's remove manually [\#106](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/106) ([PacoDw](https://github.com/PacoDw))
- Updating `git clone` command to reference current repository [\#103](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/103) ([macintacos](https://github.com/macintacos))
- Cluster label and plugin attribute [\#102](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/102) ([PacoDw](https://github.com/PacoDw))
- Added functions to handle labels attribute in some resources [\#101](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/101) ([PacoDw](https://github.com/PacoDw))
- Added labels attr for Database User resource [\#100](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/100) ([PacoDw](https://github.com/PacoDw))
- Alert configuration resource and data source [\#99](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/99) ([PacoDw](https://github.com/PacoDw))
- Update database\_user.html.markdown [\#98](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/98) ([themantissa](https://github.com/themantissa))
- update containers and ip whitelist doc [\#96](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/96) ([themantissa](https://github.com/themantissa))
- Add provider\_name to containers data source [\#95](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/95) ([marinsalinas](https://github.com/marinsalinas))
- Whitelist [\#94](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/94) ([PacoDw](https://github.com/PacoDw))
- Network Peering RS: remove provider\_name=AWS as default, use Required=true instead i… [\#92](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/92) ([marinsalinas](https://github.com/marinsalinas))
- Update project.html.markdown [\#91](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/91) ([themantissa](https://github.com/themantissa))
- Feat: Global Cluster Configuration Resource and Data Source. [\#90](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/90) ([marinsalinas](https://github.com/marinsalinas))
- fix: validate if mongo\_db\_major\_version is set in cluster resource [\#85](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/85) ([marinsalinas](https://github.com/marinsalinas))
- Auditing Resource and Data Source [\#82](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/82) ([PacoDw](https://github.com/PacoDw))
- Feat: Team Resource and Data Source [\#79](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/79) ([marinsalinas](https://github.com/marinsalinas))
- Maintenance window ds [\#78](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/78) ([PacoDw](https://github.com/PacoDw))
- Added default Disk Size when it doesn't set up on cluster resource [\#77](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/77) ([PacoDw](https://github.com/PacoDw))
- Maintenance window rs [\#76](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/76) ([PacoDw](https://github.com/PacoDw))
- website: collapse data sources sidebar by default [\#75](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/75) ([marinsalinas](https://github.com/marinsalinas))
- Improvements to Peering Resources [\#73](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/73) ([themantissa](https://github.com/themantissa))
- Remove dupe argument in docs [\#69](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/69) ([ktmorgan](https://github.com/ktmorgan))
- Clarify Azure Option in Doc [\#66](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/66) ([themantissa](https://github.com/themantissa))

## [v0.3.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.3.1) (2019-11-11)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.3.0...v0.3.1)

**Fixed bugs:**

- Confirmation on timelimit for a terraform apply [\#57](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/57)

**Closed issues:**

- Not able to create M0 clusters [\#64](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/64)
- No way to modify advanced configuration options for a cluster [\#61](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/61)
- mongodbatlas\_network\_peering outputting invalid json [\#59](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/59)
- Syntax are not mandatory and creates confusion [\#58](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/58)
- data source mongodbatlas\_network\_peering retrieves the same for id and connection\_id [\#56](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/56)
- Add resource for maintenance window [\#55](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/55)
- Error encryption\_at\_rest  rpc unavailable desc [\#54](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/54)
- specify oplog size? [\#52](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/52)
- Add resource for custom database roles [\#50](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/50)
- An invalid enumeration value US\_EAST\_1 was specified. [\#49](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/49)
- Version 0.3.0 [\#47](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/47)
- plugin.terraform-provider-mongodbatlas\_v0.2.0\_x4: panic: runtime error: index out of range [\#36](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/36)

**Merged pull requests:**

- chore: add changelog file for 0.3.1 version [\#65](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/65) ([marinsalinas](https://github.com/marinsalinas))
- Added format function to handle the mongo\_db\_major\_version attr [\#63](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/63) ([PacoDw](https://github.com/PacoDw))
- Added cast func to avoid panic by nil value [\#62](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/62) ([PacoDw](https://github.com/PacoDw))
- Cluster advanced configuration Options [\#60](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/60) ([PacoDw](https://github.com/PacoDw))

## [v0.3.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.3.0) (2019-10-08)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.2.0...v0.3.0)

**Closed issues:**

- Upgrade from M2 to M10 fails [\#42](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/42)
- GCP Peering endless terraform apply [\#41](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/41)
- AWS clusters default provider\_encrypt\_ebs\_volume to false [\#40](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/40)
- mongodbatlas\_network\_peering Internal Servier Error [\#35](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/35)
- Problem encryption\_at\_rest [\#33](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/33)
- Problem destroying network peering container [\#30](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/30)
- Bug VPC Peering between GCP and Atlas [\#29](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/29)

**Merged pull requests:**

- chore: add changelog file for 0.3.0 version [\#48](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/48) ([marinsalinas](https://github.com/marinsalinas))
- Clarify Doc Examples and Text [\#46](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/46) ([themantissa](https://github.com/themantissa))
- fix-\#40: added true value by defualt on provider\_encrypt\_ebs\_volume attr [\#45](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/45) ([PacoDw](https://github.com/PacoDw))
- make provider\_name forced new to avoid patch problems [\#44](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/44) ([marinsalinas](https://github.com/marinsalinas))
- Network peering [\#43](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/43) ([PacoDw](https://github.com/PacoDw))
- Update readme with more info [\#39](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/39) ([themantissa](https://github.com/themantissa))
- Fix: Network container [\#38](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/38) ([PacoDw](https://github.com/PacoDw))
- Doc updates [\#37](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/37) ([themantissa](https://github.com/themantissa))

## [v0.2.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.2.0) (2019-09-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.1.1...v0.2.0)

**Closed issues:**

- Unable to create project with peering only connections [\#24](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/24)
- importing a mongodbatlas\_project\_ip\_whitelist resource does not save project\_id to state [\#21](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/21)
- Support the vscode terraform extension [\#19](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/19)
- Bug: VPC Peering Atlas-GCP [\#17](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/17)
- PATCH network peering failed with no peer found [\#14](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/14)

**Merged pull requests:**

- chore: add changelog for new release [\#34](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/34) ([marinsalinas](https://github.com/marinsalinas))
- Add Private IP Mode Resource. [\#32](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/32) ([marinsalinas](https://github.com/marinsalinas))
- Moved provider\_name values to the correct section [\#31](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/31) ([kgriffiths](https://github.com/kgriffiths))
- website: add links to Atlas Region name reference. [\#28](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/28) ([themantissa](https://github.com/themantissa))
- Encryption at rest fix [\#27](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/27) ([marinsalinas](https://github.com/marinsalinas))
- website: make resources side nav expanded as default [\#25](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/25) ([marinsalinas](https://github.com/marinsalinas))
- fix: importing a mongodbatlas\_project\_ip\_whitelist resource does not save project\_id to state [\#23](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/23) ([PacoDw](https://github.com/PacoDw))
- Fix \#14: PATCH network peering failed with no peer found [\#22](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/22) ([PacoDw](https://github.com/PacoDw))
- fix: change the test configuration for AWS and GCP [\#20](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/20) ([PacoDw](https://github.com/PacoDw))

## [v0.1.1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.1.1) (2019-09-05)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/v0.1.0...v0.1.1)

**Fixed bugs:**

- panic: runtime error: index out of range [\#1](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/1)

**Closed issues:**

- GCP peering problem [\#16](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/16)
- Cluster creation with Azure provider failed [\#15](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/15)
- Error creating MongoDB Cluster [\#9](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/9)
- Failed to create Atlas network peering container [\#7](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/7)
- Bug: Invalid attribute diskIOPS specified. [\#2](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/2)

**Merged pull requests:**

- chore: update changelog [\#18](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/18) ([marinsalinas](https://github.com/marinsalinas))
- website: fix typo [\#13](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/13) ([heimweh](https://github.com/heimweh))
- fix: add the correct func to check the env variables on peering datasources [\#12](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/12) ([PacoDw](https://github.com/PacoDw))
- Fix diskIOPS attribute for GCP and Azure [\#11](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/11) ([PacoDw](https://github.com/PacoDw))
- website: fix data sources sidebar always collapsed [\#10](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/10) ([marinsalinas](https://github.com/marinsalinas))
- mongodbatlas\_network\_\(peering and container\): add more testing case [\#8](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/8) ([PacoDw](https://github.com/PacoDw))
- website: fix typo in MongoDB Atlas Services [\#5](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/5) ([marinsalinas](https://github.com/marinsalinas))
- Ip whitelist entries: removing all entries whitelist by terraform user [\#4](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/4) ([PacoDw](https://github.com/PacoDw))
- Refactored import function to get all ip\_addresses and cird\_blocks entries [\#3](https://github.com/mongodb/terraform-provider-mongodbatlas/pull/3) ([PacoDw](https://github.com/PacoDw))

## [v0.1.0](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v0.1.0) (2019-08-19)

[Full Changelog](https://github.com/mongodb/terraform-provider-mongodbatlas/compare/3e1c5c44b56aee2f153ec618c804dd170bbefbd4...v0.1.0)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
