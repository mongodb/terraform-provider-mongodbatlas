package projectipaccesslist

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
)

var _ list.ListResource = &IPAccessListListResource{}
var _ list.ListResourceWithConfigure = &IPAccessListListResource{}

type IPAccessListListResource struct {
	client *config.MongoDBClient
}

type listConfigModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

func ListResource() list.ListResource {
	return &IPAccessListListResource{}
}

func (r *IPAccessListListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_ip_access_list"
}

func (r *IPAccessListListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"project_id": listschema.StringAttribute{
				Required:    true,
				Description: "Unique 24-hexadecimal digit string that identifies your project.",
			},
		},
	}
}

func (r *IPAccessListListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*config.MongoDBClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected List Resource Configure Type",
			fmt.Sprintf("Expected *config.MongoDBClient, got: %T.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *IPAccessListListResource) List(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream) {
	var cfg listConfigModel
	if diags := req.Config.Get(ctx, &cfg); diags.HasError() {
		resp.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	projectID := cfg.ProjectID.ValueString()
	allEntries, err := r.fetchAllEntries(ctx, projectID)
	if err != nil {
		var d diag.Diagnostics
		d.AddError("Error listing IP access list entries", err.Error())
		resp.Results = list.ListResultsStreamDiagnostics(d)
		return
	}

	resp.Results = func(push func(list.ListResult) bool) {
		for i := range allEntries {
			entry := &allEntries[i]
			result := req.NewListResult(ctx)

			entryVal := AccessListEntryValue(entry)
			result.Diagnostics.Append(result.Identity.SetAttribute(ctx, path.Root("project_id"), projectID)...)
			result.Diagnostics.Append(result.Identity.SetAttribute(ctx, path.Root("entry"), entryVal)...)
			result.DisplayName = fmt.Sprintf("%s (%s)", entryVal, projectID)

			if !push(result) {
				return
			}
		}
	}
}

func (r *IPAccessListListResource) fetchAllEntries(ctx context.Context, projectID string) ([]admin.NetworkPermissionEntry, error) {
	return FetchAllEntries(ctx, r.client.AtlasV2.ProjectIPAccessListApi, projectID)
}

func FetchAllEntries(ctx context.Context, api admin.ProjectIPAccessListApi, projectID string) ([]admin.NetworkPermissionEntry, error) {
	return dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.NetworkPermissionEntry], *http.Response, error) {
		return api.ListAccessListEntries(ctx, projectID).ItemsPerPage(500).PageNum(pageNum).Execute()
	})
}

func AccessListEntryValue(entry *admin.NetworkPermissionEntry) string {
	if v := entry.GetIpAddress(); v != "" {
		return v
	}
	if v := entry.GetCidrBlock(); v != "" {
		return v
	}
	return entry.GetAwsSecurityGroup()
}
