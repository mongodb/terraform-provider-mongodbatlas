package streamconnection

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312021/admin"
)

const streamConnectionFailoverName = "stream_connection_failover"

var _ resource.ResourceWithConfigure = &streamConnectionFailoverRS{}
var _ resource.ResourceWithImportState = &streamConnectionFailoverRS{}

func FailoverResource() resource.Resource {
	return &streamConnectionFailoverRS{
		RSCommon: config.RSCommon{
			ResourceName: streamConnectionFailoverName,
		},
	}
}

type streamConnectionFailoverRS struct {
	config.RSCommon
}

// TFStreamConnectionFailoverModel is the resource model for a single failover (regional-alternate)
// stream connection. A failover connection shares its primary connection's name and is uniquely
// identified by its region and failover connection id.
type TFStreamConnectionFailoverModel struct {
	ID                           types.String `tfsdk:"id"`
	ProjectID                    types.String `tfsdk:"project_id"`
	InstanceName                 types.String `tfsdk:"instance_name"`
	WorkspaceName                types.String `tfsdk:"workspace_name"`
	ConnectionName               types.String `tfsdk:"connection_name"`
	Region                       types.String `tfsdk:"region"`
	Type                         types.String `tfsdk:"type"`
	ClusterName                  types.String `tfsdk:"cluster_name"`
	ClusterProjectID             types.String `tfsdk:"cluster_project_id"`
	DBRoleToExecute              types.Object `tfsdk:"db_role_to_execute"`
	BootstrapServers             types.String `tfsdk:"bootstrap_servers"`
	Authentication               types.Object `tfsdk:"authentication"`
	Config                       types.Map    `tfsdk:"config"`
	Security                     types.Object `tfsdk:"security"`
	Networking                   types.Object `tfsdk:"networking"`
	AWS                          types.Object `tfsdk:"aws"`
	Azure                        types.Object `tfsdk:"azure"`
	GCP                          types.Object `tfsdk:"gcp"`
	URL                          types.String `tfsdk:"url"`
	Headers                      types.Map    `tfsdk:"headers"`
	SchemaRegistryProvider       types.String `tfsdk:"schema_registry_provider"`
	SchemaRegistryURLs           types.List   `tfsdk:"schema_registry_urls"`
	SchemaRegistryAuthentication types.Object `tfsdk:"schema_registry_authentication"`
}

func (m *TFStreamConnectionFailoverModel) workspaceOrInstanceName() string {
	if !m.WorkspaceName.IsNull() && !m.WorkspaceName.IsUnknown() {
		return m.WorkspaceName.ValueString()
	}
	if !m.InstanceName.IsNull() && !m.InstanceName.IsUnknown() {
		return m.InstanceName.ValueString()
	}
	return ""
}

// toConnModel maps the failover model onto a TFStreamConnectionModel so the shared connection
// conversion helpers can be reused. The failover connection's body name equals the primary
// connection name.
func (m *TFStreamConnectionFailoverModel) toConnModel() *TFStreamConnectionModel {
	return &TFStreamConnectionModel{
		TFStreamConnectionCommonModel: TFStreamConnectionCommonModel{
			ConnectionName:               m.ConnectionName,
			Type:                         m.Type,
			ClusterName:                  m.ClusterName,
			ClusterProjectID:             m.ClusterProjectID,
			DBRoleToExecute:              m.DBRoleToExecute,
			BootstrapServers:             m.BootstrapServers,
			Authentication:               m.Authentication,
			Config:                       m.Config,
			Security:                     m.Security,
			Networking:                   m.Networking,
			AWS:                          m.AWS,
			Azure:                        m.Azure,
			GCP:                          m.GCP,
			URL:                          m.URL,
			Headers:                      m.Headers,
			SchemaRegistryProvider:       m.SchemaRegistryProvider,
			SchemaRegistryURLs:           m.SchemaRegistryURLs,
			SchemaRegistryAuthentication: m.SchemaRegistryAuthentication,
		},
	}
}

