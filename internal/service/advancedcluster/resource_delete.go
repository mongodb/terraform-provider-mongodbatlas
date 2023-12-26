package advancedcluster

import (
	"context"
	"log"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func (r *advancedClusterRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	conn := r.Client.Atlas
	var state tfAdvancedClusterRSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := deleteCluster(ctx, conn, &state, timeout); err != nil {
		resp.Diagnostics.AddError("Unable to DELETE cluster. An error occurred when deleting cluster in Atlas", err.Error())
		return
	}
}

func deleteCluster(ctx context.Context, conn *matlas.Client, state *tfAdvancedClusterRSModel, timeout time.Duration) error {
	projectID, clusterName := decodeClusterID(state.ID.ValueString())

	var options *matlas.DeleteAdvanceClusterOptions
	if v := state.RetainBackupsEnabled; !v.IsNull() {
		options = &matlas.DeleteAdvanceClusterOptions{
			RetainBackups: v.ValueBoolPointer(),
		}
	}

	if _, err := conn.AdvancedClusters.Delete(ctx, projectID, clusterName, options); err != nil {
		return err
	}

	log.Println("[INFO] Waiting for MongoDB ClusterAdvanced to be destroyed")
	err := waitClusterDelete(ctx, conn, timeout, projectID, clusterName)

	return err
}

func waitClusterDelete(ctx context.Context, conn *matlas.Client, timeout time.Duration, projectID, clusterName string) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
