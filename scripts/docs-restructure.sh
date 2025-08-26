#!/usr/bin/env bash
set -euo pipefail

base="docs"
mkdir -p \
  "$base/Clusters/Resources" "$base/Clusters/Data Sources" \
  "$base/Cluster Outage Simulation/Resources" "$base/Cluster Outage Simulation/Data Sources" \
  "$base/Cloud Backups/Resources" "$base/Cloud Backups/Data Sources" \
  "$base/Cloud Provider Snapshots/Resources" "$base/Cloud Provider Snapshots/Data Sources" \
  "$base/Cloud Provider Access/Resources" "$base/Cloud Provider Access/Data Sources" \
  "$base/Projects/Resources" "$base/Projects/Data Sources" \
  "$base/Organizations/Resources" "$base/Organizations/Data Sources" \
  "$base/Network Peering/Resources" "$base/Network Peering/Data Sources" \
  "$base/Private Endpoint Services/Resources" "$base/Private Endpoint Services/Data Sources" \
  "$base/Data Federation/Resources" "$base/Data Federation/Data Sources" \
  "$base/Data Lake Pipelines/Resources" "$base/Data Lake Pipelines/Data Sources" \
  "$base/Search/Resources" "$base/Search/Data Sources" \
  "$base/Streams/Resources" "$base/Streams/Data Sources" \
  "$base/Serverless Instances/Resources" "$base/Serverless Instances/Data Sources" \
  "$base/Serverless Private Endpoints/Resources" "$base/Serverless Private Endpoints/Data Sources" \
  "$base/Database Users/Resources" "$base/Database Users/Data Sources" \
  "$base/LDAP Configuration/Resources" "$base/LDAP Configuration/Data Sources" \
  "$base/Auditing/Resources" "$base/Auditing/Data Sources" \
  "$base/X.509 Authentication/Resources" "$base/X.509 Authentication/Data Sources" \
  "$base/Alert Configurations/Resources" "$base/Alert Configurations/Data Sources" \
  "$base/Third-Party Integrations/Resources" "$base/Third-Party Integrations/Data Sources" \
  "$base/Online Archive/Resources" "$base/Online Archive/Data Sources" \
  "$base/Event Trigger/Resources" "$base/Event Trigger/Data Sources" \
  "$base/AWS Clusters DNS/Resources" "$base/AWS Clusters DNS/Data Sources" \
  "$base/Root/Resources" "$base/Root/Data Sources" \
  "$base/Programmatic API Keys/Resources" "$base/Programmatic API Keys/Data Sources" \
  "$base/MongoDB Employee Access/Resources" "$base/MongoDB Employee Access/Data Sources" \
  "$base/Push-Based Log Export/Resources" "$base/Push-Based Log Export/Data Sources" \
  "$base/Resource Policies/Resources" "$base/Resource Policies/Data Sources" \
  "$base/Teams/Resources" "$base/Teams/Data Sources" \
  "$base/Atlas Users (Deprecated)/Resources" "$base/Atlas Users (Deprecated)/Data Sources" \
  "$base/Flex Clusters/Resources" "$base/Flex Clusters/Data Sources" \
  "$base/Flex Restore Jobs/Resources" "$base/Flex Restore Jobs/Data Sources" \
  "$base/Flex Snapshots/Resources" "$base/Flex Snapshots/Data Sources" \
  "$base/Auditing/Resources" "$base/Auditing/Data Sources" \
  "$base/Shared-Tier Restore Jobs/Resources" "$base/Shared-Tier Restore Jobs/Data Sources" \
  "$base/Shared-Tier Snapshots/Resources" "$base/Shared-Tier Snapshots/Data Sources" \
  "$base/Global Clusters/Resources" "$base/Global Clusters/Data Sources" \
  "$base/Custom Database Roles/Resources" "$base/Custom Database Roles/Data Sources" \
  "$base/Federated Authentication/Resources" "$base/Federated Authentication/Data Sources" \
  "$base/X.509 Authentication/Resources" "$base/X.509 Authentication/Data Sources" \
  "$base/Maintenance Windows/Resources" "$base/Maintenance Windows/Data Sources" \
  "$base/Teams/Resources" "$base/Teams/Data Sources" \
  "$base/Project IP Access List/Resources" "$base/Project IP Access List/Data Sources" \
  "$base/Encryption at Rest using Customer Key Management/Resources" "$base/Encryption at Rest using Customer Key Management/Data Sources" \