func (r *streamConnectionFailoverRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("workspace_name")),
				},
				DeprecationMessage: fmt.Sprintf(constant.DeprecationParamWithReplacement, "workspace_name"),
			},
			"workspace_name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("instance_name")),
				},
			},
			"connection_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
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
			"cluster_name": schema.StringAttribute{
				Optional: true,
			},
			"cluster_project_id": schema.StringAttribute{
				Optional: true,
			},
			"db_role_to_execute": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"role": schema.StringAttribute{Required: true},
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("BUILT_IN", "CUSTOM"),
						},
					},
				},
			},
			"bootstrap_servers": schema.StringAttribute{
				Optional: true,
			},
			"authentication": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"mechanism":                   schema.StringAttribute{Optional: true},
					"method":                      schema.StringAttribute{Optional: true},
					"password":                    schema.StringAttribute{Optional: true, Sensitive: true},
					"username":                    schema.StringAttribute{Optional: true},
					"token_endpoint_url":          schema.StringAttribute{Optional: true},
					"client_id":                   schema.StringAttribute{Optional: true},
					"client_secret":               schema.StringAttribute{Optional: true, Sensitive: true},
					"scope":                       schema.StringAttribute{Optional: true},
					"sasl_oauthbearer_extensions": schema.StringAttribute{Optional: true},
				},
			},
			"config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"security": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"broker_public_certificate": schema.StringAttribute{Optional: true},
					"protocol":                  schema.StringAttribute{Optional: true},
				},
			},
			"networking": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"access": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"type":          schema.StringAttribute{Required: true},
							"connection_id": schema.StringAttribute{Optional: true},
						},
					},
				},
			},
			"aws": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"role_arn": schema.StringAttribute{Required: true},
				},
			},
			"azure": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"service_principal_id": schema.StringAttribute{Required: true},
					"storage_account_name": schema.StringAttribute{Required: true},
					"region":               schema.StringAttribute{Optional: true},
				},
			},
			"gcp": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"service_account_id": schema.StringAttribute{Required: true},
				},
			},
			"url": schema.StringAttribute{
				Optional: true,
			},
			"headers": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"schema_registry_provider": schema.StringAttribute{
				Optional: true,
			},
			"schema_registry_urls": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"schema_registry_authentication": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"type":     schema.StringAttribute{Optional: true},
					"username": schema.StringAttribute{Optional: true},
					"password": schema.StringAttribute{Optional: true, Sensitive: true},
				},
			},
		},
	}
}

func (r *streamConnectionFailoverRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFStreamConnectionFailoverModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceName := plan.workspaceOrInstanceName()
	if workspaceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}

	connReq, diags := newStreamConnectionFailoverReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	connV2 := r.Client.AtlasV2
	created, _, err := connV2.StreamsApi.CreateFailoverConnection(ctx, plan.ProjectID.ValueString(), workspaceName, plan.ConnectionName.ValueString(), connReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating failover connection", err.Error())
		return
	}

	// Build state from the plan (source of truth for sensitive/config fields) and only inject the
	// computed failover connection id from the create response, which can be sparse.
	plan.ID = types.StringPointerValue(created.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *streamConnectionFailoverRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFStreamConnectionFailoverModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceName := state.workspaceOrInstanceName()
	if workspaceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}

	connV2 := r.Client.AtlasV2
	apiResp, getResp, err := connV2.StreamsApi.GetStreamFailoverConnection(ctx, state.ProjectID.ValueString(), workspaceName, state.ConnectionName.ValueString(), state.ID.ValueString()).Execute()
	if err != nil {
		if validate.StatusNotFound(getResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error reading failover connection", err.Error())
		return
	}

	newModel, diags := newTFStreamConnectionFailover(ctx, &state, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newModel)...)
}

func (r *streamConnectionFailoverRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFStreamConnectionFailoverModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state TFStreamConnectionFailoverModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceName := plan.workspaceOrInstanceName()
	if workspaceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}

	connReq, diags := newStreamConnectionFailoverReq(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	connV2 := r.Client.AtlasV2
	if _, _, err := connV2.StreamsApi.UpdateStreamFailoverConnection(ctx, plan.ProjectID.ValueString(), workspaceName, plan.ConnectionName.ValueString(), state.ID.ValueString(), connReq).Execute(); err != nil {
		resp.Diagnostics.AddError("error updating failover connection", err.Error())
		return
	}

	// Preserve the existing id and build state from the plan (see Create).
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *streamConnectionFailoverRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFStreamConnectionFailoverModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceName := state.workspaceOrInstanceName()
	if workspaceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}

	connV2 := r.Client.AtlasV2
	if _, err := connV2.StreamsApi.DeleteStreamFailoverConnection(ctx, state.ProjectID.ValueString(), workspaceName, state.ConnectionName.ValueString(), state.ID.ValueString()).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting failover connection", err.Error())
		return
	}
}

