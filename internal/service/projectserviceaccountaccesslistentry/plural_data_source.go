package projectserviceaccountaccesslistentry

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	serviceaccountaccesslistentry "github.com/mongodb/terraform-provider-mongodbatlas/internal/service/serviceaccountaccesslistentry"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
)

var _ datasource.DataSource = &pluralDS{}
var _ datasource.DataSourceWithConfigure = &pluralDS{}

const pluralDatasourceName = "project_service_account_access_list_entries"

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
		RequiredFields: []string{"project_id", "client_id"},
	})
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var conf TFProjectServiceAccountAccessListEntriesPluralDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &conf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := conf.ProjectID.ValueString()
	clientID := conf.ClientID.ValueString()

	api := d.Client.AtlasV2.ServiceAccountsApi
	entries, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.ServiceAccountIPAccessListEntry], *http.Response, error) {
		return api.ListAccessList(ctx, projectID, clientID).PageNum(pageNum).ItemsPerPage(serviceaccountaccesslistentry.ItemsPerPage).Execute()
	})
	if err != nil {
		resp.Diagnostics.AddError("error fetching list", err.Error())
		return
	}

	newProjectServiceAccountAccessListsModel := NewTFProjectServiceAccountAccessListEntriesPluralDSModel(projectID, clientID, entries)
	resp.Diagnostics.Append(resp.State.Set(ctx, newProjectServiceAccountAccessListsModel)...)
}
