package advancedcluster

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func (r *advancedClusterRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.Client.Atlas

	var state tfAdvancedClusterRSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID, clusterName := decodeClusterID(state.ID.ValueString())

	cluster, response, err := conn.AdvancedClusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to READ cluster. An error occurred when getting cluster details from Atlas", err.Error())
		return
	}

	newClusterModel, diags := newTfAdvancedClusterRSModel(ctx, conn, cluster, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterModel)...)
}

func getAdvancedClusterContainerID(containers []matlas.Container, cluster *matlas.AdvancedRegionConfig) string {
	if len(containers) != 0 {
		for i := range containers {
			if cluster.ProviderName == "GCP" {
				return containers[i].ID
			}

			if containers[i].ProviderName == cluster.ProviderName &&
				containers[i].Region == cluster.RegionName || // For Azure
				containers[i].RegionName == cluster.RegionName { // For AWS
				return containers[i].ID
			}
		}
	}

	return ""
}
