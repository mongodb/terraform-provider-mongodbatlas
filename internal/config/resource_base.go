package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	errorConfigureSummary = "Unexpected Resource Configure Type"
	errorConfigure        = "expected *MongoDBClient, got: %T. Please report this issue to the provider developers"
)

type ProviderMeta struct {
	ModuleName     types.String `tfsdk:"module_name"`
	ModuleVersion  types.String `tfsdk:"module_version"`
	UserAgentExtra types.Map    `tfsdk:"user_agent_extra"`
}

type ImplementedResource interface {
	resource.ResourceWithImportState
	// Additional methods such as upgrade state & plan modifier are optional
	SetClient(*MongoDBClient)
	GetName() string
}

func AnalyticsResourceFunc(iResource resource.Resource) func() resource.Resource {
	commonResource, ok := iResource.(ImplementedResource)
	if !ok {
		panic(fmt.Sprintf("resource %T didn't comply with the ImplementedResource interface", iResource))
	}
	return func() resource.Resource {
		return analyticsResource(commonResource)
	}
}

// analyticsResource wraps an ImplementedResource with RSCommon to add analytics tracking.
// We cannot return iResource directly because we need to intercept all CRUD operations
// to inject provider_meta information into the context before calling the actual resource methods.
func analyticsResource(iResource ImplementedResource) resource.Resource {
	return &RSCommon{
		ResourceName:        iResource.GetName(),
		ImplementedResource: iResource,
	}
}

// RSCommon is used as an embedded struct for all framework resources. Implements the following plugin-framework defined functions:
// - Metadata
// - Configure
// Client is left empty and populated by the framework when envoking Configure method.
// ResourceName must be defined when creating an instance of a resource.
type RSCommon struct {
	ImplementedResource
	Client       *MongoDBClient
	ResourceName string
}

func (r *RSCommon) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.ResourceName)
}

func (r *RSCommon) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := configureClient(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	r.ImplementedResource.SetClient(client)
}

func (r *RSCommon) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	extra := asUserAgentExtraFromProviderMeta(ctx, r.ResourceName, UserAgentOperationValueCreate, false, req.ProviderMeta)
	ctx = AddUserAgentExtra(ctx, extra)
	r.ImplementedResource.Create(ctx, req, resp)
}

func (r *RSCommon) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	extra := asUserAgentExtraFromProviderMeta(ctx, r.ResourceName, UserAgentOperationValueRead, false, req.ProviderMeta)
	ctx = AddUserAgentExtra(ctx, extra)
	r.ImplementedResource.Read(ctx, req, resp)
}

func (r *RSCommon) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	extra := asUserAgentExtraFromProviderMeta(ctx, r.ResourceName, UserAgentOperationValueUpdate, false, req.ProviderMeta)
	ctx = AddUserAgentExtra(ctx, extra)
	r.ImplementedResource.Update(ctx, req, resp)
}

func (r *RSCommon) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	extra := asUserAgentExtraFromProviderMeta(ctx, r.ResourceName, UserAgentOperationValueDelete, false, req.ProviderMeta)
	ctx = AddUserAgentExtra(ctx, extra)
	r.ImplementedResource.Delete(ctx, req, resp)
}

// Optional interfaces for resource.Resource
func (r *RSCommon) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// req resource.ImportStateRequest doesn't have ProviderMeta
	ctx = AddUserAgentExtra(ctx, UserAgentExtra{
		Name:      r.ResourceName,
		Operation: UserAgentOperationValueImport,
	})
	r.ImplementedResource.ImportState(ctx, req, resp)
}

func (r *RSCommon) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	resourceWithModifier, ok := r.ImplementedResource.(resource.ResourceWithModifyPlan)
	if !ok {
		return
	}
	extra := asUserAgentExtraFromProviderMeta(ctx, r.ResourceName, UserAgentOperationValuePlanModify, false, req.ProviderMeta)
	ctx = AddUserAgentExtra(ctx, extra)
	resourceWithModifier.ModifyPlan(ctx, req, resp)
}

func (r *RSCommon) MoveState(ctx context.Context) []resource.StateMover {
	resourceWithMoveState, ok := r.ImplementedResource.(resource.ResourceWithMoveState)
	if !ok {
		return nil
	}
	ctx = AddUserAgentExtra(ctx, UserAgentExtra{
		Name:      r.ResourceName,
		Operation: UserAgentOperationValueMoveState,
	})
	return resourceWithMoveState.MoveState(ctx)
}

func (r *RSCommon) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	resourceWithUpgradeState, ok := r.ImplementedResource.(resource.ResourceWithUpgradeState)
	if !ok {
		return nil
	}
	ctx = AddUserAgentExtra(ctx, UserAgentExtra{
		Name:      r.ResourceName,
		Operation: UserAgentOperationValueUpgradeState,
	})
	return resourceWithUpgradeState.UpgradeState(ctx)
}

// Extra methods not found on resource.Resource
func (r *RSCommon) GetName() string {
	return r.ResourceName
}

func (r *RSCommon) SetClient(client *MongoDBClient) {
	r.Client = client
}
