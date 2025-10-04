package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/accesslistapikey"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/apikey"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/auditing"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/backupcompliancepolicy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupschedule"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupsnapshot"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupsnapshotexportbucket"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupsnapshotexportjob"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudbackupsnapshotrestorejob"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cloudprovideraccess"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clusteroutagesimulation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/customdbrole"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/customdnsconfigurationclusteraws"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/eventtrigger"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federateddatabaseinstance"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedquerylimit"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsidentityprovider"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsorgconfig"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsorgrolemapping"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/globalclusterconfig"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/ldapconfiguration"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/ldapverify"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/maintenancewindow"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/networkcontainer"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/networkpeering"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/onlinearchive"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/organization"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/orginvitation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privateendpointregionalmode"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpoint"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointservice"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointservicedatafederationonlinearchive"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectapikey"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectinvitation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/rolesorgid"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchindex"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/serverlessinstance"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/sharedtier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/team"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/thirdpartyintegration"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/x509authenticationdatabaseuser"
)

// NewSdkV2Provider returns the provider to be use by the code.
func NewSdkV2Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MongoDB Atlas Programmatic Public Key",
			},
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MongoDB Atlas Programmatic Private Key",
				Sensitive:   true,
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MongoDB Atlas Base URL",
			},
			"realm_base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MongoDB Realm Base URL",
			},
			"is_mongodbgov_cloud": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "MongoDB Atlas Base URL default to gov",
			},
			"assume_role": assumeRoleSchema(),
			"secret_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of secret stored in AWS Secret Manager.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region where secret is stored as part of AWS Secret Manager.",
			},
			"sts_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Security Token Service endpoint. Required for cross-AWS region or cross-AWS account secrets.",
			},
			"aws_access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS API Access Key.",
			},
			"aws_secret_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS API Access Secret Key.",
			},
			"aws_session_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Security Token Service provided session token.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MongoDB Atlas Client ID for Service Account.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MongoDB Atlas Client Secret for Service Account.",
			},
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MongoDB Atlas Access Token for Service Account.",
			},
		},
		DataSourcesMap: getDataSourcesMap(),
		ResourcesMap:   getResourcesMap(),
	}
	provider.ConfigureContextFunc = providerConfigure(provider)
	provider.ProviderMetaSchema = map[string]*schema.Schema{
		ProviderMetaModuleName: {
			Type:        schema.TypeString,
			Description: ProviderMetaModuleNameDesc,
			Optional:    true,
		},
		ProviderMetaModuleVersion: {
			Type:        schema.TypeString,
			Description: ProviderMetaModuleVersionDesc,
			Optional:    true,
		},
		ProviderMetaUserAgentExtra: {
			Type:        schema.TypeMap,
			Elem:        schema.TypeString,
			Description: ProviderMetaUserAgentExtraDesc,
			Optional:    true,
		},
	}
	return provider
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Amazon Resource Name (ARN) of an IAM Role to assume prior to making API calls.",
				},
			},
		},
	}
}

