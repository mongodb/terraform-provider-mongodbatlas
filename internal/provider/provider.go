package provider

import (
	"context"
	"log"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/alertconfiguration"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/apikeyprojectassignment"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/atlasuser"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserorgassignment"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserprojectassignment"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserteamassignment"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/controlplaneipaddresses"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/databaseuser"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrest"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrestprivateendpoint"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexrestorejob"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexsnapshot"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/mongodbemployeeaccessgrant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectipaccesslist"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/projectipaddresses"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/pushbasedlogexport"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/resourcepolicy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchdeployment"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamaccountdetails"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streaminstance"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprivatelinkendpoint"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamworkspace"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/teamprojectassignment"
	"github.com/mongodb/terraform-provider-mongodbatlas/version"
)

const (
	ProviderConfigError            = "error in configuring the provider."
	MissingAuthAttrError           = "either AWS Secrets Manager, Service Accounts or Atlas Programmatic API Keys attributes must be set"
	ProviderMetaUserAgentExtra     = "user_agent_extra"
	ProviderMetaUserAgentExtraDesc = "You can extend the user agent header for each request made by the provider to the Atlas Admin API. The Key Values will be formatted as {key}/{value}."
	ProviderMetaModuleName         = "module_name"
	ProviderMetaModuleNameDesc     = "The name of the module using the provider"
	ProviderMetaModuleVersion      = "module_version"
	ProviderMetaModuleVersionDesc  = "The version of the module using the provider"
)

type MongodbtlasProvider struct {
}

type tfModel struct {
	Region               types.String        `tfsdk:"region"`
	PrivateKey           types.String        `tfsdk:"private_key"`
	BaseURL              types.String        `tfsdk:"base_url"`
	RealmBaseURL         types.String        `tfsdk:"realm_base_url"`
	SecretName           types.String        `tfsdk:"secret_name"`
	PublicKey            types.String        `tfsdk:"public_key"`
	StsEndpoint          types.String        `tfsdk:"sts_endpoint"`
	AwsAccessKeyID       types.String        `tfsdk:"aws_access_key_id"`
	AwsSecretAccessKeyID types.String        `tfsdk:"aws_secret_access_key"`
	AwsSessionToken      types.String        `tfsdk:"aws_session_token"`
	ClientID             types.String        `tfsdk:"client_id"`
	ClientSecret         types.String        `tfsdk:"client_secret"`
	AccessToken          types.String        `tfsdk:"access_token"`
	AssumeRole           []tfAssumeRoleModel `tfsdk:"assume_role"`
	IsMongodbGovCloud    types.Bool          `tfsdk:"is_mongodbgov_cloud"`
}

type tfAssumeRoleModel struct {
	RoleARN types.String `tfsdk:"role_arn"`
}

func (p *MongodbtlasProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mongodbatlas"
	resp.Version = version.ProviderVersion
}

