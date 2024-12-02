package resourcepolicy

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &resourcePolicysDS{}
var _ datasource.DataSourceWithConfigure = &resourcePolicysDS{}

const (
	dataSourcePluralName     = "resource_policies"
	fullDataSourcePluralName = "mongodbatlas_" + dataSourcePluralName
	errorReadDSP             = "error reading plural data source " + fullDataSourcePluralName
)

func PluralDataSource() datasource.DataSource {
	return &resourcePolicysDS{
		DSCommon: config.DSCommon{
			DataSourceName: dataSourcePluralName,
		},
	}
}

type resourcePolicysDS struct {
	config.DSCommon
}

func (d *resourcePolicysDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: THIS WILL BE REMOVED BEFORE MERGING, check old data source schema and new auto-generated schema are the same
	ds1 := DataSourcePluralSchema(ctx)
	conversion.UpdateSchemaDescription(&ds1)
	requiredFields := []string{"org_id"}
	ds2 := conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), requiredFields, nil, nil, "", false)
	clone := conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), requiredFields, nil, nil, "", false)
	resourcePolicies := clone.Attributes["results"].(schema.ListNestedAttribute)
	resourcePolicies.DeprecationMessage = fmt.Sprintf(constant.DeprecationParamWithReplacement, "`results`")
	resourcePolicies.Description = ""
	resourcePolicies.MarkdownDescription = ""
	ds2.Attributes["resource_policies"] = resourcePolicies
	if diff := cmp.Diff(ds1, ds2); diff != "" {
		log.Fatal(diff)
	}
	resp.Schema = ds2
}

func (d *resourcePolicysDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg TFModelDSP
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	orgID := cfg.OrgID.ValueString()
	apiResp, _, err := connV2.ResourcePoliciesApi.GetAtlasResourcePolicies(ctx, orgID).Execute()

	if err != nil {
		resp.Diagnostics.AddError(errorReadDSP, err.Error())
		return
	}

	newResourcePolicysModel, diags := NewTFModelDSP(ctx, orgID, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicysModel)...)
}