# Clusters
git mv docs/resources/advanced_cluster.md "$base/Clusters/Resources/" || true
git mv docs/data-sources/advanced_cluster.md "$base/Clusters/Data Sources/" || true
git mv docs/data-sources/advanced_clusters.md "$base/Clusters/Data Sources/" || true
git mv docs/resources/mongodb_employee_access_grant.md "$base/Clusters/Resources/" || true
git mv docs/data-sources/mongodb_employee_access_grant.md "$base/Clusters/Data Sources/" || true
# Legacy Clusters (deprecated) moved under Clusters
git mv docs/resources/cluster.md "$base/Clusters/Resources/" || true
git mv docs/data-sources/cluster.md "$base/Clusters/Data Sources/" || true
git mv docs/data-sources/clusters.md "$base/Clusters/Data Sources/" || true





git mv docs/data-sources/flex_cluster.md "$base/Flex Clusters/Data Sources/" || true
git mv docs/data-sources/flex_clusters.md "$base/Flex Clusters/Data Sources/" || true
git mv docs/resources/flex_cluster.md "$base/Flex Clusters/Resources/" || true

git mv docs/data-sources/global_cluster_config.md "$base/Global Clusters/Data Sources/" || true
git mv docs/resources/global_cluster_config.md "$base/Global Clusters/Resources/" || true



# Cloud Backups
git mv docs/resources/backup_compliance_policy.md "$base/Cloud Backups/Resources/" || true
git mv docs/resources/cloud_backup_schedule.md "$base/Cloud Backups/Resources/" || true
git mv docs/resources/cloud_backup_snapshot.md "$base/Cloud Backups/Resources/" || true
git mv docs/resources/cloud_backup_snapshot_restore_job.md "$base/Cloud Backups/Resources/" || true
git mv docs/resources/cloud_backup_snapshot_export_bucket.md "$base/Cloud Backups/Resources/" || true
git mv docs/resources/cloud_backup_snapshot_export_job.md "$base/Cloud Backups/Resources/" || true
git mv docs/data-sources/backup_compliance_policy.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_schedule.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_snapshot.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_snapshots.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_snapshot_restore_job.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_snapshot_restore_jobs.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_snapshot_export_bucket.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_snapshot_export_buckets.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_snapshot_export_job.md "$base/Cloud Backups/Data Sources/" || true
git mv docs/data-sources/cloud_backup_snapshot_export_jobs.md "$base/Cloud Backups/Data Sources/" || true

# Cloud Provider Snapshots
git mv docs/resources/cloud_provider_snapshot.md "$base/Cloud Provider Snapshots/Resources/" || true
git mv docs/resources/cloud_provider_snapshot_backup_policy.md "$base/Cloud Provider Snapshots/Resources/" || true
git mv docs/resources/cloud_provider_snapshot_restore_job.md "$base/Cloud Provider Snapshots/Resources/" || true
git mv docs/data-sources/cloud_provider_snapshot.md "$base/Cloud Provider Snapshots/Data Sources/" || true
git mv docs/data-sources/cloud_provider_snapshots.md "$base/Cloud Provider Snapshots/Data Sources/" || true
git mv docs/data-sources/cloud_provider_snapshot_backup_policy.md "$base/Cloud Provider Snapshots/Data Sources/" || true
git mv docs/data-sources/cloud_provider_snapshot_restore_job.md "$base/Cloud Provider Snapshots/Data Sources/" || true
git mv docs/data-sources/cloud_provider_snapshot_restore_jobs.md "$base/Cloud Provider Snapshots/Data Sources/" || true
git mv docs/data-sources/cloud_provider_shared_tier_snapshot.md "$base/Shared-Tier Snapshots/Data Sources/" || true
git mv docs/data-sources/cloud_provider_shared_tier_snapshots.md "$base/Shared-Tier Snapshots/Data Sources/" || true
git mv docs/data-sources/cloud_provider_shared_tier_restore_job.md "$base/Shared-Tier Restore Jobs/Data Sources/" || true
git mv docs/data-sources/cloud_provider_shared_tier_restore_jobs.md "$base/Shared-Tier Restore Jobs/Data Sources/" || true
git mv docs/data-sources/flex_snapshot.md "$base/Flex Snapshots/Data Sources/" || true
git mv docs/data-sources/flex_snapshots.md "$base/Flex Snapshots/Data Sources/" || true
git mv docs/data-sources/flex_restore_job.md "$base/Flex Restore Jobs/Data Sources/" || true
git mv docs/data-sources/flex_restore_jobs.md "$base/Flex Restore Jobs/Data Sources/" || true