func (p *MongodbtlasProvider) MetaSchema(ctx context.Context, req provider.MetaSchemaRequest, resp *provider.MetaSchemaResponse) {
	resp.Schema = metaschema.Schema{
		Attributes: map[string]metaschema.Attribute{
			ProviderMetaModuleName: metaschema.StringAttribute{
				Description: ProviderMetaModuleNameDesc,
				Optional:    true,
			},
			ProviderMetaModuleVersion: metaschema.StringAttribute{
				Description: ProviderMetaModuleVersionDesc,
				Optional:    true,
			},
			ProviderMetaUserAgentExtra: metaschema.MapAttribute{
				Description: ProviderMetaUserAgentExtraDesc,
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (p *MongodbtlasProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"assume_role": fwAssumeRoleSchema,
		},
		Attributes: map[string]schema.Attribute{
			"public_key": schema.StringAttribute{
				Optional:    true,
				Description: "MongoDB Atlas Programmatic Public Key",
			},
			"private_key": schema.StringAttribute{
				Optional:    true,
				Description: "MongoDB Atlas Programmatic Private Key",
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "MongoDB Atlas Base URL",
			},
			"realm_base_url": schema.StringAttribute{
				Optional:    true,
				Description: "MongoDB Realm Base URL",
			},
			"is_mongodbgov_cloud": schema.BoolAttribute{
				Optional:    true,
				Description: "MongoDB Atlas Base URL default to gov",
			},
			"secret_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of secret stored in AWS Secret Manager.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "Region where secret is stored as part of AWS Secret Manager.",
			},
			"sts_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "AWS Security Token Service endpoint. Required for cross-AWS region or cross-AWS account secrets.",
			},
			"aws_access_key_id": schema.StringAttribute{
				Optional:    true,
				Description: "AWS API Access Key.",
			},
			"aws_secret_access_key": schema.StringAttribute{
				Optional:    true,
				Description: "AWS API Access Secret Key.",
			},
			"aws_session_token": schema.StringAttribute{
				Optional:    true,
				Description: "AWS Security Token Service provided session token.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "MongoDB Atlas Client ID for Service Account.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Description: "MongoDB Atlas Client Secret for Service Account.",
			},
			"access_token": schema.StringAttribute{
				Optional:    true,
				Description: "MongoDB Atlas Access Token for Service Account.",
			},
		},
	}
}

var fwAssumeRoleSchema = schema.ListNestedBlock{
	Validators: []validator.List{listvalidator.SizeAtMost(1)},
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"role_arn": schema.StringAttribute{
				Optional:    true,
				Description: "Amazon Resource Name (ARN) of an IAM Role to assume prior to making API calls.",
			},
		},
	},
}

