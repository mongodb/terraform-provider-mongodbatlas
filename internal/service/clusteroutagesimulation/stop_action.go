package clusteroutagesimulation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

type StopOutageSimulationAction struct {
	client *config.MongoDBClient
}

type stopOutageSimulationActionModel struct {
	ProjectID   types.String `tfsdk:"project_id"`
	ClusterName types.String `tfsdk:"cluster_name"`
	Timeout     types.String `tfsdk:"timeout"`
}

func StopAction() action.Action {
	return &StopOutageSimulationAction{}
}

func (a *StopOutageSimulationAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stop_outage_simulation"
}

func (a *StopOutageSimulationAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		Description: "Stops an active cluster outage simulation and waits for the cluster to recover.",
		Attributes: map[string]actionschema.Attribute{
			"project_id": actionschema.StringAttribute{
				Required:    true,
				Description: "Unique 24-hexadecimal digit string that identifies the project.",
			},
			"cluster_name": actionschema.StringAttribute{
				Required:    true,
				Description: "Name of the cluster whose outage simulation to stop.",
			},
			"timeout": actionschema.StringAttribute{
				Optional:    true,
				Description: "Duration to wait for the simulation to reach DELETED state. Valid units: ns, us, ms, s, m, h (e.g. 25m, 1h30m). Defaults to 25m.",
			},
		},
	}
}

func (a *StopOutageSimulationAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

func (a *StopOutageSimulationAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var model stopOutageSimulationActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := model.ProjectID.ValueString()
	clusterName := model.ClusterName.ValueString()

	deleteTimeout, err := ParseActionTimeout(model.Timeout)
	if err != nil {
		resp.Diagnostics.AddError("Invalid timeout value", err.Error())
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Stopping outage simulation for cluster %s in project %s.", clusterName, projectID),
	})

	tc := retrystrategy.TimeConfig{Timeout: deleteTimeout, MinTimeout: timeout, Delay: timeout}
	if err := StopSimulation(ctx, a.client.AtlasV2.ClusterOutageSimulationApi, projectID, clusterName, tc); err != nil {
		resp.Diagnostics.AddError("Error Stopping Outage Simulation", err.Error())
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Outage simulation stopped. Cluster %s has recovered.", clusterName),
	})
}
