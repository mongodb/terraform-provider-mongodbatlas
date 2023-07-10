package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/framework/utils"
)

const (
	endPointSTSDefault                    = "https://sts.amazonaws.com"
	DeprecationMessage                    = "this resource is deprecated and will be removed in %s, please transition to %s"
	DeprecationMessageParameterToResource = "this parameter is deprecated and will be removed in %s, please transition to %s"
)

var _ provider.Provider = (*MongodbtlasProvider)(nil)

type MongodbtlasProvider struct {
}

type MongodbtlasProviderModel struct {
	PublicKey         types.String `tfsdk:"public_key"`
	PrivateKey        types.String `tfsdk:"private_key"`
	BaseURL           types.String `tfsdk:"base_url"`
	RealmBaseURL      types.String `tfsdk:"realm_base_url"`
	IsMongodbGovCloud types.Bool   `tfsdk:"is_mongodbgov_cloud"`

	// AssumeRole types.Object `tfsdk:"assume_role"`
	SecretName           types.String `tfsdk:"secret_name"`
	Region               types.String `tfsdk:"region"`
	StsEndpoint          types.String `tfsdk:"sts_endpoint"`
	AwsAccessKeyID       types.String `tfsdk:"aws_access_key_id"`
	AwsSecretAccessKeyID types.String `tfsdk:"aws_secret_access_key"`
	AwsSessionToken      types.String `tfsdk:"aws_session_token"`
}

type SecretData struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func (p *MongodbtlasProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mongodbatlas"
}

func (p *MongodbtlasProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
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
				Optional: true,
			},
			"region": schema.StringAttribute{
				Optional: true,
			},
			"sts_endpoint": schema.StringAttribute{
				Optional: true,
			},
			"aws_access_key_id": schema.StringAttribute{
				Optional: true,
				// Sensitive: true,
			},
			"aws_secret_access_key": schema.StringAttribute{
				Optional: true,
				// Sensitive: true,
			},
			"aws_session_token": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *MongodbtlasProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "configuring client")
	var (
		data    MongodbtlasProviderModel
		baseURL string
	)

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseURL = data.BaseURL.ValueString()
	mongodbgovCloud := data.IsMongodbGovCloud.ValueBool()
	if mongodbgovCloud {
		baseURL = "https://cloud.mongodbgov.com"
	}

	if data.PublicKey.ValueString() == "" {
		data.PublicKey = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"MONGODB_ATLAS_PUBLIC_KEY",
			"MCLI_PUBLIC_API_KEY",
		}, "").(string))
		if data.PublicKey.ValueString() == "" {
			resp.Diagnostics.AddError(utils.ProviderConfigError, fmt.Sprintf(utils.AttrNotSetError, "public_key"))
		}
	}

	if data.PrivateKey.ValueString() == "" {
		data.PrivateKey = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"MONGODB_ATLAS_PRIVATE_KEY",
			"MCLI_PRIVATE_API_KEY",
		}, "").(string))
		if data.PrivateKey.ValueString() == "" {
			resp.Diagnostics.AddError(utils.ProviderConfigError, fmt.Sprintf(utils.AttrNotSetError, "private_key"))
		}
	}

	if data.BaseURL.ValueString() == "" {
		data.BaseURL = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"MONGODB_ATLAS_BASE_URL",
			"MCLI_OPS_MANAGER_URL",
		}, "").(string))
	}

	if data.RealmBaseURL.ValueString() == "" {
		data.RealmBaseURL = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"MONGODB_REALM_BASE_URL",
		}, "").(string))
	}

	if data.Region.ValueString() == "" {
		data.Region = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"AWS_REGION",
			"TF_VAR_AWS_REGION",
		}, "").(string))
	}

	if data.StsEndpoint.ValueString() == "" {
		data.StsEndpoint = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"STS_ENDPOINT",
			"TF_VAR_STS_ENDPOINT",
		}, "").(string))
	}

	if data.AwsAccessKeyID.ValueString() == "" {
		data.AwsAccessKeyID = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"AWS_ACCESS_KEY_ID",
			"TF_VAR_AWS_ACCESS_KEY_ID",
		}, "").(string))
	}

	if data.AwsSessionToken.ValueString() == "" {
		data.AwsSessionToken = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"AWS_SESSION_TOKEN",
			"TF_VAR_AWS_SESSION_TOKEN",
		}, "").(string))
	}

	if data.AwsSecretAccessKeyID.ValueString() == "" {
		data.AwsSecretAccessKeyID = types.StringValue(utils.MultiEnvDefaultFunc([]string{
			"AWS_SECRET_ACCESS_KEY",
			"TF_VAR_AWS_SECRET_ACCESS_KEY",
		}, "").(string))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	config := Config{
		PublicKey:    data.PublicKey.ValueString(),
		PrivateKey:   data.PrivateKey.ValueString(),
		BaseURL:      baseURL,
		RealmBaseURL: data.RealmBaseURL.ValueString(),
	}

	client, _ := config.NewClient(ctx)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MongodbtlasProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func (p *MongodbtlasProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
		NewMongoDBAtlasProjectResource,
	}
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &MongodbtlasProvider{}
	}
}