# Cloud Provider Access
git mv docs/resources/cloud_provider_access.md "$base/Cloud Provider Access/Resources/" || true
git mv docs/data-sources/cloud_provider_access_setup.md "$base/Cloud Provider Access/Data Sources/" || true

# Projects
git mv docs/resources/project.md "$base/Projects/Resources/" || true
git mv docs/resources/project_api_key.md "$base/Projects/Resources/" || true
git mv docs/resources/project_invitation.md "$base/Projects/Resources/" || true
git mv docs/resources/cloud_user_project_assignment.md "$base/Projects/Resources/" || true
git mv docs/data-sources/project.md "$base/Projects/Data Sources/" || true
git mv docs/data-sources/projects.md "$base/Projects/Data Sources/" || true
git mv docs/data-sources/project_api_key.md "$base/Projects/Data Sources/" || true
git mv docs/data-sources/project_api_keys.md "$base/Projects/Data Sources/" || true
git mv docs/data-sources/project_ip_addresses.md "$base/Projects/Data Sources/" || true
git mv docs/data-sources/project_invitation.md "$base/Projects/Data Sources/" || true
git mv docs/data-sources/cloud_user_project_assignment.md "$base/Projects/Data Sources/" || true


git mv docs/data-sources/project_ip_access_list.md "$base/Project IP Access List/Data Sources/" || true
git mv docs/resources/project_ip_access_list.md "$base/Project IP Access List/Resources/" || true


git mv docs/data-sources/resource_policy.md "$base/Resource Policies/Data Sources/" || true
git mv docs/data-sources/resource_policies.md "$base/Resource Policies/Data Sources/" || true
git mv docs/resources/resource_policy.md "$base/Resource Policies/Resources/" || true


git mv docs/data-sources/push_based_log_export.md "$base/Push-Based Log Export/Data Sources/" || true
git mv docs/resources/push_based_log_export.md "$base/Push-Based Log Export/Resources/" || true

git mv docs/data-sources/maintenance_window.md "$base/Maintenance Windows/Data Sources/" || true
git mv docs/resources/maintenance_window.md "$base/Maintenance Windows/Resources/" || true


# Push-Based Log Export
git mv docs/resources/push_based_log_export.md "$base/Push-Based Log Export/Resources/" || true
git mv docs/data-sources/push_based_log_export.md "$base/Push-Based Log Export/Data Sources/" || true

# Resource Policies
git mv docs/resources/resource_policy.md "$base/Resource Policies/Resources/" || true
git mv docs/data-sources/resource_policy.md "$base/Resource Policies/Data Sources/" || true
git mv docs/data-sources/resource_policies.md "$base/Resource Policies/Data Sources/" || true

# Organizations (incl. exception roles_org_id)
git mv docs/resources/organization.md "$base/Organizations/Resources/" || true
git mv docs/resources/org_invitation.md "$base/Organizations/Resources/" || true
git mv docs/resources/cloud_user_org_assignment.md "$base/Organizations/Resources/" || true
git mv docs/data-sources/organization.md "$base/Organizations/Data Sources/" || true
git mv docs/data-sources/organizations.md "$base/Organizations/Data Sources/" || true
git mv docs/data-sources/team.md "$base/Teams/Data Sources/" || true
git mv docs/data-sources/teams.md "$base/Teams/Data Sources/" || true
git mv docs/data-sources/roles_org_id.md "$base/Organizations/Data Sources/" || true
git mv docs/data-sources/org_invitation.md "$base/Organizations/Data Sources/" || true
git mv docs/data-sources/cloud_user_org_assignment.md "$base/Organizations/Data Sources/" || true


