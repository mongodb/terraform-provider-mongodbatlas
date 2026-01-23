package serviceaccountaccesslistentry

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
)

var _ datasource.DataSource = &ds{}
var _ datasource.DataSourceWithConfigure = &ds{}

func DataSource() datasource.DataSource {
	return &ds{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type ds struct {
	config.DSCommon
}

func (d *ds) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"org_id", "client_id"},
		OverridenFields: map[string]dsschema.Attribute{
			"cidr_block": dsschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: cidrBlockDesc,
				Validators: []validator.String{
					validate.ValidCIDR(),
					stringvalidator.ConflictsWith(path.MatchRoot("ip_address")),
				},
			},
			"ip_address": dsschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: ipAddressDesc,
				Validators: []validator.String{
					validate.ValidIP(),
					stringvalidator.ConflictsWith(path.MatchRoot("cidr_block")),
				},
			},
		},
	})
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var conf TFServiceAccountAccessListEntryModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if conf.CIDRBlock.ValueString() == "" && conf.IPAddress.ValueString() == "" {
		resp.Diagnostics.AddError("validation error", "cidr_block or ip_address must be provided")
		return
	}

	orgID := conf.OrgID.ValueString()
	clientID := conf.ClientID.ValueString()
	cidrOrIP := getCidrOrIP(&conf)

	connV2 := d.Client.AtlasV2
	listPageFunc := func(ctx context.Context, pageNum int) (*admin.PaginatedServiceAccountIPAccessEntry, *http.Response, error) {
		return connV2.ServiceAccountsApi.ListOrgAccessList(ctx, orgID, clientID).PageNum(pageNum).ItemsPerPage(ItemsPerPage).Execute()
	}
	entry, _, err := ReadAccessListEntry(ctx, nil, listPageFunc, cidrOrIP)
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}
	if entry == nil {
		resp.Diagnostics.AddError("Resource not found", "The requested resource does not exist")
		return
	}

	accessListModel := NewTFServiceAccountAccessListModel(orgID, clientID, entry)
	resp.Diagnostics.Append(resp.State.Set(ctx, accessListModel)...)
}
