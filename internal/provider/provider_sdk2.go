package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/accesslistapikey"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
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
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/datalakepipeline"
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
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointserverless"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointservice"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointservicedatafederationonlinearchive"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/privatelinkendpointserviceserverless"
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

type SecretData struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

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
		"mongodbatlas_teams":                             team.LegacyTeamsDataSource(),
		"mongodbatlas_global_cluster_config":             globalclusterconfig.DataSource(),
		"mongodbatlas_x509_authentication_database_user": x509authenticationdatabaseuser.DataSource(),
		"mongodbatlas_private_endpoint_regional_mode":    privateendpointregionalmode.DataSource(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive":  privatelinkendpointservicedatafederationonlinearchive.DataSource(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archives": privatelinkendpointservicedatafederationonlinearchive.PluralDataSource(),
		"mongodbatlas_privatelink_endpoint":                                         privatelinkendpoint.DataSource(),
		"mongodbatlas_privatelink_endpoint_service":                                 privatelinkendpointservice.DataSource(),
		"mongodbatlas_privatelink_endpoint_service_serverless":                      privatelinkendpointserviceserverless.DataSource(),
		"mongodbatlas_privatelink_endpoints_service_serverless":                     privatelinkendpointserviceserverless.PluralDataSource(),
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
		"mongodbatlas_data_lake_pipeline_run":                                       datalakepipeline.DataSourceRun(),
		"mongodbatlas_data_lake_pipeline_runs":                                      datalakepipeline.PluralDataSourceRun(),
		"mongodbatlas_data_lake_pipeline":                                           datalakepipeline.DataSource(),
		"mongodbatlas_data_lake_pipelines":                                          datalakepipeline.PluralDataSource(),
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
	if !config.PreviewProviderV2AdvancedCluster() {
		dataSourcesMap["mongodbatlas_advanced_cluster"] = advancedcluster.DataSource()
		dataSourcesMap["mongodbatlas_advanced_clusters"] = advancedcluster.PluralDataSource()
	}
	return dataSourcesMap
}

func getResourcesMap() map[string]*schema.Resource {
	resourcesMap := map[string]*schema.Resource{
		"mongodbatlas_api_key":                           apikey.Resource(),
		"mongodbatlas_access_list_api_key":               accesslistapikey.Resource(),
		"mongodbatlas_project_api_key":                   projectapikey.Resource(),
		"mongodbatlas_custom_db_role":                    config.NewAnalyticsResourceSDKv2(customdbrole.Resource(), "custom_db_role"),
		"mongodbatlas_cluster":                           cluster.Resource(),
		"mongodbatlas_network_container":                 networkcontainer.Resource(),
		"mongodbatlas_network_peering":                   networkpeering.Resource(),
		"mongodbatlas_maintenance_window":                maintenancewindow.Resource(),
		"mongodbatlas_auditing":                          auditing.Resource(),
		"mongodbatlas_team":                              team.Resource(),
		"mongodbatlas_teams":                             team.LegacyTeamsResource(),
		"mongodbatlas_global_cluster_config":             globalclusterconfig.Resource(),
		"mongodbatlas_x509_authentication_database_user": x509authenticationdatabaseuser.Resource(),
		"mongodbatlas_private_endpoint_regional_mode":    privateendpointregionalmode.Resource(),
		"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive": privatelinkendpointservicedatafederationonlinearchive.Resource(),
		"mongodbatlas_privatelink_endpoint":                                        privatelinkendpoint.Resource(),
		"mongodbatlas_privatelink_endpoint_serverless":                             privatelinkendpointserverless.Resource(),
		"mongodbatlas_privatelink_endpoint_service":                                privatelinkendpointservice.Resource(),
		"mongodbatlas_privatelink_endpoint_service_serverless":                     privatelinkendpointserviceserverless.Resource(),
		"mongodbatlas_third_party_integration":                                     thirdpartyintegration.Resource(),
		"mongodbatlas_online_archive":                                              onlinearchive.Resource(),
		"mongodbatlas_custom_dns_configuration_cluster_aws":                        customdnsconfigurationclusteraws.Resource(),
		"mongodbatlas_ldap_configuration":                                          ldapconfiguration.Resource(),
		"mongodbatlas_ldap_verify":                                                 ldapverify.Resource(),
		"mongodbatlas_cloud_provider_access_setup":                                 cloudprovideraccess.ResourceSetup(),
		"mongodbatlas_cloud_provider_access_authorization":                         cloudprovideraccess.ResourceAuthorization(),
		"mongodbatlas_search_index":                                                searchindex.Resource(),
		"mongodbatlas_data_lake_pipeline":                                          datalakepipeline.Resource(),
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
	if !config.PreviewProviderV2AdvancedCluster() {
		resourcesMap["mongodbatlas_advanced_cluster"] = advancedcluster.Resource()
	}
	return resourcesMap
}

func providerConfigure(provider *schema.Provider) func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		diagnostics := setDefaultsAndValidations(d)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		cfg := config.Config{
			PublicKey:        d.Get("public_key").(string),
			PrivateKey:       d.Get("private_key").(string),
			BaseURL:          d.Get("base_url").(string),
			RealmBaseURL:     d.Get("realm_base_url").(string),
			TerraformVersion: provider.TerraformVersion,
		}

		assumeRoleValue, ok := d.GetOk("assume_role")
		awsRoleDefined := ok && len(assumeRoleValue.([]any)) > 0 && assumeRoleValue.([]any)[0] != nil
		if awsRoleDefined {
			cfg.AssumeRole = expandAssumeRole(assumeRoleValue.([]any)[0].(map[string]any))
			secret := d.Get("secret_name").(string)
			region := conversion.MongoDBRegionToAWSRegion(d.Get("region").(string))
			awsAccessKeyID := d.Get("aws_access_key_id").(string)
			awsSecretAccessKey := d.Get("aws_secret_access_key").(string)
			awsSessionToken := d.Get("aws_session_token").(string)
			endpoint := d.Get("sts_endpoint").(string)
			var err error
			cfg, err = configureCredentialsSTS(&cfg, secret, region, awsAccessKeyID, awsSecretAccessKey, awsSessionToken, endpoint)
			if err != nil {
				return nil, append(diagnostics, diag.FromErr(err)...)
			}
		}

		client, err := cfg.NewClient(ctx)
		if err != nil {
			return nil, append(diagnostics, diag.FromErr(err)...)
		}
		return client, diagnostics
	}
}