git mv docs/data-sources/cloud_user_team_assignment.md "$base/Teams/Data Sources/" || true
git mv docs/resources/cloud_user_team_assignment.md "$base/Teams/Resources/" || true
git mv docs/resources/team.md "$base/Teams/Resources/" || true
git mv docs/resources/teams.md "$base/Teams/Resources/" || true
git mv docs/data-sources/team_project_assignment.md "$base/Teams/Data Sources/" || true
git mv docs/resources/team_project_assignment.md "$base/Teams/Resources/" || true

# Network Peering
git mv docs/resources/network_peering.md "$base/Network Peering/Resources/" || true
git mv docs/resources/network_container.md "$base/Network Peering/Resources/" || true
git mv docs/data-sources/network_peering.md "$base/Network Peering/Data Sources/" || true
git mv docs/data-sources/network_peerings.md "$base/Network Peering/Data Sources/" || true
git mv docs/data-sources/network_container.md "$base/Network Peering/Data Sources/" || true
git mv docs/data-sources/network_containers.md "$base/Network Peering/Data Sources/" || true

# Private Endpoint Services
git mv docs/resources/privatelink_endpoint.md "$base/Private Endpoint Services/Resources/" || true
git mv docs/resources/privatelink_endpoint_service.md "$base/Private Endpoint Services/Resources/" || true
git mv docs/resources/private_endpoint_regional_mode.md "$base/Private Endpoint Services/Resources/" || true
git mv docs/data-sources/privatelink_endpoint.md "$base/Private Endpoint Services/Data Sources/" || true
git mv docs/data-sources/privatelink_endpoint_service.md "$base/Private Endpoint Services/Data Sources/" || true
git mv docs/data-sources/privatelink_endpoints_service_adl.md "$base/Private Endpoint Services/Data Sources/" || true
git mv docs/data-sources/private_endpoint_regional_mode.md "$base/PrivateLPrivate Endpoint Servicesink/Data Sources/" || true

# Data Federation
git mv docs/resources/federated_database_instance.md "$base/Data Federation/Resources/" || true
git mv docs/resources/privatelink_endpoint_service_data_federation_online_archive.md "$base/Data Federation/Resources/" || true
git mv docs/resources/federated_query_limit.md "$base/Data Federation/Resources/" || true
git mv docs/data-sources/federated_database_instance.md "$base/Data Federation/Data Sources/" || true
git mv docs/data-sources/federated_database_instances.md "$base/Data Federation/Data Sources/" || true
git mv docs/data-sources/privatelink_endpoint_service_data_federation_online_archive.md "$base/Data Federation/Data Sources/" || true
git mv docs/data-sources/privatelink_endpoint_service_data_federation_online_archives.md "$base/Data Federation/Data Sources/" || true
git mv docs/data-sources/federated_query_limit.md "$base/Data Federation/Data Sources/" || true
git mv docs/data-sources/federated_query_limits.md "$base/Data Federation/Data Sources/" || true


git mv docs/resources/federated_settings_identity_provider.md "$base/Federated Authentication/Resources/" || true
git mv docs/resources/federated_settings_org_config.md "$base/Federated Authentication/Resources/" || true
git mv docs/resources/federated_settings_org_role_mapping.md "$base/Federated Authentication/Resources/" || true
git mv docs/data-sources/federated_settings_identity_providers.md "$base/Data Federation/Data Sources/" || true
git mv docs/data-sources/federated_settings_org_role_mapping.md "$base/Federated Authentication/Data Sources/" || true
git mv docs/data-sources/federated_settings_org_role_mappings.md "$base/Federated Authentication/Data Sources/" || true
git mv docs/data-sources/federated_settings_identity_provider.md "$base/Federated Authentication/Data Sources/" || true
git mv docs/data-sources/federated_settings_org_config.md "$base/Federated Authentication/Data Sources/" || true
git mv docs/data-sources/federated_settings_org_configs.md "$base/Federated Authentication/Data Sources/" || true
git mv docs/data-sources/federated_settings.md "$base/Federated Authentication/Data Sources/" || true



