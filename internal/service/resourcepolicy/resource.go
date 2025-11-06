package resourcepolicy

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &resourcePolicyRS{}
var _ resource.ResourceWithImportState = &resourcePolicyRS{}
var _ resource.ResourceWithModifyPlan = &resourcePolicyRS{}

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

func (r *resourcePolicyRS) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var policies []TFPolicyModel
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("policies"), &policies)...)
	sdkPolicies := NewAdminPolicies(ctx, policies)
	var orgID, name *string
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("org_id"), &orgID)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if name == nil || orgID == nil {
		return
	}
	sdkCreate := &admin.ApiAtlasResourcePolicyCreate{
		Name:     *name,
		Policies: sdkPolicies,
	}
	connV2 := r.Client.AtlasV2
	_, _, err := connV2.ResourcePoliciesApi.ValidateResourcePolicies(ctx, *orgID, sdkCreate).Execute()
	if err != nil {
		conversion.AddJSONBodyErrorToDiagnostics(fmt.Sprintf("Policy Validation failed (name=%s): ", *name), err, &resp.Diagnostics)
	}
}

func (r *resourcePolicyRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
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
	policySDK, _, err := connV2.ResourcePoliciesApi.CreateOrgResourcePolicy(ctx, orgID, &admin.ApiAtlasResourcePolicyCreate{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		Policies:    policies,
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
	policySDK, apiResp, err := connV2.ResourcePoliciesApi.GetOrgResourcePolicy(ctx, orgID, resourcePolicyID).Execute()

	if err != nil {
		if validate.StatusNotFound(apiResp) {
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
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgID := plan.OrgID.ValueString()
	resourcePolicyID := plan.ID.ValueString()
	connV2 := r.Client.AtlasV2
	policies := NewAdminPolicies(ctx, plan.Policies)
	editAdmin := admin.ApiAtlasResourcePolicyEdit{
		Name: plan.Name.ValueStringPointer(),
		// description is an optional attribute (i.e. null by default), if it is removed from the config during an update
		// (i.e. user wants to remove the existing description from resource policy), we send an empty string ("") as the value in API request for update
		Description: conversion.Pointer(plan.Description.ValueString()),
		Policies:    &policies,
	}
	policySDK, _, err := connV2.ResourcePoliciesApi.UpdateOrgResourcePolicy(ctx, orgID, resourcePolicyID, &editAdmin).Execute()

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
	resourcePolicyAPI := connV2.ResourcePoliciesApi
	if _, err := resourcePolicyAPI.DeleteOrgResourcePolicy(ctx, orgID, resourcePolicyID).Execute(); err != nil {
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
