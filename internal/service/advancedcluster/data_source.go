package advancedcluster

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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
	resp.Schema = dataSourceSchema(ctx)
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModelDS
	diags := &resp.Diagnostics
	diags.Append(req.Config.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	model := d.readCluster(ctx, diags, &state)
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (d *ds) readCluster(ctx context.Context, diags *diag.Diagnostics, modelDS *TFModelDS) *TFModelDS {
	clusterName := modelDS.Name.ValueString()
	projectID := modelDS.ProjectID.ValueString()
	clusterResp, flexClusterResp := GetClusterDetails(ctx, diags, projectID, clusterName, d.Client, false, modelDS.UseEffectiveFields.ValueBool())
	if diags.HasError() {
		return nil
	}
	if flexClusterResp == nil && clusterResp == nil {
		return nil
	}
	var result *TFModelDS
	if flexClusterResp != nil {
		result = convertFlexClusterToDS(ctx, diags, flexClusterResp)
	} else {
		result = convertBasicClusterToDS(ctx, diags, d.Client, clusterResp)
	}
	if result != nil {
		result.UseEffectiveFields = modelDS.UseEffectiveFields
	}
	return result
}

func convertFlexClusterToDS(ctx context.Context, diags *diag.Diagnostics, flexCluster *admin.FlexClusterDescription20241113) *TFModelDS {
	clusterDesc := FlexDescriptionToClusterDescription(flexCluster, nil)
	modelOutDS := newTFModelDS(ctx, clusterDesc, diags, nil)
	if diags.HasError() {
		return nil
	}
	modelOutDS.AdvancedConfiguration = types.ObjectNull(AdvancedConfigurationObjType.AttrTypes)
	return modelOutDS
}

func convertBasicClusterToDS(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, clusterResp *admin.ClusterDescription20240805) *TFModelDS {
	containerIDs := resolveContainerIDsOrError(ctx, diags, clusterResp, client.AtlasV2.NetworkPeeringApi)
	if diags.HasError() {
		return nil
	}
	modelOutDS := newTFModelDS(ctx, clusterResp, diags, containerIDs)
	if diags.HasError() {
		return nil
	}
	updateModelAdvancedConfigDS(ctx, diags, client, modelOutDS, &ProcessArgs{
		ArgsDefault:           nil,
		ClusterAdvancedConfig: clusterResp.AdvancedConfiguration,
	})
	if diags.HasError() {
		return nil
	}
	return modelOutDS
}
