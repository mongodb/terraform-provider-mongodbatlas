package mongodbatlas

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	sdkv2schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cstmvalidator "github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/validator"
	"github.com/mongodb/terraform-provider-mongodbatlas/version"
)

const (
	DeprecationMessageParameterToResource = "this parameter is deprecated and will be removed in %s, please transition to %s"
	DeprecationByDateMessageParameter     = "this parameter is deprecated and will be removed by %s"
	DeprecationByDateWithReplacement      = "this parameter is deprecated and will be removed by %s, please transition to %s"
	DeprecationMessage                    = "this resource is deprecated and will be removed in %s, please transition to %s"
	endPointSTSDefault                    = "https://sts.amazonaws.com"
	MissingAuthAttrError                  = "either Atlas Programmatic API Keys or AWS Secrets Manager attributes must be set"
	ProviderConfigError                   = "error in configuring the provider."
	AWS                                   = "AWS"
	AZURE                                 = "AZURE"
	GCP
	errorConfigureSummary = "Unexpected Resource Configure Type"
	errorConfigure        = "expected *MongoDBClient, got: %T. Please report this issue to the provider developers"
)

type MongodbtlasProvider struct{}

type tfMongodbAtlasProviderModel struct {
	AssumeRole           types.List   `tfsdk:"assume_role"`
	PublicKey            types.String `tfsdk:"public_key"`
	PrivateKey           types.String `tfsdk:"private_key"`
	BaseURL              types.String `tfsdk:"base_url"`
	RealmBaseURL         types.String `tfsdk:"realm_base_url"`
	SecretName           types.String `tfsdk:"secret_name"`
	Region               types.String `tfsdk:"region"`
	StsEndpoint          types.String `tfsdk:"sts_endpoint"`
	AwsAccessKeyID       types.String `tfsdk:"aws_access_key_id"`
	AwsSecretAccessKeyID types.String `tfsdk:"aws_secret_access_key"`
	AwsSessionToken      types.String `tfsdk:"aws_session_token"`
	IsMongodbGovCloud    types.Bool   `tfsdk:"is_mongodbgov_cloud"`
}

type tfAssumeRoleModel struct {
	PolicyARNs        types.Set    `tfsdk:"policy_arns"`
	TransitiveTagKeys types.Set    `tfsdk:"transitive_tag_keys"`
	Tags              types.Map    `tfsdk:"tags"`
	Duration          types.String `tfsdk:"duration"`
	ExternalID        types.String `tfsdk:"external_id"`
	Policy            types.String `tfsdk:"policy"`
	RoleARN           types.String `tfsdk:"role_arn"`
	SessionName       types.String `tfsdk:"session_name"`
	SourceIdentity    types.String `tfsdk:"source_identity"`
}

func (p *MongodbtlasProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mongodbatlas"
	resp.Version = version.ProviderVersion
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
		},
	}
}

var fwAssumeRoleSchema = schema.ListNestedBlock{
	Validators: []validator.List{listvalidator.SizeAtMost(1)},
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"duration": schema.StringAttribute{
				Optional:    true,
				Description: "The duration, between 15 minutes and 12 hours, of the role session. Valid time units are ns, us (or µs), ms, s, h, or m.",
				Validators: []validator.String{
					cstmvalidator.ValidDurationBetween(15, 12*60),
				},
			},
			"external_id": schema.StringAttribute{
				Optional:    true,
				Description: "A unique identifier that might be required when you assume a role in another account.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 1224),
					stringvalidator.RegexMatches(regexp.MustCompile(`[\w+=,.@:/\-]*`), ""),
				},
			},
			"policy": schema.StringAttribute{
				Optional:    true,
				Description: "IAM Policy JSON describing further restricting permissions for the IAM Role being assumed.",
				Validators: []validator.String{
					cstmvalidator.StringIsJSON(),
				},
			},
			"policy_arns": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed.",
			},
			"role_arn": schema.StringAttribute{
				Optional:    true,
				Description: "Amazon Resource Name (ARN) of an IAM Role to assume prior to making API calls.",
			},
			"session_name": schema.StringAttribute{
				Optional:    true,
				Description: "An identifier for the assumed role session.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`[\w+=,.@\-]*`), ""),
				},
			},
			"source_identity": schema.StringAttribute{
				Optional:    true,
				Description: "Source identity specified by the principal assuming the role.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 64),
					stringvalidator.RegexMatches(regexp.MustCompile(`[\w+=,.@\-]*`), ""),
				},
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Assume role session tags.",
			},
			"transitive_tag_keys": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Assume role session tag keys to pass to any subsequent sessions.",
			},
		},
	},
}