func setDefaultsAndValidations(d *schema.ResourceData) diag.Diagnostics {
	diagnostics := []diag.Diagnostic{}

	mongodbgovCloud := conversion.Pointer(d.Get("is_mongodbgov_cloud").(bool))
	if *mongodbgovCloud {
		if !isGovBaseURLConfiguredForSDK2Provider(d) {
			if err := d.Set("base_url", MongodbGovCloudURL); err != nil {
				return append(diagnostics, diag.FromErr(err)...)
			}
		}
	}

	if err := setValueFromConfigOrEnv(d, "base_url", []string{
		"MONGODB_ATLAS_BASE_URL",
		"MCLI_OPS_MANAGER_URL",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	awsRoleDefined := false
	assumeRoles := d.Get("assume_role").([]any)
	if len(assumeRoles) == 0 {
		roleArn := MultiEnvDefaultFunc([]string{
			"ASSUME_ROLE_ARN",
			"TF_VAR_ASSUME_ROLE_ARN",
		}, "").(string)
		if roleArn != "" {
			awsRoleDefined = true
			if err := d.Set("assume_role", []map[string]any{{"role_arn": roleArn}}); err != nil {
				return append(diagnostics, diag.FromErr(err)...)
			}
		}
	} else {
		awsRoleDefined = true
	}

	if err := setValueFromConfigOrEnv(d, "public_key", []string{
		"MONGODB_ATLAS_PUBLIC_API_KEY",
		"MONGODB_ATLAS_PUBLIC_KEY",
		"MCLI_PUBLIC_API_KEY",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}
	if d.Get("public_key").(string) == "" && !awsRoleDefined {
		diagnostics = append(diagnostics, diag.Diagnostic{Severity: diag.Warning, Summary: MissingAuthAttrError})
	}

	if err := setValueFromConfigOrEnv(d, "private_key", []string{
		"MONGODB_ATLAS_PRIVATE_API_KEY",
		"MONGODB_ATLAS_PRIVATE_KEY",
		"MCLI_PRIVATE_API_KEY",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	if d.Get("private_key").(string) == "" && !awsRoleDefined {
		diagnostics = append(diagnostics, diag.Diagnostic{Severity: diag.Warning, Summary: MissingAuthAttrError})
	}

	if err := setValueFromConfigOrEnv(d, "realm_base_url", []string{
		"MONGODB_REALM_BASE_URL",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	if err := setValueFromConfigOrEnv(d, "region", []string{
		"AWS_REGION",
		"TF_VAR_AWS_REGION",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	if err := setValueFromConfigOrEnv(d, "sts_endpoint", []string{
		"STS_ENDPOINT",
		"TF_VAR_STS_ENDPOINT",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	if err := setValueFromConfigOrEnv(d, "aws_access_key_id", []string{
		"AWS_ACCESS_KEY_ID",
		"TF_VAR_AWS_ACCESS_KEY_ID",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	if err := setValueFromConfigOrEnv(d, "aws_secret_access_key", []string{
		"AWS_SECRET_ACCESS_KEY",
		"TF_VAR_AWS_SECRET_ACCESS_KEY",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	if err := setValueFromConfigOrEnv(d, "secret_name", []string{
		"SECRET_NAME",
		"TF_VAR_SECRET_NAME",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	if err := setValueFromConfigOrEnv(d, "aws_session_token", []string{
		"AWS_SESSION_TOKEN",
		"TF_VAR_AWS_SESSION_TOKEN",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}

	return diagnostics
}

func setValueFromConfigOrEnv(d *schema.ResourceData, attrName string, envVars []string) error {
	var val = d.Get(attrName).(string)
	if val == "" {
		val = MultiEnvDefaultFunc(envVars, "").(string)
	}
	return d.Set(attrName, val)
}

// assumeRoleSchema From aws provider.go
func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"duration": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "The duration, between 15 minutes and 12 hours, of the role session. Valid time units are ns, us (or µs), ms, s, h, or m.",
					ValidateFunc: validAssumeRoleDuration,
				},
				"external_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A unique identifier that might be required when you assume a role in another account.",
					ValidateFunc: validation.All(
						validation.StringLenBetween(2, 1224),
						validation.StringMatch(regexp.MustCompile(`[\w+=,.@:/\-]*`), ""),
					),
				},
				"policy": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "IAM Policy JSON describing further restricting permissions for the IAM Role being assumed.",
					ValidateFunc: validation.StringIsJSON,
				},
				"policy_arns": {
					Type:        schema.TypeSet,
					Optional:    true,
					Description: "Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Amazon Resource Name (ARN) of an IAM Role to assume prior to making API calls.",
				},
				"session_name": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "An identifier for the assumed role session.",
					ValidateFunc: validAssumeRoleSessionName,
				},
				"source_identity": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Source identity specified by the principal assuming the role.",
					ValidateFunc: validAssumeRoleSourceIdentity,
				},
				"tags": {
					Type:        schema.TypeMap,
					Optional:    true,
					Description: "Assume role session tags.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"transitive_tag_keys": {
					Type:        schema.TypeSet,
					Optional:    true,
					Description: "Assume role session tag keys to pass to any subsequent sessions.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

var validAssumeRoleSessionName = validation.All(
	validation.StringLenBetween(2, 64),
	validation.StringMatch(regexp.MustCompile(`[\w+=,.@\-]*`), ""),
)

var validAssumeRoleSourceIdentity = validation.All(
	validation.StringLenBetween(2, 64),
	validation.StringMatch(regexp.MustCompile(`[\w+=,.@\-]*`), ""),
)

// validAssumeRoleDuration validates a string can be parsed as a valid time.Duration
// and is within a minimum of 15 minutes and maximum of 12 hours
func validAssumeRoleDuration(v any, k string) (ws []string, errorResults []error) {
	duration, err := time.ParseDuration(v.(string))

	if err != nil {
		errorResults = append(errorResults, fmt.Errorf("%q cannot be parsed as a duration: %w", k, err))
		return
	}

	if duration.Minutes() < 15 || duration.Hours() > 12 {
		errorResults = append(errorResults, fmt.Errorf("duration %q must be between 15 minutes (15m) and 12 hours (12h), inclusive", k))
	}

	return
}

func expandAssumeRole(tfMap map[string]any) *config.AssumeRole {
	if tfMap == nil {
		return nil
	}

	assumeRole := config.AssumeRole{}

	if v, ok := tfMap["duration"].(string); ok && v != "" {
		duration, _ := time.ParseDuration(v)
		assumeRole.Duration = duration
	}

	if v, ok := tfMap["external_id"].(string); ok && v != "" {
		assumeRole.ExternalID = v
	}

	if v, ok := tfMap["policy"].(string); ok && v != "" {
		assumeRole.Policy = v
	}

	if v, ok := tfMap["policy_arns"].(*schema.Set); ok && v.Len() > 0 {
		assumeRole.PolicyARNs = conversion.ExpandStringList(v.List())
	}

	if v, ok := tfMap["role_arn"].(string); ok && v != "" {
		assumeRole.RoleARN = v
	}

	if v, ok := tfMap["session_name"].(string); ok && v != "" {
		assumeRole.SessionName = v
	}

	if v, ok := tfMap["source_identity"].(string); ok && v != "" {
		assumeRole.SourceIdentity = v
	}

	if v, ok := tfMap["transitive_tag_keys"].(*schema.Set); ok && v.Len() > 0 {
		assumeRole.TransitiveTagKeys = conversion.ExpandStringList(v.List())
	}

	return &assumeRole
}

func isGovBaseURLConfiguredForSDK2Provider(d *schema.ResourceData) bool {
	return isGovBaseURLConfigured(d.Get("base_url").(string))
}
