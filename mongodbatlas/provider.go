package mongodbatlas

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	"github.com/zclconf/go-cty/cty"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/util"
)

var (
	ProviderEnableBeta, _ = strconv.ParseBool(os.Getenv("MONGODB_ATLAS_ENABLE_BETA"))
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
		DataSourcesMap:       getDataSourcesMap(),
		ResourcesMap:         getResourcesMap(),
		ConfigureContextFunc: providerConfigure,
	}
	addBetaFeatures(provider)
	return provider
}

func getDataSourcesMap() map[string]*schema.Resource {
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

func getResourcesMap() map[string]*schema.Resource {
	resourcesMap := map[string]*schema.Resource{
		// "mongodbatlas_cluster":                           resourceMongoDBAtlasCluster(),
		"mongodbatlas_advanced_cluster":                                            resourceMongoDBAtlasAdvancedCluster(),
		"mongodbatlas_api_key":                                                     resourceMongoDBAtlasAPIKey(),
		"mongodbatlas_access_list_api_key":                                         resourceMongoDBAtlasAccessListAPIKey(),
		"mongodbatlas_project_api_key":                                             resourceMongoDBAtlasProjectAPIKey(),
		"mongodbatlas_custom_db_role":                                              resourceMongoDBAtlasCustomDBRole(),
		"mongodbatlas_network_container":                                           resourceMongoDBAtlasNetworkContainer(),
		"mongodbatlas_network_peering":                                             resourceMongoDBAtlasNetworkPeering(),
		"mongodbatlas_maintenance_window":                                          resourceMongoDBAtlasMaintenanceWindow(),
		"mongodbatlas_auditing":                                                    resourceMongoDBAtlasAuditing(),
		"mongodbatlas_team":                                                        resourceMongoDBAtlasTeam(),
		"mongodbatlas_teams":                                                       resourceMongoDBAtlasTeam(),
		"mongodbatlas_global_cluster_config":                                       resourceMongoDBAtlasGlobalCluster(),
		"mongodbatlas_x509_authentication_database_user":                           resourceMongoDBAtlasX509AuthDBUser(),
		"mongodbatlas_private_endpoint_regional_mode":                              resourceMongoDBAtlasPrivateEndpointRegionalMode(),
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

func addBetaFeatures(provider *schema.Provider) {
	if ProviderEnableBeta {
		return
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	diagnostics := setDefaultsAndValidations(d)
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	config := Config{
		PublicKey:    d.Get("public_key").(string),
		PrivateKey:   d.Get("private_key").(string),
		BaseURL:      d.Get("base_url").(string),
		RealmBaseURL: d.Get("realm_base_url").(string),
	}

	assumeRoleValue, ok := d.GetOk("assume_role")
	awsRoleDefined := ok && len(assumeRoleValue.([]any)) > 0 && assumeRoleValue.([]any)[0] != nil
	if awsRoleDefined {
		config.AssumeRole = expandAssumeRole(assumeRoleValue.([]any)[0].(map[string]any))
		secret := d.Get("secret_name").(string)
		region := util.MongoDBRegionToAWSRegion(d.Get("region").(string))
		awsAccessKeyID := d.Get("aws_access_key_id").(string)
		awsSecretAccessKey := d.Get("aws_secret_access_key").(string)
		awsSessionToken := d.Get("aws_session_token").(string)
		endpoint := d.Get("sts_endpoint").(string)
		var err error
		config, err = configureCredentialsSTS(config, secret, region, awsAccessKeyID, awsSecretAccessKey, awsSessionToken, endpoint)
		if err != nil {
			return nil, append(diagnostics, diag.FromErr(err)...)
		}
	}

	client, err := config.NewClient(ctx)
	if err != nil {
		return nil, append(diagnostics, diag.FromErr(err)...)
	}
	return client, diagnostics
}

func setDefaultsAndValidations(d *schema.ResourceData) diag.Diagnostics {
	diagnostics := []diag.Diagnostic{}

	mongodbgovCloud := pointy.Bool(d.Get("is_mongodbgov_cloud").(bool))
	if *mongodbgovCloud {
		if err := d.Set("base_url", MongodbGovCloudURL); err != nil {
			return append(diagnostics, diag.FromErr(err)...)
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
		"MONGODB_ATLAS_PUBLIC_KEY",
		"MCLI_PUBLIC_API_KEY",
	}); err != nil {
		return append(diagnostics, diag.FromErr(err)...)
	}
	if d.Get("public_key").(string) == "" && !awsRoleDefined {
		diagnostics = append(diagnostics, diag.Diagnostic{Severity: diag.Warning, Summary: MissingAuthAttrError})
	}

	if err := setValueFromConfigOrEnv(d, "private_key", []string{
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

func MultiEnvDefaultFunc(ks []string, def any) any {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return def
}

func configureCredentialsSTS(config Config, secret, region, awsAccessKeyID, awsSecretAccessKey, awsSessionToken, endpoint string) (Config, error) {
	ep, err := endpoints.GetSTSRegionalEndpoint("regional")
	if err != nil {
		log.Printf("GetSTSRegionalEndpoint error: %s", err)
		return config, err
	}

	defaultResolver := endpoints.DefaultResolver()
	stsCustResolverFn := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.StsServiceID {
			if endpoint == "" {
				return endpoints.ResolvedEndpoint{
					URL:           endPointSTSDefault,
					SigningRegion: region,
				}, nil
			}
			return endpoints.ResolvedEndpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}

		return defaultResolver.EndpointFor(service, region, optFns...)
	}

	cfg := aws.Config{
		Region:              aws.String(region),
		Credentials:         credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, awsSessionToken),
		STSRegionalEndpoint: ep,
		EndpointResolver:    endpoints.ResolverFunc(stsCustResolverFn),
	}

	sess := session.Must(session.NewSession(&cfg))

	creds := stscreds.NewCredentials(sess, config.AssumeRole.RoleARN)

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		log.Printf("Session get credentials error: %s", err)
		return config, err
	}
	_, err = creds.Get()
	if err != nil {
		log.Printf("STS get credentials error: %s", err)
		return config, err
	}
	secretString, err := secretsManagerGetSecretValue(sess, &aws.Config{Credentials: creds, Region: aws.String(region)}, secret)
	if err != nil {
		log.Printf("Get Secrets error: %s", err)
		return config, err
	}

	var secretData SecretData
	err = json.Unmarshal([]byte(secretString), &secretData)
	if err != nil {
		return config, err
	}
	if secretData.PrivateKey == "" {
		return config, fmt.Errorf("secret missing value for credential PrivateKey")
	}

	if secretData.PublicKey == "" {
		return config, fmt.Errorf("secret missing value for credential PublicKey")
	}

	config.PublicKey = secretData.PublicKey
	config.PrivateKey = secretData.PrivateKey
	return config, nil
}

func secretsManagerGetSecretValue(sess *session.Session, creds *aws.Config, secret string) (string, error) {
	svc := secretsmanager.New(sess, creds)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secret),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				log.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				log.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				log.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				log.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				log.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return "", err
	}

	return *result.SecretString, err
}

func encodeStateID(values map[string]string) string {
	encode := func(e string) string { return base64.StdEncoding.EncodeToString([]byte(e)) }
	encodedValues := make([]string, 0)

	// sort to make sure the same encoding is returned in case of same input
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		encodedValues = append(encodedValues, fmt.Sprintf("%s:%s", encode(key), encode(values[key])))
	}

	return strings.Join(encodedValues, "-")
}

func decodeStateID(stateID string) map[string]string {
	decode := func(d string) string {
		decodedString, err := base64.StdEncoding.DecodeString(d)
		if err != nil {
			log.Printf("[WARN] error decoding state ID: %s", err)
		}

		return string(decodedString)
	}
	decodedValues := make(map[string]string)
	encodedValues := strings.Split(stateID, "-")

	for _, value := range encodedValues {
		keyValue := strings.Split(value, ":")
		if len(keyValue) > 1 {
			decodedValues[decode(keyValue[0])] = decode(keyValue[1])
		}
	}

	return decodedValues
}

func valRegion(reg any, opt ...string) (string, error) {
	region, err := cast.ToStringE(reg)
	if err != nil {
		return "", err
	}

	if region == "" {
		return "", fmt.Errorf("region must be set")
	}

	/*
		We need to check if the option will be similar to network_peering word
		 (this comes in from the same resource) because network_pering resource
		 has not the standard region name pattern "US_EAST_1",
		 instead it needs the following one: "us-east-1".
	*/
	if len(opt) > 0 && strings.EqualFold("network_peering", opt[0]) {
		return strings.ToLower(strings.ReplaceAll(region, "_", "-")), nil
	}

	return strings.ReplaceAll(region, "-", "_"), nil
}

func expandStringList(list []any) (res []string) {
	for _, v := range list {
		res = append(res, v.(string))
	}

	return
}

func getEncodedID(stateID, keyPosition string) string {
	id := ""
	if !hasMultipleValues(stateID) {
		return stateID
	}

	decoded := decodeStateID(stateID)
	id = decoded[keyPosition]

	return id
}

func hasMultipleValues(value string) bool {
	if strings.Contains(value, "-") && strings.Contains(value, ":") {
		return true
	}

	return false
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

func appendBlockWithCtyValues(body *hclwrite.Body, name string, labels []string, values map[string]cty.Value) {
	if len(values) == 0 {
		return
	}

	keys := make([]string, 0, len(values))

	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	body.AppendNewline()
	block := body.AppendNewBlock(name, labels).Body()

	for _, k := range keys {
		block.SetAttributeValue(k, values[k])
	}
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

type AssumeRole struct {
	Tags              map[string]string
	RoleARN           string
	ExternalID        string
	Policy            string
	SessionName       string
	SourceIdentity    string
	PolicyARNs        []string
	TransitiveTagKeys []string
	Duration          time.Duration
}

func expandAssumeRole(tfMap map[string]any) *AssumeRole {
	if tfMap == nil {
		return nil
	}

	assumeRole := AssumeRole{}

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
		assumeRole.PolicyARNs = expandStringList(v.List())
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
		assumeRole.TransitiveTagKeys = expandStringList(v.List())
	}

	return &assumeRole
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
