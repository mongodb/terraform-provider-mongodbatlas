package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/framework/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	endPointSTSDefault                    = "https://sts.amazonaws.com"
	DeprecationMessage                    = "this resource is deprecated and will be removed in %s, please transition to %s"
	DeprecationMessageParameterToResource = "this parameter is deprecated and will be removed in %s, please transition to %s"
)

// var _ provider.Provider = &MongodbtlasProvider{}

var _ provider.Provider = (*MongodbtlasProvider)(nil)

type MongodbtlasProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	// Version string
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
			},
			"aws_secret_access_key": schema.StringAttribute{
				Optional: true,
			},
			"aws_session_token": schema.StringAttribute{
				Optional: true,
			},
		},
		// TODO: validAssumeRoleDuration
		//       validation.StringIsJSON
		// Blocks: map[string]schema.Block{
		// 	// Optional: true,
		// 	"assume_role": schema.ListNestedBlock{
		// 		// MaxItems: 1,
		// 		Validators: []validator.List{
		// 			listvalidator.SizeAtMost(1),
		// 		},
		// 		NestedObject: schema.NestedBlockObject{
		// 			Attributes: map[string]schema.Attribute{
		// 				"duration": schema.StringAttribute{
		// 					Optional:    true,
		// 					Description: "The duration, between 15 minutes and 12 hours, of the role session. Valid time units are ns, us (or µs), ms, s, h, or m.",
		// 					// ValidateFunc:  validAssumeRoleDuration,
		// 					Validators: []validator.String{
		// 						stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("duration_seconds")),
		// 					},
		// 				},
		// 				"duration_seconds": schema.Int64Attribute{
		// 					Optional:           true,
		// 					DeprecationMessage: "Use assume_role.duration instead",
		// 					Description:        "The duration, in seconds, of the role session.",
		// 					Validators: []validator.Int64{
		// 						int64validator.ConflictsWith(path.MatchRelative().AtParent().AtName("duration")),
		// 						int64validator.Between(900, 43200),
		// 					},
		// 				},
		// 				"external_id": schema.StringAttribute{
		// 					Optional:    true,
		// 					Description: "A unique identifier that might be required when you assume a role in another account.",
		// 					Validators: []validator.String{
		// 						stringvalidator.LengthBetween(2, 1224),
		// 						stringvalidator.RegexMatches(regexp.MustCompile(`[\w+=,.@:/\-]*`), ""),
		// 					},
		// 				},
		// 				"policy": schema.StringAttribute{
		// 					Optional:    true,
		// 					Description: "IAM Policy JSON describing further restricting permissions for the IAM Role being assumed.",
		// 					// ValidateFunc: validation.StringIsJSON,
		// 				},
		// 				"policy_arns": schema.SetAttribute{
		// 					Optional:    true,
		// 					Description: "Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed.",
		// 					ElementType: types.StringType,
		// 				},
		// 				"role_arn": schema.StringAttribute{
		// 					Optional:    true,
		// 					Description: "Amazon Resource Name (ARN) of an IAM Role to assume prior to making API calls.",
		// 				},
		// 				"session_name": schema.StringAttribute{
		// 					Optional:    true,
		// 					Description: "An identifier for the assumed role session.",
		// 					Validators: []validator.String{
		// 						stringvalidator.LengthBetween(2, 64),
		// 						stringvalidator.RegexMatches(regexp.MustCompile(`[\w+=,.@\-]*`), ""),
		// 					},
		// 				},
		// 				"source_identity": schema.StringAttribute{
		// 					Optional:    true,
		// 					Description: "Source identity specified by the principal assuming the role.",
		// 					Validators: []validator.String{
		// 						stringvalidator.LengthBetween(2, 64),
		// 						stringvalidator.RegexMatches(regexp.MustCompile(`[\w+=,.@\-]*`), ""),
		// 					},
		// 				},
		// 				"tags": schema.MapAttribute{
		// 					Optional:    true,
		// 					Description: "Assume role session tags.",
		// 					ElementType: types.StringType,
		// 				},
		// 				"transitive_tag_keys": schema.SetAttribute{
		// 					Optional:    true,
		// 					Description: "Assume role session tag keys to pass to any subsequent sessions.",
		// 					ElementType: types.StringType,
		// 				},
		// 			},
		// 		},
		// 	},
		// },
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

	mongodbgovCloud := data.IsMongodbGovCloud.ValueBool()
	if mongodbgovCloud {
		baseURL = "https://cloud.mongodbgov.com"
	} else {
		baseURL = data.BaseURL.ValueString()
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

	// if v, ok := d.GetOk("assume_role"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
	// config.AssumeRole = expandAssumeRole(v.([]interface{})[0].(map[string]interface{}))
	// secret := data.SecretName.ValueString()
	// region := data.Region.ValueString()
	// awsAccessKeyID := data.AwsAccessKeyID.ValueString()
	// awsSecretAccessKey := data.AwsSecretAccessKeyID.ValueString()
	// awsSessionToken := data.AwsSessionToken.ValueString()
	// endpoint := data.StsEndpoint.ValueString()
	// var err error
	// config, err = ConfigureCredentialsSTS(&config, secret, region, awsAccessKeyID, awsSecretAccessKey, awsSessionToken, endpoint)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"failed to initialize a new client",
	// 		err.Error(),
	// 	)
	// 	return
	// }

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
		func() resource.Resource {
			return &ExampleResource{}
		},
	}
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &MongodbtlasProvider{}
	}
}