func getDataSourcesMap() map[string]*schema.Resource {
	dataSourcesMap := map[string]*schema.Resource{
		"mongodbatlas_custom_db_role":                    customdbrole.DataSource(),
		"mongodbatlas_custom_db_roles":                   customdbrole.PluralDataSource(),
		"mongodbatlas_api_key":                           apikey.DataSource(),
		"mongodbatlas_api_keys":                          apikey.PluralDataSource(),
		"mongodbatlas_access_list_api_key":               accesslistapikey.DataSource(),
		"mongodbatlas_access_list_api_keys":              accesslistapikey.PluralDataSource(),
		"mongodbatlas_project_api_key":                   projectapikey.DataSource(),
		"mongodbatlas_project_api_keys":                  projectapikey.PluralDataSource(),
		"mongodbatlas_roles_org_id":                      rolesorgid.DataSource(),
		"mongodbatlas_cluster":                           cluster.DataSource(),
		"mongodbatlas_clusters":                          cluster.PluralDataSource(),
		"mongodbatlas_network_container":                 networkcontainer.DataSource(),
		"mongodbatlas_network_containers":                networkcontainer.PluralDataSource(),
		"mongodbatlas_network_peering":                   networkpeering.DataSource(),
		"mongodbatlas_network_peerings":                  networkpeering.PluralDataSource(),
		"mongodbatlas_maintenance_window":                maintenancewindow.DataSource(),
		"mongodbatlas_auditing":                          auditing.DataSource(),
		"mongodbatlas_team":                              team.DataSource(),
		"mongodbatlas_global_cluster_config":             globalclusterconfig.DataSource(),
		"mongodbatlas_x509_authentication_database_user": x509authenticationdatabaseuser.DataSource(),
		"mongodbatlas_private_endpoint_regional_mode":    privateendpointregionalmode.DataSource(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive":  privatelinkendpointservicedatafederationonlinearchive.DataSource(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archives": privatelinkendpointservicedatafederationonlinearchive.PluralDataSource(),
		"mongodbatlas_privatelink_endpoint":                                         privatelinkendpoint.DataSource(),
		"mongodbatlas_privatelink_endpoint_service":                                 privatelinkendpointservice.DataSource(),
		"mongodbatlas_third_party_integration":                                      thirdpartyintegration.DataSource(),
		"mongodbatlas_third_party_integrations":                                     thirdpartyintegration.PluralDataSource(),
		"mongodbatlas_cloud_provider_access_setup":                                  cloudprovideraccess.DataSourceSetup(),
		"mongodbatlas_custom_dns_configuration_cluster_aws":                         customdnsconfigurationclusteraws.DataSource(),
		"mongodbatlas_online_archive":                                               onlinearchive.DataSource(),
		"mongodbatlas_online_archives":                                              onlinearchive.PluralDataSource(),
		"mongodbatlas_ldap_configuration":                                           ldapconfiguration.DataSource(),
		"mongodbatlas_ldap_verify":                                                  ldapverify.DataSource(),
		"mongodbatlas_search_index":                                                 searchindex.DataSource(),
		"mongodbatlas_search_indexes":                                               searchindex.PluralDataSource(),
		"mongodbatlas_event_trigger":                                                eventtrigger.DataSource(),
		"mongodbatlas_event_triggers":                                               eventtrigger.PluralDataSource(),
		"mongodbatlas_project_invitation":                                           projectinvitation.DataSource(),
		"mongodbatlas_org_invitation":                                               orginvitation.DataSource(),
		"mongodbatlas_organization":                                                 organization.DataSource(),
		"mongodbatlas_organizations":                                                organization.PluralDataSource(),
		"mongodbatlas_backup_compliance_policy":                                     backupcompliancepolicy.DataSource(),
		"mongodbatlas_cloud_backup_schedule":                                        cloudbackupschedule.DataSource(),
		"mongodbatlas_cloud_backup_snapshot":                                        cloudbackupsnapshot.DataSource(),
		"mongodbatlas_cloud_backup_snapshots":                                       cloudbackupsnapshot.PluralDataSource(),
		"mongodbatlas_cloud_backup_snapshot_export_bucket":                          cloudbackupsnapshotexportbucket.DataSource(),
		"mongodbatlas_cloud_backup_snapshot_export_buckets":                         cloudbackupsnapshotexportbucket.PluralDataSource(),
		"mongodbatlas_cloud_backup_snapshot_export_job":                             cloudbackupsnapshotexportjob.DataSource(),
		"mongodbatlas_cloud_backup_snapshot_export_jobs":                            cloudbackupsnapshotexportjob.PluralDataSource(),
		"mongodbatlas_cloud_backup_snapshot_restore_job":                            cloudbackupsnapshotrestorejob.DataSource(),
		"mongodbatlas_cloud_backup_snapshot_restore_jobs":                           cloudbackupsnapshotrestorejob.PluralDataSource(),
		"mongodbatlas_federated_settings_identity_provider":                         federatedsettingsidentityprovider.DataSource(),
		"mongodbatlas_federated_settings_identity_providers":                        federatedsettingsidentityprovider.PluralDataSource(),
		"mongodbatlas_federated_settings":                                           federatedsettingsorgconfig.DataSourceSettings(),
		"mongodbatlas_federated_settings_org_config":                                federatedsettingsorgconfig.DataSource(),
		"mongodbatlas_federated_settings_org_configs":                               federatedsettingsorgconfig.PluralDataSource(),
		"mongodbatlas_federated_settings_org_role_mapping":                          federatedsettingsorgrolemapping.DataSource(),
		"mongodbatlas_federated_settings_org_role_mappings":                         federatedsettingsorgrolemapping.PluralDataSource(),
		"mongodbatlas_federated_database_instance":                                  federateddatabaseinstance.DataSource(),
		"mongodbatlas_federated_database_instances":                                 federateddatabaseinstance.PluralDataSource(),
		"mongodbatlas_federated_query_limit":                                        federatedquerylimit.DataSource(),
		"mongodbatlas_federated_query_limits":                                       federatedquerylimit.PluralDataSource(),
		"mongodbatlas_serverless_instance":                                          serverlessinstance.DataSource(),
		"mongodbatlas_serverless_instances":                                         serverlessinstance.PluralDataSource(),
		"mongodbatlas_cluster_outage_simulation":                                    clusteroutagesimulation.DataSource(),
		"mongodbatlas_shared_tier_restore_job":                                      sharedtier.DataSourceRestoreJob(),
		"mongodbatlas_shared_tier_restore_jobs":                                     sharedtier.PluralDataSourceRestoreJob(),
		"mongodbatlas_shared_tier_snapshot":                                         sharedtier.DataSourceSnapshot(),
		"mongodbatlas_shared_tier_snapshots":                                        sharedtier.PluralDataSourceSnapshot(),
	}
	return dataSourcesMap
}

func getResourcesMap() map[string]*schema.Resource {
	resourcesMap := map[string]*schema.Resource{
		"mongodbatlas_api_key":                                                     apikey.Resource(),
		"mongodbatlas_access_list_api_key":                                         accesslistapikey.Resource(),
		"mongodbatlas_project_api_key":                                             projectapikey.Resource(),
		"mongodbatlas_custom_db_role":                                              customdbrole.Resource(),
		"mongodbatlas_cluster":                                                     cluster.Resource(),
		"mongodbatlas_network_container":                                           networkcontainer.Resource(),
		"mongodbatlas_network_peering":                                             networkpeering.Resource(),
		"mongodbatlas_maintenance_window":                                          maintenancewindow.Resource(),
		"mongodbatlas_auditing":                                                    auditing.Resource(),
		"mongodbatlas_team":                                                        team.Resource(),
		"mongodbatlas_global_cluster_config":                                       globalclusterconfig.Resource(),
		"mongodbatlas_x509_authentication_database_user":                           x509authenticationdatabaseuser.Resource(),
		"mongodbatlas_private_endpoint_regional_mode":                              privateendpointregionalmode.Resource(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive": privatelinkendpointservicedatafederationonlinearchive.Resource(),
		"mongodbatlas_privatelink_endpoint":                                        privatelinkendpoint.Resource(),
		"mongodbatlas_privatelink_endpoint_service":                                privatelinkendpointservice.Resource(),
		"mongodbatlas_third_party_integration":                                     thirdpartyintegration.Resource(),
		"mongodbatlas_online_archive":                                              onlinearchive.Resource(),
		"mongodbatlas_custom_dns_configuration_cluster_aws":                        customdnsconfigurationclusteraws.Resource(),
		"mongodbatlas_ldap_configuration":                                          ldapconfiguration.Resource(),
		"mongodbatlas_ldap_verify":                                                 ldapverify.Resource(),
		"mongodbatlas_cloud_provider_access_setup":                                 cloudprovideraccess.ResourceSetup(),
		"mongodbatlas_cloud_provider_access_authorization":                         cloudprovideraccess.ResourceAuthorization(),
		"mongodbatlas_search_index":                                                searchindex.Resource(),
		"mongodbatlas_event_trigger":                                               eventtrigger.Resource(),
		"mongodbatlas_project_invitation":                                          projectinvitation.Resource(),
		"mongodbatlas_org_invitation":                                              orginvitation.Resource(),
		"mongodbatlas_organization":                                                organization.Resource(),
		"mongodbatlas_backup_compliance_policy":                                    backupcompliancepolicy.Resource(),
		"mongodbatlas_cloud_backup_schedule":                                       cloudbackupschedule.Resource(),
		"mongodbatlas_cloud_backup_snapshot":                                       cloudbackupsnapshot.Resource(),
		"mongodbatlas_cloud_backup_snapshot_export_bucket":                         cloudbackupsnapshotexportbucket.Resource(),
		"mongodbatlas_cloud_backup_snapshot_export_job":                            cloudbackupsnapshotexportjob.Resource(),
		"mongodbatlas_cloud_backup_snapshot_restore_job":                           cloudbackupsnapshotrestorejob.Resource(),
		"mongodbatlas_federated_settings_org_config":                               federatedsettingsorgconfig.Resource(),
		"mongodbatlas_federated_settings_org_role_mapping":                         federatedsettingsorgrolemapping.Resource(),
		"mongodbatlas_federated_settings_identity_provider":                        federatedsettingsidentityprovider.Resource(),
		"mongodbatlas_federated_database_instance":                                 federateddatabaseinstance.Resource(),
		"mongodbatlas_federated_query_limit":                                       federatedquerylimit.Resource(),
		"mongodbatlas_serverless_instance":                                         serverlessinstance.Resource(),
		"mongodbatlas_cluster_outage_simulation":                                   clusteroutagesimulation.Resource(),
	}
	return resourcesMap
}

func providerConfigure(provider *schema.Provider) func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		diagnostics := diag.Diagnostics{}

		providerVars := getSDKv2ProviderVars(d)

		// TODO: refactor, it's similar to the other provider

		envVars := config.NewEnvVars()

		awsCredentials, err := getAWSCredentials(envVars.GetAWS())
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("error getting AWS credentials: %w", err))
		}

		_, _ = providerVars, awsCredentials

		// TODO: chooose the credentials between AWS, SA or PAK
		client, err := config.NewClient(envVars.GetCredentials(), provider.TerraformVersion)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("error initializing provider: %w", err))
		}
		return client, diagnostics
	}

	// TODO gov look former code
}

// TODO: implement this
func getSDKv2ProviderVars(d *schema.ResourceData) *config.Vars {
	assumeRoleARN := ""
	assumeRoles := d.Get("assume_role").([]any)
	if len(assumeRoles) > 0 {
		assumeRoleARN = assumeRoles[0].(map[string]any)["role_arn"].(string)
	}
	return &config.Vars{
		AccessToken:        d.Get("access_token").(string),
		ClientID:           d.Get("client_id").(string),
		ClientSecret:       d.Get("client_secret").(string),
		PublicKey:          d.Get("public_key").(string),
		PrivateKey:         d.Get("private_key").(string),
		BaseURL:            d.Get("base_url").(string),
		RealmBaseURL:       d.Get("realm_base_url").(string),
		AWSAssumeRoleARN:   assumeRoleARN,
		AWSSecretName:      d.Get("secret_name").(string),
		AWSRegion:          d.Get("region").(string),
		AWSAccessKeyID:     d.Get("aws_access_key_id").(string),
		AWSSecretAccessKey: d.Get("aws_secret_access_key").(string),
		AWSSessionToken:    d.Get("aws_session_token").(string),
		AWSEndpoint:        d.Get("sts_endpoint").(string),
	}
}
