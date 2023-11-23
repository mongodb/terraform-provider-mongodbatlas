package mongodbatlas

import (
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	DeprecationMessageParameterToResource = config.DeprecationMessageParameterToResource
	DeprecationByDateMessageParameter     = config.DeprecationByDateMessageParameter
	DeprecationByDateWithReplacement      = config.DeprecationByDateWithReplacement
	DeprecationByVersionMessageParameter  = config.DeprecationByVersionMessageParameter
	DeprecationMessage                    = config.DeprecationMessage
	AWS                                   = config.AWS
	AZURE                                 = config.AZURE
	errorProjectSetting                   = config.ErrorProjectSetting
	errorGetRead                          = "error reading cloud provider access %s"
)

type MongoDBClient = config.MongoDBClient

func encodeStateID(values map[string]string) string {
	return config.EncodeStateID(values)
}

func getEncodedID(stateID, keyPosition string) string {
	return config.GetEncodedID(stateID, keyPosition)
}

func decodeStateID(stateID string) map[string]string {
	return config.DecodeStateID(stateID)
}

func valRegion(reg any, opt ...string) (string, error) {
	return config.ValRegion(reg, opt...)
}

func removeLabel(list []matlas.Label, item matlas.Label) []matlas.Label {
	return config.RemoveLabel(list, item)
}

func pointer[T any](x T) *T {
	return &x
}

func intPtr(v int) *int {
	if v != 0 {
		return &v
	}
	return nil
}

func stringPtr(v string) *string {
	if v != "" {
		return &v
	}
	return nil
}

// HashCodeString hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func HashCodeString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func expandStringList(list []any) (res []string) {
	return config.ExpandStringList(list)
}