// ImportState expects an id of the form "{workspaceName}-{projectID}-{connectionName}-{failoverConnectionId}".
// projectID and failoverConnectionId are 24-hex, which anchors the split so that workspace and
// connection names may themselves contain dashes.
func (r *streamConnectionFailoverRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	workspaceName, projectID, connectionName, failoverConnectionID, err := splitStreamConnectionFailoverImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting failover connection import ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_name"), workspaceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_name"), workspaceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connection_name"), connectionName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), failoverConnectionID)...)
}

func splitStreamConnectionFailoverImportID(id string) (workspaceName, projectID, connectionName, failoverConnectionID string, err error) {
	re := regexp.MustCompile(`^(.*)-([0-9a-fA-F]{24})-(.*)-([0-9a-fA-F]{24})$`)
	parts := re.FindStringSubmatch(id)
	if len(parts) != 5 {
		err = errors.New("use the format {workspace_name}-{project_id}-{connection_name}-{failover_connection_id}")
		return
	}
	workspaceName = parts[1]
	projectID = parts[2]
	connectionName = parts[3]
	failoverConnectionID = parts[4]
	return
}

// newStreamConnectionFailoverReq builds the StreamsConnection request body for a failover connection,
// reusing the primary connection request builder and adding the failover region.
func newStreamConnectionFailoverReq(ctx context.Context, plan *TFStreamConnectionFailoverModel) (*admin.StreamsConnection, diag.Diagnostics) {
	conn, diags := NewStreamConnectionReq(ctx, plan.toConnModel())
	if diags.HasError() {
		return nil, diags
	}
	conn.Region = plan.Region.ValueStringPointer()
	return conn, diags
}

// newTFStreamConnectionFailover maps a StreamsConnection API response onto the failover resource model,
// reusing the primary connection conversion and preserving parent references + sensitive fields from prior.
func newTFStreamConnectionFailover(ctx context.Context, prior *TFStreamConnectionFailoverModel, apiResp *admin.StreamsConnection) (*TFStreamConnectionFailoverModel, diag.Diagnostics) {
	model, diags := NewTFStreamConnection(ctx, "", "", "", &prior.Authentication, &prior.SchemaRegistryAuthentication, apiResp, nil)
	if diags.HasError() {
		return nil, diags
	}
	// Resolve to a single workspace field: import sets both instance_name and workspace_name, and
	// only one should end up in state (workspace_name is preferred; instance_name is deprecated).
	instanceName, workspaceName := prior.InstanceName, prior.WorkspaceName
	if !workspaceName.IsNull() && !workspaceName.IsUnknown() && workspaceName.ValueString() != "" {
		instanceName = types.StringNull()
	} else {
		workspaceName = types.StringNull()
	}

	c := model.TFStreamConnectionCommonModel
	return &TFStreamConnectionFailoverModel{
		ID:                           types.StringPointerValue(apiResp.Id),
		ProjectID:                    prior.ProjectID,
		InstanceName:                 instanceName,
		WorkspaceName:                workspaceName,
		ConnectionName:               prior.ConnectionName,
		Region:                       types.StringPointerValue(apiResp.Region),
		Type:                         c.Type,
		ClusterName:                  c.ClusterName,
		ClusterProjectID:             c.ClusterProjectID,
		DBRoleToExecute:              c.DBRoleToExecute,
		BootstrapServers:             c.BootstrapServers,
		Authentication:               c.Authentication,
		Config:                       c.Config,
		Security:                     c.Security,
		Networking:                   c.Networking,
		AWS:                          c.AWS,
		Azure:                        c.Azure,
		GCP:                          c.GCP,
		URL:                          c.URL,
		Headers:                      c.Headers,
		SchemaRegistryProvider:       c.SchemaRegistryProvider,
		SchemaRegistryURLs:           c.SchemaRegistryURLs,
		SchemaRegistryAuthentication: c.SchemaRegistryAuthentication,
	}, nil
}
