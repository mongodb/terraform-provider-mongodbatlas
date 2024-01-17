package streamconnection

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const streamConnectionName = "stream_connection"

var _ resource.ResourceWithConfigure = &streamConnectionRS{}
var _ resource.ResourceWithImportState = &streamConnectionRS{}

func Resource() resource.Resource {
	return &streamConnectionRS{
		RSCommon: config.RSCommon{
			ResourceName: streamConnectionName,
		},
	}
}

type streamConnectionRS struct {
	config.RSCommon
}

type TFStreamConnectionModel struct {
	ID               types.String `tfsdk:"id"`
	ProjectID        types.String `tfsdk:"project_id"`
	InstanceName     types.String `tfsdk:"instance_name"`
	ConnectionName   types.String `tfsdk:"connection_name"`
	Type             types.String `tfsdk:"type"`
	ClusterName      types.String `tfsdk:"cluster_name"`
	Authentication   types.Object `tfsdk:"authentication"`
	BootstrapServers types.String `tfsdk:"bootstrap_servers"`
	Config           types.Map    `tfsdk:"config"`
	Security         types.Object `tfsdk:"security"`
}

type TFConnectionAuthenticationModel struct {
	Mechanism types.String `tfsdk:"mechanism"`
	Password  types.String `tfsdk:"password"`
	Username  types.String `tfsdk:"username"`
}

var ConnectionAuthenticationObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"mechanism": types.StringType,
	"password":  types.StringType,
	"username":  types.StringType,
}}

type TFConnectionSecurityModel struct {
	BrokerPublicCertificate types.String `tfsdk:"broker_public_certificate"`
	Protocol                types.String `tfsdk:"protocol"`
}

var ConnectionSecurityObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"broker_public_certificate": types.StringType,
	"protocol":                  types.StringType,
}}

func (r *streamConnectionRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connection_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// cluster type specific
			"cluster_name": schema.StringAttribute{
				Optional: true,
			},

			// kafka type specific
			"authentication": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"mechanism": schema.StringAttribute{
						Optional: true,
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"username": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			"bootstrap_servers": schema.StringAttribute{
				Optional: true,
			},
			"config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"security": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"broker_public_certificate": schema.StringAttribute{
						Optional: true,
					},
					"protocol": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *streamConnectionRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var streamConnectionPlan TFStreamConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamConnectionPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamConnectionPlan.ProjectID.ValueString()
	instanceName := streamConnectionPlan.InstanceName.ValueString()
	streamConnectionReq, diags := NewStreamConnectionReq(ctx, &streamConnectionPlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.CreateStreamConnection(ctx, projectID, instanceName, streamConnectionReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	newStreamConnectionModel, diags := NewTFStreamConnection(ctx, projectID, instanceName, &streamConnectionPlan.Authentication, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}

func (r *streamConnectionRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var streamConnectionState TFStreamConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamConnectionState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamConnectionState.ProjectID.ValueString()
	instanceName := streamConnectionState.InstanceName.ValueString()
	connectionName := streamConnectionState.ConnectionName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamConnection(ctx, projectID, instanceName, connectionName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newStreamConnectionModel, diags := NewTFStreamConnection(ctx, projectID, instanceName, &streamConnectionState.Authentication, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}

func (r *streamConnectionRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var streamConnectionPlan TFStreamConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamConnectionPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamConnectionPlan.ProjectID.ValueString()
	instanceName := streamConnectionPlan.InstanceName.ValueString()
	connectionName := streamConnectionPlan.ConnectionName.ValueString()
	streamConnectionReq, diags := NewStreamConnectionReq(ctx, &streamConnectionPlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.UpdateStreamConnection(ctx, projectID, instanceName, connectionName, streamConnectionReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	newStreamConnectionModel, diags := NewTFStreamConnection(ctx, projectID, instanceName, &streamConnectionPlan.Authentication, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}

func (r *streamConnectionRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var streamConnectionState *TFStreamConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamConnectionState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamConnectionState.ProjectID.ValueString()
	instanceName := streamConnectionState.InstanceName.ValueString()
	connectionName := streamConnectionState.ConnectionName.ValueString()
	if _, _, err := connV2.StreamsApi.DeleteStreamConnection(ctx, projectID, instanceName, connectionName).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *streamConnectionRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	instanceName, projectID, connectionName, err := splitStreamConnectionImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting stream connection import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_name"), instanceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connection_name"), connectionName)...)
}

func splitStreamConnectionImportID(id string) (instanceName, projectID, connectionName string, err error) {
	var re = regexp.MustCompile(`^(.*)-([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("use the format {instance_name}-{project_id}-{connection_name}")
		return
	}

	instanceName = parts[1]
	projectID = parts[2]
	connectionName = parts[3]
	return
}