func GetDataSourcesMap() map[string]*schema.Resource {
	dataSourcesMap := map[string]*schema.Resource{
		"mongodbatlas_advanced_cluster":                  dataSourceMongoDBAtlasAdvancedCluster(),
		"mongodbatlas_advanced_clusters":                 dataSourceMongoDBAtlasAdvancedClusters(),
		"mongodbatlas_custom_db_role":                    dataSourceMongoDBAtlasCustomDBRole(),
		"mongodbatlas_custom_db_roles":                   dataSourceMongoDBAtlasCustomDBRoles(),
		"mongodbatlas_api_key":                           dataSourceMongoDBAtlasAPIKey(),
		"mongodbatlas_api_keys":                          dataSourceMongoDBAtlasAPIKeys(),
		"mongodbatlas_access_list_api_key":               dataSourceMongoDBAtlasAccessListAPIKey(),
		"mongodbatlas_access_list_api_keys":              dataSourceMongoDBAtlasAccessListAPIKeys(),
		"mongodbatlas_project_api_key":                   dataSourceMongoDBAtlasProjectAPIKey(),
		"mongodbatlas_project_api_keys":                  dataSourceMongoDBAtlasProjectAPIKeys(),
		"mongodbatlas_roles_org_id":                      dataSourceMongoDBAtlasOrgID(),
		"mongodbatlas_cluster":                           dataSourceMongoDBAtlasCluster(),
		"mongodbatlas_clusters":                          dataSourceMongoDBAtlasClusters(),
		"mongodbatlas_network_container":                 dataSourceMongoDBAtlasNetworkContainer(),
		"mongodbatlas_network_containers":                dataSourceMongoDBAtlasNetworkContainers(),
		"mongodbatlas_network_peering":                   dataSourceMongoDBAtlasNetworkPeering(),
		"mongodbatlas_network_peerings":                  dataSourceMongoDBAtlasNetworkPeerings(),
		"mongodbatlas_maintenance_window":                dataSourceMongoDBAtlasMaintenanceWindow(),
		"mongodbatlas_auditing":                          dataSourceMongoDBAtlasAuditing(),
		"mongodbatlas_team":                              dataSourceMongoDBAtlasTeam(),
		"mongodbatlas_teams":                             dataSourceMongoDBAtlasTeam(),
		"mongodbatlas_global_cluster_config":             dataSourceMongoDBAtlasGlobalCluster(),
		"mongodbatlas_x509_authentication_database_user": dataSourceMongoDBAtlasX509AuthDBUser(),
		"mongodbatlas_private_endpoint_regional_mode":    dataSourceMongoDBAtlasPrivateEndpointRegionalMode(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive":  dataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchive(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archives": dataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchives(),
		"mongodbatlas_privatelink_endpoint":                                         dataSourceMongoDBAtlasPrivateLinkEndpoint(),
		"mongodbatlas_privatelink_endpoint_service":                                 dataSourceMongoDBAtlasPrivateEndpointServiceLink(),
		"mongodbatlas_privatelink_endpoint_service_serverless":                      dataSourceMongoDBAtlasPrivateLinkEndpointServerless(),
		"mongodbatlas_privatelink_endpoints_service_serverless":                     dataSourceMongoDBAtlasPrivateLinkEndpointsServiceServerless(),
		"mongodbatlas_cloud_backup_schedule":                                        dataSourceMongoDBAtlasCloudBackupSchedule(),
		"mongodbatlas_third_party_integrations":                                     dataSourceMongoDBAtlasThirdPartyIntegrations(),
		"mongodbatlas_third_party_integration":                                      dataSourceMongoDBAtlasThirdPartyIntegration(),
		"mongodbatlas_cloud_provider_access":                                        dataSourceMongoDBAtlasCloudProviderAccessList(),
		"mongodbatlas_cloud_provider_access_setup":                                  dataSourceMongoDBAtlasCloudProviderAccessSetup(),
		"mongodbatlas_custom_dns_configuration_cluster_aws":                         dataSourceMongoDBAtlasCustomDNSConfigurationAWS(),
		"mongodbatlas_online_archive":                                               dataSourceMongoDBAtlasOnlineArchive(),
		"mongodbatlas_online_archives":                                              dataSourceMongoDBAtlasOnlineArchives(),
		"mongodbatlas_ldap_configuration":                                           dataSourceMongoDBAtlasLDAPConfiguration(),
		"mongodbatlas_ldap_verify":                                                  dataSourceMongoDBAtlasLDAPVerify(),
		"mongodbatlas_search_index":                                                 dataSourceMongoDBAtlasSearchIndex(),
		"mongodbatlas_search_indexes":                                               dataSourceMongoDBAtlasSearchIndexes(),
		"mongodbatlas_data_lake_pipeline_run":                                       dataSourceMongoDBAtlasDataLakePipelineRun(),
		"mongodbatlas_data_lake_pipeline_runs":                                      dataSourceMongoDBAtlasDataLakePipelineRuns(),
		"mongodbatlas_data_lake_pipeline":                                           dataSourceMongoDBAtlasDataLakePipeline(),
		"mongodbatlas_data_lake_pipelines":                                          dataSourceMongoDBAtlasDataLakePipelines(),
		"mongodbatlas_event_trigger":                                                dataSourceMongoDBAtlasEventTrigger(),
		"mongodbatlas_event_triggers":                                               dataSourceMongoDBAtlasEventTriggers(),
		"mongodbatlas_project_invitation":                                           dataSourceMongoDBAtlasProjectInvitation(),
		"mongodbatlas_org_invitation":                                               dataSourceMongoDBAtlasOrgInvitation(),
		"mongodbatlas_organization":                                                 dataSourceMongoDBAtlasOrganization(),
		"mongodbatlas_organizations":                                                dataSourceMongoDBAtlasOrganizations(),
		"mongodbatlas_cloud_backup_snapshot":                                        dataSourceMongoDBAtlasCloudBackupSnapshot(),
		"mongodbatlas_cloud_backup_snapshots":                                       dataSourceMongoDBAtlasCloudBackupSnapshots(),
		"mongodbatlas_backup_compliance_policy":                                     dataSourceMongoDBAtlasBackupCompliancePolicy(),
		"mongodbatlas_cloud_backup_snapshot_restore_job":                            dataSourceMongoDBAtlasCloudBackupSnapshotRestoreJob(),
		"mongodbatlas_cloud_backup_snapshot_restore_jobs":                           dataSourceMongoDBAtlasCloudBackupSnapshotRestoreJobs(),
		"mongodbatlas_cloud_backup_snapshot_export_bucket":                          datasourceMongoDBAtlasCloudBackupSnapshotExportBucket(),
		"mongodbatlas_cloud_backup_snapshot_export_buckets":                         datasourceMongoDBAtlasCloudBackupSnapshotExportBuckets(),
		"mongodbatlas_cloud_backup_snapshot_export_job":                             datasourceMongoDBAtlasCloudBackupSnapshotExportJob(),
		"mongodbatlas_cloud_backup_snapshot_export_jobs":                            datasourceMongoDBAtlasCloudBackupSnapshotExportJobs(),
		"mongodbatlas_federated_settings":                                           dataSourceMongoDBAtlasFederatedSettings(),
		"mongodbatlas_federated_settings_identity_provider":                         dataSourceMongoDBAtlasFederatedSettingsIdentityProvider(),
		"mongodbatlas_federated_settings_identity_providers":                        dataSourceMongoDBAtlasFederatedSettingsIdentityProviders(),
		"mongodbatlas_federated_settings_org_config":                                dataSourceMongoDBAtlasFederatedSettingsOrganizationConfig(),
		"mongodbatlas_federated_settings_org_configs":                               dataSourceMongoDBAtlasFederatedSettingsOrganizationConfigs(),
		"mongodbatlas_federated_settings_org_role_mapping":                          dataSourceMongoDBAtlasFederatedSettingsOrganizationRoleMapping(),
		"mongodbatlas_federated_settings_org_role_mappings":                         dataSourceMongoDBAtlasFederatedSettingsOrganizationRoleMappings(),
		"mongodbatlas_federated_database_instance":                                  dataSourceMongoDBAtlasFederatedDatabaseInstance(),
		"mongodbatlas_federated_database_instances":                                 dataSourceMongoDBAtlasFederatedDatabaseInstances(),
		"mongodbatlas_federated_query_limit":                                        dataSourceMongoDBAtlasFederatedDatabaseQueryLimit(),
		"mongodbatlas_federated_query_limits":                                       dataSourceMongoDBAtlasFederatedDatabaseQueryLimits(),
		"mongodbatlas_serverless_instance":                                          dataSourceMongoDBAtlasServerlessInstance(),
		"mongodbatlas_serverless_instances":                                         dataSourceMongoDBAtlasServerlessInstances(),
		"mongodbatlas_cluster_outage_simulation":                                    dataSourceMongoDBAtlasClusterOutageSimulation(),
		"mongodbatlas_shared_tier_restore_job":                                      dataSourceMongoDBAtlasCloudSharedTierRestoreJob(),
		"mongodbatlas_shared_tier_restore_jobs":                                     dataSourceMongoDBAtlasCloudSharedTierRestoreJobs(),
		"mongodbatlas_shared_tier_snapshot":                                         dataSourceMongoDBAtlasSharedTierSnapshot(),
		"mongodbatlas_shared_tier_snapshots":                                        dataSourceMongoDBAtlasSharedTierSnapshots(),
	}
	return dataSourcesMap
}

func GetResourcesMap() map[string]*schema.Resource {
	resourcesMap := map[string]*schema.Resource{
		"mongodbatlas_advanced_cluster":                  ResourceMongoDBAtlasAdvancedCluster(),
		"mongodbatlas_api_key":                           resourceMongoDBAtlasAPIKey(),
		"mongodbatlas_access_list_api_key":               resourceMongoDBAtlasAccessListAPIKey(),
		"mongodbatlas_project_api_key":                   resourceMongoDBAtlasProjectAPIKey(),
		"mongodbatlas_custom_db_role":                    resourceMongoDBAtlasCustomDBRole(),
		"mongodbatlas_cluster":                           resourceMongoDBAtlasCluster(),
		"mongodbatlas_network_container":                 resourceMongoDBAtlasNetworkContainer(),
		"mongodbatlas_network_peering":                   resourceMongoDBAtlasNetworkPeering(),
		"mongodbatlas_maintenance_window":                resourceMongoDBAtlasMaintenanceWindow(),
		"mongodbatlas_auditing":                          resourceMongoDBAtlasAuditing(),
		"mongodbatlas_team":                              resourceMongoDBAtlasTeam(),
		"mongodbatlas_teams":                             resourceMongoDBAtlasTeam(),
		"mongodbatlas_global_cluster_config":             resourceMongoDBAtlasGlobalCluster(),
		"mongodbatlas_x509_authentication_database_user": resourceMongoDBAtlasX509AuthDBUser(),
		"mongodbatlas_private_endpoint_regional_mode":    resourceMongoDBAtlasPrivateEndpointRegionalMode(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive": resourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchive(),
		"mongodbatlas_privatelink_endpoint":                                        resourceMongoDBAtlasPrivateLinkEndpoint(),
		"mongodbatlas_privatelink_endpoint_serverless":                             resourceMongoDBAtlasPrivateLinkEndpointServerless(),
		"mongodbatlas_privatelink_endpoint_service":                                resourceMongoDBAtlasPrivateEndpointServiceLink(),
		"mongodbatlas_privatelink_endpoint_service_serverless":                     resourceMongoDBAtlasPrivateLinkEndpointServiceServerless(),
		"mongodbatlas_third_party_integration":                                     resourceMongoDBAtlasThirdPartyIntegration(),
		"mongodbatlas_cloud_provider_access":                                       resourceMongoDBAtlasCloudProviderAccess(),
		"mongodbatlas_online_archive":                                              resourceMongoDBAtlasOnlineArchive(),
		"mongodbatlas_custom_dns_configuration_cluster_aws":                        resourceMongoDBAtlasCustomDNSConfiguration(),
		"mongodbatlas_ldap_configuration":                                          resourceMongoDBAtlasLDAPConfiguration(),
		"mongodbatlas_ldap_verify":                                                 resourceMongoDBAtlasLDAPVerify(),
		"mongodbatlas_cloud_provider_access_setup":                                 resourceMongoDBAtlasCloudProviderAccessSetup(),
		"mongodbatlas_cloud_provider_access_authorization":                         resourceMongoDBAtlasCloudProviderAccessAuthorization(),
		"mongodbatlas_search_index":                                                resourceMongoDBAtlasSearchIndex(),
		"mongodbatlas_data_lake_pipeline":                                          resourceMongoDBAtlasDataLakePipeline(),
		"mongodbatlas_event_trigger":                                               resourceMongoDBAtlasEventTriggers(),
		"mongodbatlas_cloud_backup_schedule":                                       resourceMongoDBAtlasCloudBackupSchedule(),
		"mongodbatlas_project_invitation":                                          resourceMongoDBAtlasProjectInvitation(),
		"mongodbatlas_org_invitation":                                              resourceMongoDBAtlasOrgInvitation(),
		"mongodbatlas_organization":                                                resourceMongoDBAtlasOrganization(),
		"mongodbatlas_cloud_backup_snapshot":                                       resourceMongoDBAtlasCloudBackupSnapshot(),
		"mongodbatlas_backup_compliance_policy":                                    resourceMongoDBAtlasBackupCompliancePolicy(),
		"mongodbatlas_cloud_backup_snapshot_restore_job":                           resourceMongoDBAtlasCloudBackupSnapshotRestoreJob(),
		"mongodbatlas_cloud_backup_snapshot_export_bucket":                         resourceMongoDBAtlasCloudBackupSnapshotExportBucket(),
		"mongodbatlas_cloud_backup_snapshot_export_job":                            resourceMongoDBAtlasCloudBackupSnapshotExportJob(),
		"mongodbatlas_federated_settings_org_config":                               resourceMongoDBAtlasFederatedSettingsOrganizationConfig(),
		"mongodbatlas_federated_settings_org_role_mapping":                         resourceMongoDBAtlasFederatedSettingsOrganizationRoleMapping(),
		"mongodbatlas_federated_settings_identity_provider":                        resourceMongoDBAtlasFederatedSettingsIdentityProvider(),
		"mongodbatlas_federated_database_instance":                                 resourceMongoDBAtlasFederatedDatabaseInstance(),
		"mongodbatlas_federated_query_limit":                                       resourceMongoDBAtlasFederatedDatabaseQueryLimit(),
		"mongodbatlas_serverless_instance":                                         resourceMongoDBAtlasServerlessInstance(),
		"mongodbatlas_cluster_outage_simulation":                                   resourceMongoDBAtlasClusterOutageSimulation(),
	}
	return resourcesMap
}
