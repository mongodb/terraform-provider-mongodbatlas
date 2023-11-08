package mongodbatlas

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20231001002/admin"
)

var _ resource.ResourceWithConfigure = &ProjectRS{}
var _ resource.ResourceWithImportState = &ProjectRS{}

func NewSearchNodeRS() resource.Resource {
	return &SearchNodeRS{
		RSCommon: RSCommon{
			resourceName: "search_node",
		},
	}
}

type SearchNodeRS struct {
	RSCommon
}

type tfSearchNodeRSModel struct {
	ID          types.String `tfsdk:"id"`
	ClusterName types.String `tfsdk:"cluster_name"`
	ProjectID   types.String `tfsdk:"project_id"`
	Specs       types.List   `tfsdk:"specs"`
	StateName   types.String `tfsdk:"state_name"`
}

type tfSearchNodeSpecModel struct {
	InstanceSize types.String `tfsdk:"instance_size"`
	NodeCount    types.Int64  `tfsdk:"node_count"`
}

var SpecObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"instance_size": types.StringType,
	"node_count":    types.Int64Type,
}}

func (r *SearchNodeRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"cluster_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"specs": schema.ListNestedAttribute{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_size": schema.StringAttribute{
							Required: true,
						},
						"node_count": schema.Int64Attribute{
							Required: true,
						},
					},
				},
				Required: true,
			},
			"state_name": schema.StringAttribute{ // TODO: use for sync creation?
				Computed: true,
			},
		},
	}
}

func (r *SearchNodeRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var searchNodePlan tfSearchNodeRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &searchNodePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.client.AtlasV2
	projectID := searchNodePlan.ProjectID.ValueString()
	clusterName := searchNodePlan.ClusterName.ValueString()
	searchDeploymentReq := newSearchDeploymentReq(ctx, &searchNodePlan)
	deploymentResp, _, err := connV2.AtlasSearchApi.CreateAtlasSearchDeployment(ctx, projectID, clusterName, &searchDeploymentReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error during search node creation", err.Error())
		return
	}

	newSearchNodeModel, diagnostics := newTFSearchDeployment(ctx, clusterName, deploymentResp)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *SearchNodeRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var searchNodePlan tfSearchNodeRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &searchNodePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.client.AtlasV2
	projectID := searchNodePlan.ProjectID.ValueString()
	clusterName := searchNodePlan.ClusterName.ValueString()
	deploymentResp, _, err := connV2.AtlasSearchApi.GetAtlasSearchDeployment(ctx, projectID, clusterName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting search node information", err.Error())
		return
	}

	newSearchNodeModel, diagnostics := newTFSearchDeployment(ctx, clusterName, deploymentResp)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *SearchNodeRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var searchNodePlan tfSearchNodeRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &searchNodePlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.client.AtlasV2
	projectID := searchNodePlan.ProjectID.ValueString()
	clusterName := searchNodePlan.ClusterName.ValueString()
	searchDeploymentReq := newSearchDeploymentReq(ctx, &searchNodePlan)
	deploymentResp, _, err := connV2.AtlasSearchApi.UpdateAtlasSearchDeployment(ctx, projectID, clusterName, &searchDeploymentReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error during search node update", err.Error())
		return
	}

	newSearchNodeModel, diagnostics := newTFSearchDeployment(ctx, clusterName, deploymentResp)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *SearchNodeRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var searchNodeState *tfSearchNodeRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &searchNodeState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.client.AtlasV2
	projectID := searchNodeState.ProjectID.ValueString()
	clusterName := searchNodeState.ClusterName.ValueString()
	_, err := connV2.AtlasSearchApi.DeleteAtlasSearchDeployment(ctx, projectID, clusterName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error during search node delete", err.Error())
		return
	}
}

func (r *SearchNodeRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, clusterName, err := splitSearchNodeImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting search node import ID", err.Error())
		return
	}

	searchNodeModel := tfSearchNodeRSModel{ // read operation requires projectID and clusterName to be defined
		ProjectID:   types.StringValue(projectID),
		ClusterName: types.StringValue(clusterName),
		Specs:       types.ListNull(SpecObjectType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &searchNodeModel)...)
}

func splitSearchNodeImportID(id string) (projectID, clusterName string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("use the format {project_id}-{cluster_name}")
		return
	}

	projectID = parts[1]
	clusterName = parts[2]
	return
}

func newSearchDeploymentReq(ctx context.Context, searchNodePlan *tfSearchNodeRSModel) admin.ApiSearchDeploymentRequest {
	var specs []tfSearchNodeSpecModel
	searchNodePlan.Specs.ElementsAs(ctx, &specs, true)

	resultSpecs := make([]admin.ApiSearchDeploymentSpec, len(specs))
	for i, spec := range specs {
		resultSpecs[i] = admin.ApiSearchDeploymentSpec{
			InstanceSize: spec.InstanceSize.ValueString(),
			NodeCount:    int(spec.NodeCount.ValueInt64()),
		}
	}

	return admin.ApiSearchDeploymentRequest{
		Specs: resultSpecs,
	}
}

func newTFSearchDeployment(ctx context.Context, clusterName string, deployResp *admin.ApiSearchDeploymentResponse) (*tfSearchNodeRSModel, diag.Diagnostics) {
	result := tfSearchNodeRSModel{
		ID:          types.StringPointerValue(deployResp.Id),
		ClusterName: types.StringValue(clusterName),
		ProjectID:   types.StringPointerValue(deployResp.GroupId),
		StateName:   types.StringPointerValue(deployResp.StateName),
	}

	specsList, diagnostics := types.ListValueFrom(ctx, SpecObjectType, newTFSpecsModel(deployResp.Specs))
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	result.Specs = specsList
	return &result, nil
}

func newTFSpecsModel(specs []admin.ApiSearchDeploymentSpec) []tfSearchNodeSpecModel {
	result := make([]tfSearchNodeSpecModel, len(specs))
	for i, v := range specs {
		result[i] = tfSearchNodeSpecModel{
			InstanceSize: types.StringValue(v.InstanceSize),
			NodeCount:    types.Int64Value(int64(v.NodeCount)),
		}
	}

	return result
}
