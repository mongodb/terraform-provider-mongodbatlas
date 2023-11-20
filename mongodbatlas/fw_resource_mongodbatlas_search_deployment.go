package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework"
	retrystrategy "github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/util"
	"go.mongodb.org/atlas-sdk/v20231115001/admin"
)

var _ resource.ResourceWithConfigure = &SearchDeploymentRS{}
var _ resource.ResourceWithImportState = &SearchDeploymentRS{}

const (
	searchDeploymentDoesNotExistsError = "ATLAS_FTS_DEPLOYMENT_DOES_NOT_EXIST"
	searchDeploymentName               = "search_deployment"
)

func NewSearchDeploymentRS() resource.Resource {
	return &SearchDeploymentRS{
		RSCommon: framework.RSCommon{
			ResourceName: searchDeploymentName,
		},
	}
}

type SearchDeploymentRS struct {
	framework.RSCommon
}

type tfSearchDeploymentRSModel struct {
	ID          types.String   `tfsdk:"id"`
	ClusterName types.String   `tfsdk:"cluster_name"`
	ProjectID   types.String   `tfsdk:"project_id"`
	Specs       types.List     `tfsdk:"specs"`
	StateName   types.String   `tfsdk:"state_name"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

type tfSearchNodeSpecModel struct {
	InstanceSize types.String `tfsdk:"instance_size"`
	NodeCount    types.Int64  `tfsdk:"node_count"`
}

var SpecObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"instance_size": types.StringType,
	"node_count":    types.Int64Type,
}}

func (r *SearchDeploymentRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"state_name": schema.StringAttribute{
				Computed: true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

const defaultSearchNodeTimeout time.Duration = 3 * time.Hour

func (r *SearchDeploymentRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var searchDeploymentPlan tfSearchDeploymentRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &searchDeploymentPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := searchDeploymentPlan.ProjectID.ValueString()
	clusterName := searchDeploymentPlan.ClusterName.ValueString()
	searchDeploymentReq := newSearchDeploymentReq(ctx, &searchDeploymentPlan)
	if _, _, err := connV2.AtlasSearchApi.CreateAtlasSearchDeployment(ctx, projectID, clusterName, &searchDeploymentReq).Execute(); err != nil {
		resp.Diagnostics.AddError("error during search deployment creation", err.Error())
		return
	}

	createTimeout, diags := searchDeploymentPlan.Timeouts.Create(ctx, defaultSearchNodeTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	deploymentResp, err := waitSearchNodeStateTransition(ctx, projectID, clusterName, connV2, createTimeout)
	if err != nil {
		resp.Diagnostics.AddError("error during search deployment creation", err.Error())
		return
	}
	newSearchNodeModel, diagnostics := newTFSearchDeployment(ctx, clusterName, deploymentResp, &searchDeploymentPlan.Timeouts)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *SearchDeploymentRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var searchDeploymentPlan tfSearchDeploymentRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &searchDeploymentPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := searchDeploymentPlan.ProjectID.ValueString()
	clusterName := searchDeploymentPlan.ClusterName.ValueString()
	deploymentResp, _, err := connV2.AtlasSearchApi.GetAtlasSearchDeployment(ctx, projectID, clusterName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting search deployment information", err.Error())
		return
	}

	newSearchNodeModel, diagnostics := newTFSearchDeployment(ctx, clusterName, deploymentResp, &searchDeploymentPlan.Timeouts)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *SearchDeploymentRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var searchDeploymentPlan tfSearchDeploymentRSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &searchDeploymentPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := searchDeploymentPlan.ProjectID.ValueString()
	clusterName := searchDeploymentPlan.ClusterName.ValueString()
	searchDeploymentReq := newSearchDeploymentReq(ctx, &searchDeploymentPlan)
	if _, _, err := connV2.AtlasSearchApi.UpdateAtlasSearchDeployment(ctx, projectID, clusterName, &searchDeploymentReq).Execute(); err != nil {
		resp.Diagnostics.AddError("error during search deployment update", err.Error())
		return
	}

	updateTimeout, diags := searchDeploymentPlan.Timeouts.Update(ctx, defaultSearchNodeTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	deploymentResp, err := waitSearchNodeStateTransition(ctx, projectID, clusterName, connV2, updateTimeout)
	if err != nil {
		resp.Diagnostics.AddError("error during search deployment update", err.Error())
		return
	}
	newSearchNodeModel, diagnostics := newTFSearchDeployment(ctx, clusterName, deploymentResp, &searchDeploymentPlan.Timeouts)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newSearchNodeModel)...)
}

func (r *SearchDeploymentRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var searchDeploymentState *tfSearchDeploymentRSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &searchDeploymentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := searchDeploymentState.ProjectID.ValueString()
	clusterName := searchDeploymentState.ClusterName.ValueString()
	if _, err := connV2.AtlasSearchApi.DeleteAtlasSearchDeployment(ctx, projectID, clusterName).Execute(); err != nil {
		resp.Diagnostics.AddError("error during search deployment delete", err.Error())
		return
	}

	deleteTimeout, diags := searchDeploymentState.Timeouts.Delete(ctx, defaultSearchNodeTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := waitSearchNodeDelete(ctx, projectID, clusterName, connV2, deleteTimeout); err != nil {
		resp.Diagnostics.AddError("error during search deployment delete", err.Error())
		return
	}
}

func (r *SearchDeploymentRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, clusterName, err := splitSearchNodeImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting search deployment import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_name"), clusterName)...)
	if resp.Diagnostics.HasError() {
		return
	}
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

func waitSearchNodeStateTransition(ctx context.Context, projectID, clusterName string, conn *admin.APIClient, timeout time.Duration) (*admin.ApiSearchDeploymentResponse, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyUpdatingState, retrystrategy.RetryStrategyPausedState},
		Target:     []string{retrystrategy.RetryStrategyIdleState},
		Refresh:    searchDeploymentRefreshFunc(ctx, projectID, clusterName, conn),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      1 * time.Minute,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if deploymentResp, ok := result.(*admin.ApiSearchDeploymentResponse); ok && deploymentResp != nil {
		return deploymentResp, nil
	}
	return nil, errors.New("did not obtain valid result when waiting for search deployment state transition")
}

func waitSearchNodeDelete(ctx context.Context, projectID, clusterName string, conn *admin.APIClient, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyIdleState, retrystrategy.RetryStrategyUpdatingState, retrystrategy.RetryStrategyPausedState},
		Target:     []string{retrystrategy.RetryStrategyDeletedState},
		Refresh:    searchDeploymentRefreshFunc(ctx, projectID, clusterName, conn),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func searchDeploymentRefreshFunc(ctx context.Context, projectID, clusterName string, conn *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		deploymentResp, resp, err := conn.AtlasSearchApi.GetAtlasSearchDeployment(ctx, projectID, clusterName).Execute()
		if err != nil && deploymentResp == nil && resp == nil {
			return nil, "", err
		}
		if err != nil {
			if resp.StatusCode == 400 && strings.Contains(err.Error(), searchDeploymentDoesNotExistsError) {
				return "", retrystrategy.RetryStrategyDeletedState, nil
			}
			if resp.StatusCode == 503 {
				return "", retrystrategy.RetryStrategyPendingState, nil
			}
			return nil, "", err
		}

		if util.IsStringPresent(deploymentResp.StateName) {
			tflog.Debug(ctx, fmt.Sprintf("search deployment status: %s", *deploymentResp.StateName))
			return deploymentResp, *deploymentResp.StateName, nil
		}
		return deploymentResp, "", nil
	}
}

func newSearchDeploymentReq(ctx context.Context, searchDeploymentPlan *tfSearchDeploymentRSModel) admin.ApiSearchDeploymentRequest {
	var specs []tfSearchNodeSpecModel
	searchDeploymentPlan.Specs.ElementsAs(ctx, &specs, true)

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

func newTFSearchDeployment(ctx context.Context, clusterName string, deployResp *admin.ApiSearchDeploymentResponse, timeout *timeouts.Value) (*tfSearchDeploymentRSModel, diag.Diagnostics) {
	result := tfSearchDeploymentRSModel{
		ID:          types.StringPointerValue(deployResp.Id),
		ClusterName: types.StringValue(clusterName),
		ProjectID:   types.StringPointerValue(deployResp.GroupId),
		StateName:   types.StringPointerValue(deployResp.StateName),
	}

	if timeout != nil {
		result.Timeouts = *timeout
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
