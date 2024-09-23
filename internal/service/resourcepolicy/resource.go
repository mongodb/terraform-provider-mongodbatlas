package resourcepolicy

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
)

var _ resource.ResourceWithConfigure = &resourcePolicyRS{}
var _ resource.ResourceWithImportState = &resourcePolicyRS{}

const (
	resourceName     = "resource_policy"
	fullResourceName = "mongodbatlas_" + resourceName
	errorCreate      = "error creating resource " + fullResourceName
	errorRead        = "error reading resource " + fullResourceName
	errorUpdate      = "error updating resource " + fullResourceName
)

func Resource() resource.Resource {
	return &resourcePolicyRS{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type resourcePolicyRS struct {
	config.RSCommon
}

func (r *resourcePolicyRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *resourcePolicyRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgID := plan.OrgID.ValueString()
	policies := NewAdminPolicies(ctx, plan.Policies)

	connV2 := r.Client.AtlasV2
	policySDK, _, err := connV2.AtlasResourcePoliciesApi.CreateAtlasResourcePolicy(ctx, orgID, &admin.ApiAtlasResourcePolicyCreate{
		Name:     plan.Name.ValueString(),
		Policies: policies,
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorCreate, err.Error())
		return
	}
	newResourcePolicyModel, diags := NewTFModel(ctx, policySDK)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicyModel)...)
}

func (r *resourcePolicyRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgID := state.OrgID.ValueString()
	resourcePolicyID := state.ID.ValueString()
	connV2 := r.Client.AtlasV2
	policySDK, apiResp, err := connV2.AtlasResourcePoliciesApi.GetAtlasResourcePolicy(ctx, orgID, resourcePolicyID).Execute()

	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(errorRead, err.Error())
		return
	}

	newResourcePolicyModel, diags := NewTFModel(ctx, policySDK)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicyModel)...)
}

func (r *resourcePolicyRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgID := plan.OrgID.ValueString()
	resourcePolicyID := plan.ID.ValueString()
	connV2 := r.Client.AtlasV2
	resourcePolicyAPI := connV2.AtlasResourcePoliciesApi
	editAdmin := admin.ApiAtlasResourcePolicyEdit{}
	if plan.Name.ValueString() != state.Name.ValueString() {
		editAdmin.SetName(plan.Name.ValueString())
	}
	policiesBefore := NewAdminPolicies(ctx, state.Policies)
	policiesAfter := NewAdminPolicies(ctx, plan.Policies)
	// comparing SDK models to check only the policy.Body for changes to avoid nested policy.Id updates
	if !reflect.DeepEqual(policiesBefore, policiesAfter) {
		editAdmin.SetPolicies(policiesAfter)
	}
	policySDK, _, err := resourcePolicyAPI.UpdateAtlasResourcePolicy(ctx, orgID, resourcePolicyID, &editAdmin).Execute()

	if err != nil {
		resp.Diagnostics.AddError(errorUpdate, err.Error())
		return
	}
	newResourcePolicyModel, diags := NewTFModel(ctx, policySDK)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicyModel)...)
}

func (r *resourcePolicyRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var resourcePolicyState *TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &resourcePolicyState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgID := resourcePolicyState.OrgID.ValueString()
	resourcePolicyID := resourcePolicyState.ID.ValueString()
	connV2 := r.Client.AtlasV2
	resourcePolicyAPI := connV2.AtlasResourcePoliciesApi
	if _, _, err := resourcePolicyAPI.DeleteAtlasResourcePolicy(ctx, orgID, resourcePolicyID).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *resourcePolicyRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	orgID, resourcePolicyID, err := splitImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting search deployment import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourcePolicyID)...)
}

func splitImportID(id string) (orgID, resourcePolicyID string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("use the format {org_id}-{resource_policy_id}")
		return
	}

	orgID = parts[1]
	resourcePolicyID = parts[2]
	return
}
