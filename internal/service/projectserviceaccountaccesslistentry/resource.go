package projectserviceaccountaccesslistentry

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	serviceaccountaccesslistentry "github.com/mongodb/terraform-provider-mongodbatlas/internal/service/serviceaccountaccesslistentry"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
)

const (
	resourceName = "project_service_account_access_list_entry"
)

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}

func Resource() resource.Resource {
	return &rs{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type rs struct {
	config.RSCommon
}

func (r *rs) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema()
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFProjectServiceAccountAccessListEntryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.CIDRBlock.ValueString() == "" && plan.IPAddress.ValueString() == "" {
		resp.Diagnostics.AddError("validation error", "cidr_block or ip_address must be provided")
		return
	}

	projectID := plan.ProjectID.ValueString()
	clientID := plan.ClientID.ValueString()
	accessListReq := NewMongoDBProjectServiceAccountAccessListEntry(&plan)

	connV2 := r.Client.AtlasV2
	firstPage, _, err := connV2.ServiceAccountsApi.CreateAccessList(ctx, projectID, clientID, accessListReq).ItemsPerPage(serviceaccountaccesslistentry.ItemsPerPage).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	cidrOrIP := getCidrOrIP(&plan)
	listPageFunc := func(ctx context.Context, pageNum int) (*admin.PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
		return connV2.ServiceAccountsApi.ListAccessList(ctx, projectID, clientID).PageNum(pageNum).ItemsPerPage(serviceaccountaccesslistentry.ItemsPerPage).Execute()
	}
	entry, _, err := serviceaccountaccesslistentry.ReadAccessListEntry(ctx, firstPage, listPageFunc, cidrOrIP)
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	accessListModel := NewTFProjectServiceAccountAccessListModel(projectID, clientID, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, accessListModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFProjectServiceAccountAccessListEntryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ProjectID.ValueString()
	clientID := state.ClientID.ValueString()
	cidrOrIP := getCidrOrIP(&state)

	connV2 := r.Client.AtlasV2
	listPageFunc := func(ctx context.Context, pageNum int) (*admin.PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
		return connV2.ServiceAccountsApi.ListAccessList(ctx, projectID, clientID).PageNum(pageNum).ItemsPerPage(serviceaccountaccesslistentry.ItemsPerPage).Execute()
	}
	entry, apiResp, err := serviceaccountaccesslistentry.ReadAccessListEntry(ctx, nil, listPageFunc, cidrOrIP)
	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}
	if entry == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	accessListModel := NewTFProjectServiceAccountAccessListModel(projectID, clientID, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, accessListModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFProjectServiceAccountAccessListEntryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ProjectID.ValueString()
	clientID := state.ClientID.ValueString()
	cidrOrIP := getCidrOrIP(&state)

	connV2 := r.Client.AtlasV2
	if _, err := connV2.ServiceAccountsApi.DeleteGroupAccessEntry(ctx, projectID, clientID, cidrOrIP).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"invalid import ID",
			"expected format: {project_id}/{client_id}/{cidr_block} or {project_id}/{client_id}/{ip_address}",
		)
		return
	}

	projectID, clientID, cidrOrIP := parts[0], parts[1], parts[2]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("client_id"), clientID)...)
	if strings.Contains(cidrOrIP, "/") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cidr_block"), cidrOrIP)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip_address"), cidrOrIP)...)
	}
}

func getCidrOrIP(model *TFProjectServiceAccountAccessListEntryModel) string {
	cidrOrIP := model.IPAddress.ValueString()
	if cidrOrIP == "" {
		cidrOrIP = model.CIDRBlock.ValueString()
	}
	return cidrOrIP
}