# Data Lake Pipelines
git mv docs/resources/data_lake_pipeline.md "$base/Data Lake Pipelines/Resources/" || true
git mv docs/data-sources/data_lake_pipeline.md "$base/Data Lake Pipelines/Data Sources/" || true
git mv docs/data-sources/data_lake_pipelines.md "$base/Data Lake Pipelines/Data Sources/" || true
git mv docs/data-sources/data_lake_pipeline_run.md "$base/Data Lake Pipelines/Data Sources/" || true
git mv docs/data-sources/data_lake_pipeline_runs.md "$base/Data Lake Pipelines/Data Sources/" || true

# Search
git mv docs/resources/search_deployment.md "$base/Search/Resources/" || true
git mv docs/resources/search_index.md "$base/Search/Resources/" || true
git mv docs/data-sources/search_deployment.md "$base/Search/Data Sources/" || true
git mv docs/data-sources/search_index.md "$base/Search/Data Sources/" || true
git mv docs/data-sources/search_indexes.md "$base/Search/Data Sources/" || true

# Cluster Outage Simulation
git mv docs/resources/cluster_outage_simulation.md "$base/Cluster Outage Simulation/Resources/" || true
git mv docs/data-sources/cluster_outage_simulation.md "$base/Cluster Outage Simulation/Data Sources/" || true

# Streams
git mv docs/resources/stream_instance.md "$base/Streams/Resources/" || true
git mv docs/resources/stream_processor.md "$base/Streams/Resources/" || true
git mv docs/resources/stream_connection.md "$base/Streams/Resources/" || true
git mv docs/resources/stream_privatelink_endpoint.md "$base/Streams/Resources/" || true
git mv docs/data-sources/stream_instance.md "$base/Streams/Data Sources/" || true
git mv docs/data-sources/stream_instances.md "$base/Streams/Data Sources/" || true
git mv docs/data-sources/stream_processor.md "$base/Streams/Data Sources/" || true
git mv docs/data-sources/stream_processors.md "$base/Streams/Data Sources/" || true
git mv docs/data-sources/stream_connection.md "$base/Streams/Data Sources/" || true
git mv docs/data-sources/stream_connections.md "$base/Streams/Data Sources/" || true
git mv docs/data-sources/stream_privatelink_endpoint.md "$base/Streams/Data Sources/" || true
git mv docs/data-sources/stream_privatelink_endpoints.md "$base/Streams/Data Sources/" || true
git mv docs/data-sources/stream_account_details.md "$base/Streams/Data Sources/" || true

# Serverless Instances
git mv docs/resources/serverless_instance.md "$base/Serverless Instances/Resources/" || true
git mv docs/data-sources/serverless_instance.md "$base/Serverless Instances/Data Sources/" || true
git mv docs/data-sources/serverless_instances.md "$base/Serverless Instances/Data Sources/" || true


git mv docs/data-sources/privatelink_endpoints_service_serverless.md "$base/Serverless Private Endpoints/Data Sources/" || true
git mv docs/data-sources/privatelink_endpoint_service_serverless.md "$base/Serverless Private Endpoints/Data Sources/" || true
git mv docs/resources/privatelink_endpoint_serverless.md "$base/Serverless Private Endpoints/Resources/" || true
git mv docs/resources/privatelink_endpoint_service_serverless.md "$base/Serverless Private Endpoints/Resources/" || true

# Database Users
git mv docs/resources/database_user.md "$base/Database Users/Resources/" || true
git mv docs/data-sources/database_user.md "$base/Database Users/Data Sources/" || true
git mv docs/data-sources/database_users.md "$base/Database Users/Data Sources/" || true

git mv docs/data-sources/x509_authentication_database_user.md "$base/X.509 Authentication/Data Sources/" || true
git mv docs/resources/x509_authentication_database_user.md "$base/X.509 Authentication/Resources/" || true


# LDAP Configuration
git mv docs/resources/ldap_configuration.md "$base/LDAP Configuration/Resources/" || true
git mv docs/resources/ldap_verify.md "$base/LDAP Configuration/Resources/" || true
git mv docs/data-sources/ldap_verify.md "$base/LDAP Configuration/Data Sources/" || true
git mv docs/data-sources/ldap_configuration.md "$base/LDAP Configuration/Data Sources/" || true

