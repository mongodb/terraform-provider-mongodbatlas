package serviceaccountaccesslistentry

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
)

var _ datasource.DataSource = &pluralDS{}
var _ datasource.DataSourceWithConfigure = &pluralDS{}

const pluralDatasourceName = "service_account_access_list_entries"

func PluralDataSource() datasource.DataSource {
	return &pluralDS{
		DSCommon: config.DSCommon{
			DataSourceName: pluralDatasourceName,
		},
	}
}

type pluralDS struct {
	config.DSCommon
}

func (d *pluralDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields: []string{"org_id", "client_id"},
	})
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var conf TFServiceAccountAccessListEntriesPluralDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := conf.OrgID.ValueString()
	clientID := conf.ClientID.ValueString()

	api := d.Client.AtlasV2.ServiceAccountsApi
	entries, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.ServiceAccountIPAccessListEntry], *http.Response, error) {
		return api.ListOrgAccessList(ctx, orgID, clientID).PageNum(pageNum).ItemsPerPage(ItemsPerPage).Execute()
	})
	if err != nil {
		resp.Diagnostics.AddError("error fetching list", err.Error())
		return
	}

	newServiceAccountAccessListsModel := NewTFServiceAccountAccessListEntriesPluralDSModel(orgID, clientID, entries)
	resp.Diagnostics.Append(resp.State.Set(ctx, newServiceAccountAccessListsModel)...)
}