func (p *MongodbtlasProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	providerVars := getProviderVars(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	c, err := config.GetCredentials(ctx, providerVars, config.NewEnvVars(), getAWSCredentials)
	if err != nil {
		resp.Diagnostics.AddError("Error getting credentials for provider", err.Error())
		return
	}
	if c.Errors() != "" {
		resp.Diagnostics.AddError("Error getting credentials for provider", c.Errors())
		return
	}
	if c.Warnings() != "" {
		resp.Diagnostics.AddWarning("Warning getting credentials for provider", c.Warnings())
	}
	client, err := config.NewClient(c, req.TerraformVersion)
	if err != nil {
		resp.Diagnostics.AddError("Error initializing provider", err.Error())
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func getProviderVars(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) *config.Vars {
	var data tfModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return nil
	}
	assumeRoleARN := ""
	if len(data.AssumeRole) > 0 {
		assumeRoleARN = data.AssumeRole[0].RoleARN.ValueString()
	}
	baseURL := applyGovBaseURLIfNeeded(data.BaseURL.ValueString(), data.IsMongodbGovCloud.ValueBool())
	return &config.Vars{
		AccessToken:        data.AccessToken.ValueString(),
		ClientID:           data.ClientID.ValueString(),
		ClientSecret:       data.ClientSecret.ValueString(),
		PublicKey:          data.PublicKey.ValueString(),
		PrivateKey:         data.PrivateKey.ValueString(),
		BaseURL:            baseURL,
		RealmBaseURL:       data.RealmBaseURL.ValueString(),
		AWSAssumeRoleARN:   assumeRoleARN,
		AWSSecretName:      data.SecretName.ValueString(),
		AWSRegion:          data.Region.ValueString(),
		AWSAccessKeyID:     data.AwsAccessKeyID.ValueString(),
		AWSSecretAccessKey: data.AwsSecretAccessKeyID.ValueString(),
		AWSSessionToken:    data.AwsSessionToken.ValueString(),
		AWSEndpoint:        data.StsEndpoint.ValueString(),
	}
}

func applyGovBaseURLIfNeeded(providerBaseURL string, providerIsMongodbGovCloud bool) string {
	const govURL = "https://cloud.mongodbgov.com"
	govAdditionalURLs := []string{
		"https://cloud-dev.mongodbgov.com",
		"https://cloud-qa.mongodbgov.com",
	}
	if providerIsMongodbGovCloud && !slices.Contains(govAdditionalURLs, config.NormalizeBaseURL(providerBaseURL)) {
		return govURL
	}
	return providerBaseURL
}

func (p *MongodbtlasProvider) DataSources(context.Context) []func() datasource.DataSource {
	dataSources := []func() datasource.DataSource{
		project.DataSource,
		project.PluralDataSource,
		databaseuser.DataSource,
		databaseuser.PluralDataSource,
		alertconfiguration.DataSource,
		alertconfiguration.PluralDataSource,
		projectipaccesslist.DataSource,
		atlasuser.DataSource,
		atlasuser.PluralDataSource,
		searchdeployment.DataSource,
		pushbasedlogexport.DataSource,
		streaminstance.DataSource,
		streaminstance.PluralDataSource,
		streamworkspace.DataSource,
		streamworkspace.PluralDataSource,
		streamconnection.DataSource,
		streamconnection.PluralDataSource,
		controlplaneipaddresses.DataSource,
		projectipaddresses.DataSource,
		streamprocessor.DataSource,
		streamprocessor.PluralDataSource,
		encryptionatrest.DataSource,
		encryptionatrestprivateendpoint.DataSource,
		encryptionatrestprivateendpoint.PluralDataSource,
		mongodbemployeeaccessgrant.DataSource,
		streamaccountdetails.DataSource,
		streamprivatelinkendpoint.DataSource,
		streamprivatelinkendpoint.PluralDataSource,
		flexcluster.DataSource,
		flexcluster.PluralDataSource,
		flexsnapshot.DataSource,
		flexsnapshot.PluralDataSource,
		flexrestorejob.DataSource,
		flexrestorejob.PluralDataSource,
		resourcepolicy.DataSource,
		resourcepolicy.PluralDataSource,
		clouduserorgassignment.DataSource,
		clouduserprojectassignment.DataSource,
		clouduserteamassignment.DataSource,
		teamprojectassignment.DataSource,
		apikeyprojectassignment.DataSource,
		apikeyprojectassignment.PluralDataSource,
		advancedcluster.DataSource,
		advancedcluster.PluralDataSource,
	}
	return dataSources
}

func (p *MongodbtlasProvider) Resources(context.Context) []func() resource.Resource {
	resources := []func() resource.Resource{
		project.Resource,
		encryptionatrest.Resource,
		databaseuser.Resource,
		alertconfiguration.Resource,
		projectipaccesslist.Resource,
		searchdeployment.Resource,
		pushbasedlogexport.Resource,
		streaminstance.Resource,
		streamworkspace.Resource,
		streamconnection.Resource,
		streamprocessor.Resource,
		encryptionatrestprivateendpoint.Resource,
		mongodbemployeeaccessgrant.Resource,
		streamprivatelinkendpoint.Resource,
		flexcluster.Resource,
		resourcepolicy.Resource,
		clouduserorgassignment.Resource,
		apikeyprojectassignment.Resource,
		clouduserprojectassignment.Resource,
		teamprojectassignment.Resource,
		clouduserteamassignment.Resource,
		advancedcluster.Resource,
	}
	return resources
}

func NewFrameworkProvider() provider.Provider {
	return &MongodbtlasProvider{}
}

func MuxProviderFactory() func() tfprotov6.ProviderServer {
	v2Provider := NewSdkV2Provider()
	newProvider := NewFrameworkProvider()
	ctx := context.Background()
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, v2Provider.GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}
	muxServer, err := tf6muxserver.NewMuxServer(ctx,
		func() tfprotov6.ProviderServer { return upgradedSdkProvider },
		providerserver.NewProtocol6(newProvider),
	)
	if err != nil {
		log.Fatal(err)
	}
	return muxServer.ProviderServer
}