git mv docs/resources/auditing.md "$base/Auditing/Resources/" || true

git mv docs/data-sources/auditing.md "$base/Auditing/Data Sources/" || true
git mv docs/resources/custom_db_role.md "$base/Custom Database Roles/Resources/" || true
git mv docs/data-sources/custom_db_role.md "$base/Custom Database Roles/Data Sources/" || true
git mv docs/data-sources/custom_db_roles.md "$base/Custom Database Roles/Data Sources/" || true



git mv docs/data-sources/encryption_at_rest_private_endpoints.md "$base/Encryption at Rest using Customer Key Management/Data Sources/" || true
git mv docs/data-sources/encryption_at_rest.md "$base/Encryption at Rest using Customer Key Management/Data Sources/" || true
git mv docs/data-sources/encryption_at_rest_private_endpoint.md "$base/Encryption at Rest using Customer Key Management/Data Sources/" || true
git mv docs/resources/encryption_at_rest.md "$base/Encryption at Rest using Customer Key Management/Resources/" || true
git mv docs/resources/encryption_at_rest_private_endpoint.md "$base/Encryption at Rest using Customer Key Management/Resources/" || true

# Alert Configurations
git mv docs/resources/alert_configuration.md "$base/Alert Configurations/Resources/" || true
git mv docs/data-sources/alert_configuration.md "$base/Alert Configurations/Data Sources/" || true
git mv docs/data-sources/alert_configurations.md "$base/Alert Configurations/Data Sources/" || true

# Third-Party Integrations
git mv docs/resources/third_party_integration.md "$base/Third-Party Integrations/Resources/" || true
git mv docs/data-sources/third_party_integration.md "$base/Third-Party Integrations/Data Sources/" || true
git mv docs/data-sources/third_party_integrations.md "$base/Third-Party Integrations/Data Sources/" || true

# Online Archive
git mv docs/resources/online_archive.md "$base/Online Archive/Resources/" || true
git mv docs/data-sources/online_archive.md "$base/Online Archive/Data Sources/" || true
git mv docs/data-sources/online_archives.md "$base/Online Archive/Data Sources/" || true

# Event Trigger (exception)
git mv docs/resources/event_trigger.md "$base/Event Trigger/Resources/" || true
git mv docs/data-sources/event_trigger.md "$base/Event Trigger/Data Sources/" || true
git mv docs/data-sources/event_triggers.md "$base/Event Trigger/Data Sources/" || true

# AWS Clusters DNS
git mv docs/resources/custom_dns_configuration_cluster_aws.md "$base/AWS Clusters DNS/Resources/" || true
git mv docs/data-sources/custom_dns_configuration_cluster_aws.md "$base/AWS Clusters DNS/Data Sources/" || true

# Root
git mv docs/data-sources/control_plane_ip_addresses.md "$base/Root/Data Sources/" || true

# Programmatic API Keys
git mv docs/resources/api_key.md "$base/Programmatic API Keys/Resources/" || true
git mv docs/resources/access_list_api_key.md "$base/Programmatic API Keys/Resources/" || true
git mv docs/resources/api_key_project_assignment.md "$base/Programmatic API Keys/Resources/" || true
git mv docs/data-sources/api_keys.md "$base/Programmatic API Keys/Data Sources/" || true
git mv docs/data-sources/api_key.md "$base/Programmatic API Keys/Data Sources/" || true
git mv docs/data-sources/access_list_api_key.md "$base/Programmatic API Keys/Data Sources/" || true
git mv docs/data-sources/access_list_api_keys.md "$base/Programmatic API Keys/Data Sources/" || true
git mv docs/data-sources/api_key_project_assignment.md "$base/Programmatic API Keys/Data Sources/" || true
git mv docs/data-sources/api_key_project_assignments.md "$base/Programmatic API Keys/Data Sources/" || true



# Atlas Users (Deprecated)
git mv docs/data-sources/atlas_user.md "$base/Atlas Users (Deprecated)/Data Sources/" || true
git mv docs/data-sources/atlas_users.md "$base/Atlas Users (Deprecated)/Data Sources/" || true



echo "Reorg complete. Please update any internal links and navigation as needed."
