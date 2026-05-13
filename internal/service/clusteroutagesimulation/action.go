package clusteroutagesimulation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
)

type OutageSimulationAction struct {
	client *config.MongoDBClient
}

type outageSimulationActionModel struct {
	ProjectID             types.String              `tfsdk:"project_id"`
	ClusterName           types.String              `tfsdk:"cluster_name"`
	Timeout               types.String              `tfsdk:"timeout"`
	OutageFilters         []outageFilterActionModel `tfsdk:"outage_filters"`
	DeleteOnCreateTimeout types.Bool                `tfsdk:"delete_on_create_timeout"`
}

type outageFilterActionModel struct {
	CloudProvider types.String `tfsdk:"cloud_provider"`
	RegionName    types.String `tfsdk:"region_name"`
}

func Action() action.Action {
	return &OutageSimulationAction{}
}

func (a *OutageSimulationAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_start_outage_simulation"
}

func (a *OutageSimulationAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		Description: "Starts a cluster outage simulation to test resilience. Simulates the outage of specified regions.",
		Attributes: map[string]actionschema.Attribute{
			"project_id": actionschema.StringAttribute{
				Required:    true,
				Description: "Unique 24-hexadecimal digit string that identifies the project.",
			},
			"cluster_name": actionschema.StringAttribute{
				Required:    true,
				Description: "Name of the cluster to simulate the outage on.",
			},
			"outage_filters": actionschema.ListNestedAttribute{
				Required:    true,
				Description: "List of regions to simulate outages for.",
				NestedObject: actionschema.NestedAttributeObject{
					Attributes: map[string]actionschema.Attribute{
						"cloud_provider": actionschema.StringAttribute{
							Required:    true,
							Description: "Cloud provider of the region to simulate the outage for (e.g. AWS, GCP, AZURE).",
						},
						"region_name": actionschema.StringAttribute{
							Required:    true,
							Description: "Name of the region to simulate the outage for (e.g. US_EAST_1).",
						},
					},
				},
			},
			"timeout": actionschema.StringAttribute{
				Optional:    true,
				Description: "Duration to wait for the simulation to reach SIMULATING state. Valid units: ns, us, ms, s, m, h (e.g. 25m, 1h30m). Defaults to 25m.",
			},
			"delete_on_create_timeout": actionschema.BoolAttribute{
				Optional:    true,
				Description: "If true and a timeout occurs while waiting for SIMULATING state, the simulation is ended before returning. Defaults to true.",
			},
		},
	}
}

func (a *OutageSimulationAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*config.MongoDBClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Action Configure Type",
			fmt.Sprintf("Expected *config.MongoDBClient, got: %T.", req.ProviderData),
		)
		return
	}
	a.client = client
}

func (a *OutageSimulationAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var model outageSimulationActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := model.ProjectID.ValueString()
	clusterName := model.ClusterName.ValueString()

	createTimeout, err := ParseActionTimeout(model.Timeout)
	if err != nil {
		resp.Diagnostics.AddError("Invalid timeout value", err.Error())
		return
	}

	deleteOnCreateTimeout := true
	if !model.DeleteOnCreateTimeout.IsNull() {
		deleteOnCreateTimeout = model.DeleteOnCreateTimeout.ValueBool()
	}

	filters := make([]admin.AtlasClusterOutageSimulationOutageFilter, len(model.OutageFilters))
	for i, f := range model.OutageFilters {
		filters[i] = admin.AtlasClusterOutageSimulationOutageFilter{
			CloudProvider: f.CloudProvider.ValueStringPointer(),
			RegionName:    f.RegionName.ValueStringPointer(),
			Type:          conversion.StringPtr(defaultOutageFilterType),
		}
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Starting outage simulation for cluster %s in project %s.", clusterName, projectID),
	})

	tc := retrystrategy.TimeConfig{Timeout: createTimeout, MinTimeout: timeout, Delay: timeout}
	if err := SimulateOutage(ctx, a.client.AtlasV2.ClusterOutageSimulationApi, projectID, clusterName, filters, deleteOnCreateTimeout, tc); err != nil {
		resp.Diagnostics.AddError("Error Starting Outage Simulation", err.Error())
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Outage simulation is now active for cluster %s.", clusterName),
	})
}
