package serviceaccountaccesslistentry

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
)

const (
	resourceName = "service_account_access_list_entry"
	itemsPerPage = 500 // Max items per page
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
	var plan TFServiceAccountAccessListEntryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.CIDRBlock.ValueString() == "" && plan.IPAddress.ValueString() == "" {
		resp.Diagnostics.AddError("validation error", "cidr_block or ip_address must be provided")
		return
	}

	orgID := plan.OrgID.ValueString()
	clientID := plan.ClientID.ValueString()
	accessListReq := NewMongoDBServiceAccountAccessListEntry(&plan)

	connV2 := r.Client.AtlasV2
	firstPage, _, err := connV2.ServiceAccountsApi.CreateOrgAccessList(ctx, orgID, clientID, accessListReq).ItemsPerPage(itemsPerPage).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	cidrOrIP := getCidrOrIP(&plan)
	entry, _, err := readAccessListEntry(ctx, connV2.ServiceAccountsApi, firstPage, orgID, clientID, cidrOrIP)
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	accessListModel := NewTFServiceAccountAccessListModel(orgID, clientID, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, accessListModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFServiceAccountAccessListEntryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := state.OrgID.ValueString()
	clientID := state.ClientID.ValueString()
	cidrOrIP := getCidrOrIP(&state)

	connV2 := r.Client.AtlasV2
	entry, apiResp, err := readAccessListEntry(ctx, connV2.ServiceAccountsApi, nil, orgID, clientID, cidrOrIP)
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

	accessListModel := NewTFServiceAccountAccessListModel(orgID, clientID, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, accessListModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFServiceAccountAccessListEntryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := state.OrgID.ValueString()
	clientID := state.ClientID.ValueString()
	cidrOrIP := getCidrOrIP(&state)

	connV2 := r.Client.AtlasV2
	if _, err := connV2.ServiceAccountsApi.DeleteOrgAccessEntry(ctx, orgID, clientID, cidrOrIP).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"invalid import ID",
			"expected format: {org_id}/{client_id}/{cidr_block} or {org_id}/{client_id}/{ip_address}",
		)
		return
	}

	orgID, clientID, cidrOrIP := parts[0], parts[1], parts[2]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("client_id"), clientID)...)
	if strings.Contains(cidrOrIP, "/") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cidr_block"), cidrOrIP)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip_address"), cidrOrIP)...)
	}
}

func getCidrOrIP(model *TFServiceAccountAccessListEntryModel) string {
	cidrOrIP := model.IPAddress.ValueString()
	if cidrOrIP == "" {
		cidrOrIP = model.CIDRBlock.ValueString()
	}
	return cidrOrIP
}

// readAccessListEntry Iterates through access list pages looking for the entry.
// The first page can be provided to skip an API call. Useful for Create operation, which returns the first page.
func readAccessListEntry(
	ctx context.Context,
	api admin.ServiceAccountsApi,
	firstPage *admin.PaginatedServiceAccountIPAccessEntry,
	orgID, clientID, cidrOrIP string,
) (*admin.ServiceAccountIPAccessListEntry, *http.Response, error) {
	var err error
	var apiResp *http.Response

	count := 0
	page := firstPage
	for currentPage := 1; ; currentPage++ {
		if page == nil {
			page, apiResp, err = api.ListOrgAccessList(ctx, orgID, clientID).PageNum(currentPage).ItemsPerPage(itemsPerPage).Execute()
			if err != nil {
				return nil, apiResp, err
			}
		}

		results := page.GetResults()
		count += len(results)

		for i := range results {
			entry := &results[i]
			if entry.GetIpAddress() == cidrOrIP || entry.GetCidrBlock() == cidrOrIP {
				return entry, nil, nil
			}
		}

		if len(results) == 0 || count >= page.GetTotalCount() {
			break
		}
		page = nil
	}

	return nil, nil, nil
}