func (p *MongodbtlasProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data tfMongodbAtlasProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var assumeRoles []tfAssumeRoleModel
	data.AssumeRole.ElementsAs(ctx, &assumeRoles, true)
	awsRoleDefined := len(assumeRoles) > 0

	data = setDefaultValuesWithValidations(&data, awsRoleDefined, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	config := Config{
		PublicKey:    data.PublicKey.ValueString(),
		PrivateKey:   data.PrivateKey.ValueString(),
		BaseURL:      data.BaseURL.ValueString(),
		RealmBaseURL: data.RealmBaseURL.ValueString(),
	}

	if awsRoleDefined {
		config.AssumeRole = parseTfModel(ctx, &assumeRoles[0])
		secret := data.SecretName.ValueString()
		region := data.Region.ValueString()
		awsAccessKeyID := data.AwsAccessKeyID.ValueString()
		awsSecretAccessKey := data.AwsSecretAccessKeyID.ValueString()
		awsSessionToken := data.AwsSessionToken.ValueString()
		endpoint := data.StsEndpoint.ValueString()
		var err error
		config, err = configureCredentialsSTS(config, secret, region, awsAccessKeyID, awsSecretAccessKey, awsSessionToken, endpoint)
		if err != nil {
			resp.Diagnostics.AddError("failed to configure credentials STS", err.Error())
			return
		}
	}

	client, err := config.NewClient(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"failed to initialize a new client",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// parseTfModel extracts the values from tfAssumeRoleModel creating a new instance of our internal model AssumeRole used in Config
func parseTfModel(ctx context.Context, tfAssumeRoleModel *tfAssumeRoleModel) *AssumeRole {
	assumeRole := AssumeRole{}

	if !tfAssumeRoleModel.Duration.IsNull() {
		duration, _ := time.ParseDuration(tfAssumeRoleModel.Duration.ValueString())
		assumeRole.Duration = duration
	}

	assumeRole.ExternalID = tfAssumeRoleModel.ExternalID.ValueString()
	assumeRole.Policy = tfAssumeRoleModel.Policy.ValueString()

	if !tfAssumeRoleModel.PolicyARNs.IsNull() {
		var policiesARNs []string
		tfAssumeRoleModel.PolicyARNs.ElementsAs(ctx, &policiesARNs, true)
		assumeRole.PolicyARNs = policiesARNs
	}

	assumeRole.RoleARN = tfAssumeRoleModel.RoleARN.ValueString()
	assumeRole.SessionName = tfAssumeRoleModel.SessionName.ValueString()
	assumeRole.SourceIdentity = tfAssumeRoleModel.SourceIdentity.ValueString()

	if !tfAssumeRoleModel.TransitiveTagKeys.IsNull() {
		var transitiveTagKeys []string
		tfAssumeRoleModel.TransitiveTagKeys.ElementsAs(ctx, &transitiveTagKeys, true)
		assumeRole.TransitiveTagKeys = transitiveTagKeys
	}

	return &assumeRole
}

const MongodbGovCloudURL = "https://cloud.mongodbgov.com"

func setDefaultValuesWithValidations(data *tfMongodbAtlasProviderModel, awsRoleDefined bool, resp *provider.ConfigureResponse) tfMongodbAtlasProviderModel {
	if mongodbgovCloud := data.IsMongodbGovCloud.ValueBool(); mongodbgovCloud {
		data.BaseURL = types.StringValue(MongodbGovCloudURL)
	}
	if data.BaseURL.ValueString() == "" {
		data.BaseURL = types.StringValue(MultiEnvDefaultFunc([]string{
			"MONGODB_ATLAS_BASE_URL",
			"MCLI_OPS_MANAGER_URL",
		}, "").(string))
	}

	if data.PublicKey.ValueString() == "" {
		data.PublicKey = types.StringValue(MultiEnvDefaultFunc([]string{
			"MONGODB_ATLAS_PUBLIC_KEY",
			"MCLI_PUBLIC_API_KEY",
		}, "").(string))
		if data.PublicKey.ValueString() == "" && !awsRoleDefined {
			resp.Diagnostics.AddWarning(ProviderConfigError, MissingAuthAttrError)
		}
	}

	if data.PrivateKey.ValueString() == "" {
		data.PrivateKey = types.StringValue(MultiEnvDefaultFunc([]string{
			"MONGODB_ATLAS_PRIVATE_KEY",
			"MCLI_PRIVATE_API_KEY",
		}, "").(string))
		if data.PrivateKey.ValueString() == "" && !awsRoleDefined {
			resp.Diagnostics.AddWarning(ProviderConfigError, MissingAuthAttrError)
		}
	}

	if data.RealmBaseURL.ValueString() == "" {
		data.RealmBaseURL = types.StringValue(MultiEnvDefaultFunc([]string{
			"MONGODB_REALM_BASE_URL",
		}, "").(string))
	}

	if data.Region.ValueString() == "" {
		data.Region = types.StringValue(MultiEnvDefaultFunc([]string{
			"AWS_REGION",
			"TF_VAR_AWS_REGION",
		}, "").(string))
	}

	if data.StsEndpoint.ValueString() == "" {
		data.StsEndpoint = types.StringValue(MultiEnvDefaultFunc([]string{
			"STS_ENDPOINT",
			"TF_VAR_STS_ENDPOINT",
		}, "").(string))
	}

	if data.AwsAccessKeyID.ValueString() == "" {
		data.AwsAccessKeyID = types.StringValue(MultiEnvDefaultFunc([]string{
			"AWS_ACCESS_KEY_ID",
			"TF_VAR_AWS_ACCESS_KEY_ID",
		}, "").(string))
	}

	if data.AwsSecretAccessKeyID.ValueString() == "" {
		data.AwsSecretAccessKeyID = types.StringValue(MultiEnvDefaultFunc([]string{
			"AWS_SECRET_ACCESS_KEY",
			"TF_VAR_AWS_SECRET_ACCESS_KEY",
		}, "").(string))
	}

	if data.AwsSessionToken.ValueString() == "" {
		data.AwsSessionToken = types.StringValue(MultiEnvDefaultFunc([]string{
			"AWS_SESSION_TOKEN",
			"TF_VAR_AWS_SESSION_TOKEN",
		}, "").(string))
	}

	return *data
}

func (p *MongodbtlasProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDS,
		NewProjectsDS,
		NewDatabaseUserDS,
		NewDatabaseUsersDS,
		NewAlertConfigurationDS,
		NewAlertConfigurationsDS,
		NewProjectIPAccessListDS,
		NewAtlasUserDS,
		NewAtlasUsersDS,
	}
}

func (p *MongodbtlasProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectRS,
		NewEncryptionAtRestRS,
		NewDatabaseUserRS,
		NewAlertConfigurationRS,
		NewProjectIPAccessListRS,
	}
}

func NewFrameworkProvider() provider.Provider {
	return &MongodbtlasProvider{}
}

func MuxedProviderFactory() func() tfprotov6.ProviderServer {
	return muxedProviderFactory(NewSdkV2Provider())
}

// muxedProviderFactory creates mux provider using existing sdk v2 provider passed as parameter and creating new instance of framework provider.
// Used in testing where existing sdk v2 provider has to be used.
func muxedProviderFactory(sdkV2Provider *sdkv2schema.Provider) func() tfprotov6.ProviderServer {
	fwProvider := NewFrameworkProvider()

	ctx := context.Background()
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, sdkV2Provider.GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
		providerserver.NewProtocol6(fwProvider),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}
	return muxServer.ProviderServer
}
